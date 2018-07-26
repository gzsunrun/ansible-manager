package controllers

import (
	"os"

	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/function"
	"github.com/gzsunrun/ansible-manager/core/template"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/storage"
	"github.com/satori/go.uuid"
)

// RepoController repo controller
type RepoController struct {
	BaseController
}

// List get repo list
func (c *RepoController) List() {
	defer c.ServeJSON()
	var repos []orm.RepositoryList
	err := orm.FindRepos(&repos)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, repos, 200)
}

// Create create repo
func (c *RepoController) Create() {
	defer c.ServeJSON()
	repo := orm.RepositoryInsert{}
	repo.ID = uuid.NewV4().String()
	f, _, err := c.GetFile("repo_path")
	if err != nil {
		log.Error("Getfile", err)
		c.SetResult(err, nil, 400)
		return
	}
	f.Close()
	repoPath := uuid.NewV4().String()
	repo.Path = repoPath
	err = c.SaveToFile("repo_path", config.Cfg.Common.WorkPath+"/"+repoPath)
	if err != nil {
		log.Error(err)
		c.SetResult(err, nil, 400)
		return
	}
	defer os.Remove(config.Cfg.Common.WorkPath + "/" + repoPath)
	defer os.RemoveAll(config.Cfg.Common.WorkPath + "/" + repoPath + "_dir")
	// err = function.ReadVars(config.Cfg.Common.WorkPath+"/"+repoPath, &repo)
	// if err != nil {
	// 	c.SetResult(err, nil, 400)
	// 	return
	// }
	tpls,err:=template.ReadVars(config.Cfg.Common.WorkPath+"/"+repoPath,repoPath)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}

	if c.GetString("repo_name") != "" {
		repo.Name = c.GetString("repo_name")
	}
	if c.GetString("repo_desc") != "" {
		repo.Desc = c.GetString("repo_desc")
	}
	_, err = os.Stat(config.Cfg.Common.WorkPath + "/" + repoPath + "_dir/logo.png")
	if err == nil || os.IsExist(err) {
		logoParse := storage.StorageParse{
			LocalPath:  config.Cfg.Common.WorkPath + "/" + repoPath + "_dir/logo.png",
			RemotePath: repoPath + ".png",
		}
		err = storage.Storage.Put(&logoParse)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
	} 
	repoParse := storage.StorageParse{
		LocalPath:  config.Cfg.Common.WorkPath + "/" + repoPath,
		RemotePath: repoPath,
	}
	err = storage.Storage.Put(&repoParse)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	err = orm.CreateRepos(tpls)
	if err != nil {
		c.SetResult(err, nil, 400)
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
	if !orm.GetRepoByPath(repo.Path){
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
	repo.ID = uuid.NewV4().String()
	repo.Path = c.GetString("git_url")
	repoPath := uuid.NewV4().String()
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
