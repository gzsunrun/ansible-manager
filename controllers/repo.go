package controllers

import (
	"io/ioutil"
	"os"

	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/function"
	"github.com/gzsunrun/ansible-manager/core/helm"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/storage"
	"github.com/gzsunrun/ansible-manager/core/template"
	"github.com/hashwing/log"
)

// RepoController repo controller
type RepoController struct {
	BaseController
}

// List get repo list
func (c *RepoController) List() {
	defer c.ServeJSON()
	class := c.GetString("repo_type", "ansible")
	var repos []orm.RepositoryList
	err := orm.FindRepos(&repos, class)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, repos, 200)
}

func (c *RepoController) Icon() {
	id := c.GetString("id")
	repoParse := &storage.StorageParse{
		RemotePath: id,
	}

	data, contentType, err := storage.Storage.GetIO(repoParse)
	if err != nil {
		log.Error(err)
		c.SetResult(err, nil, 500)
		return
	}
	c.Ctx.Output.ContentType(contentType)
	c.Ctx.Output.Body(data)
}

// Create create repo
func (c *RepoController) Create() {
	defer c.ServeJSON()
	repo := orm.RepositoryInsert{}
	repo.ID = function.NewUuidString()
	f, _, err := c.GetFile("repo_path")
	if err != nil {
		log.Error("Getfile", err)
		c.SetResult(err, nil, 400)
		return
	}
	f.Close()
	repoPath := function.NewUuidString()
	repo.Path = repoPath
	err = c.SaveToFile("repo_path", config.Cfg.Common.WorkPath+"/"+repoPath)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(400, err.Error())
		return
	}
	defer os.Remove(config.Cfg.Common.WorkPath + "/" + repoPath)
	//defer os.RemoveAll(config.Cfg.Common.WorkPath + "/" + repoPath + "_dir")
	// err = function.ReadVars(config.Cfg.Common.WorkPath+"/"+repoPath, &repo)
	// if err != nil {
	// 	c.SetResult(err, nil, 400)
	// 	return
	// }
	var tpls []orm.RepositoryInsert
	if c.GetString("repo_type") == "helm" {
		tpls, err = helm.ReadChart(config.Cfg.Common.WorkPath+"/"+repoPath, repoPath)
	} else {
		tpls, err = template.ReadVars(config.Cfg.Common.WorkPath+"/"+repoPath, repoPath)
	}
	if err != nil {
		log.Error(err)
		c.SetErrMsg(400, err.Error())
		return
	}

	if c.GetString("repo_name") != "" {
		repo.Name = c.GetString("repo_name")
	}
	if c.GetString("repo_desc") != "" {
		repo.Desc = c.GetString("repo_desc")
	}
	logoInfo, _ := ioutil.ReadDir(config.Cfg.Common.WorkPath + "/" + repoPath + "_dir/logo")
	for _, info := range logoInfo {
		log.Debug(info.Name())
		logoParse := storage.StorageParse{
			LocalPath:  config.Cfg.Common.WorkPath + "/" + repoPath + "_dir/logo/" + info.Name(),
			RemotePath: "logos/" + info.Name(),
		}
		err = storage.Storage.Put(&logoParse)
		if err != nil {
			c.SetErrMsg(500, err.Error())
			return
		}
	}
	repoParse := storage.StorageParse{
		LocalPath:  config.Cfg.Common.WorkPath + "/" + repoPath,
		RemotePath: repoPath,
	}
	err = storage.Storage.Put(&repoParse)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	err = orm.CreateRepos(tpls)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	// err = orm.CreateRepo(repo)
	// if err != nil {
	// 	c.SetResult(err, nil, 400)
	// 	return
	// }
	c.SetResult(nil, nil, 204)
}

// Delete delete a repo
func (c *RepoController) Delete() {
	repoID := c.GetString("repo_id")
	var repo orm.Repository
	err := orm.GetRepoByID(repoID, &repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	err = orm.DelRepoByID(repoID)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if !orm.GetRepoByPath(repo.Path) {
		logoParse := storage.StorageParse{
			RemotePath: repo.Path + ".png",
		}
		repoParse := storage.StorageParse{
			RemotePath: repo.Path,
		}
		storage.Storage.Delete(&logoParse)
		storage.Storage.Delete(&repoParse)
	}
	c.SetResult(nil, nil, 204)
}

// Vars get repo vars
func (c *RepoController) Vars() {
	defer c.ServeJSON()
	rid := c.GetString("repo_id")
	var repo orm.Repository
	err := orm.GetRepoByID(rid, &repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	data := map[string]interface{}{
		"vars":  repo.Vars,
		"group": repo.Group,
		"tag":   repo.Tag,
	}
	c.SetResult(nil, data, 200)
}

// SyncGit clone repo from git
func (c *RepoController) SyncGit() {
	repo := orm.RepositoryInsert{}
	repo.ID = function.NewUuidString()
	repo.Path = c.GetString("git_url")
	repoPath := function.NewUuidString()
	repoParse := storage.StorageParse{
		RemotePath: repo.Path,
		LocalPath:  config.Cfg.Common.WorkPath + "/" + repoPath,
	}
	err := storage.Storage.Get(&repoParse)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}

	defer os.Remove(config.Cfg.Common.WorkPath + "/" + repoPath)
	defer os.RemoveAll(config.Cfg.Common.WorkPath + "/" + repoPath + "_dir")
	err = function.ReadVars(config.Cfg.Common.WorkPath+"/"+repoPath, &repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	err = orm.CreateRepo(repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

// StorageType get the type of storage
func (c *RepoController) StorageType() {
	defer c.ServeJSON()
	c.SetResult(nil, config.Cfg.Git.Enable, 200, "status")
}

// Health check health
func (c *RepoController) Health() {
	c.Data["json"] = map[string]interface{}{"health": "ok"}
	c.ServeJSON()
}
