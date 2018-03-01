package controllers

import (
	"os"
	"strconv"
	"time"

	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/function"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/satori/go.uuid"
)

type RepoController struct {
	BaseController
}

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

func (c *RepoController) Create() {
	defer c.ServeJSON()
	repo := orm.RepositoryInsert{}
	repo.Name=c.GetString("repo_name")
	repo.Desc=c.GetString("repo_desc")
	repo.ID= uuid.NewV4().String()
	f, h, err := c.GetFile("repo_path")
	if err != nil {
		log.Error("Getfile",err)
		c.SetResult(err, nil, 400)
		return
	}
	f.Close()
	repoPath := strconv.FormatInt(time.Now().UnixNano(), 10)
	repo.Path = repoPath
	err = c.SaveToFile("repo_path", config.Cfg.Ansible.WorkPath+"/"+repoPath)
	if err != nil {
		log.Error(err)
		c.SetResult(err, nil, 400)
		return
	}
	defer os.Remove(config.Cfg.Ansible.WorkPath + "/" + repoPath)
	defer os.RemoveAll(config.Cfg.Ansible.WorkPath + "/" + repoPath + "_dir")
	err = function.ReadVars(config.Cfg.Ansible.WorkPath + "/" + repoPath,&repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	icof, err := os.Open(config.Cfg.Ansible.WorkPath + "/" + repoPath + "_dir/logo.png")
	if err != nil {
		log.Error(err)
		c.SetResult(err, nil, 400)
		return
	}
	defer icof.Close()
	fileStat, err := icof.Stat()
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	err = function.S3Put(icof, fileStat.Size(), repoPath+".png")
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	rf, err := os.Open(config.Cfg.Ansible.WorkPath + "/" + repoPath)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	defer rf.Close()
	err = function.S3Put(rf, h.Size, repoPath)
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

func (c *RepoController) Delete() {
	repoID := c.GetString("repo_id")
	var repo orm.Repository
	err := orm.GetRepoByID(repoID,&repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	err = orm.DelRepoByID(repoID)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	function.S3Delte(repo.Path)
	function.S3Delte(repo.Path + ".png")
	c.SetResult(nil, nil, 204)
}

func (c *RepoController) Vars(){
	defer c.ServeJSON()
	rid:=c.GetString("repo_id")
	var repo orm.Repository
	err:=orm.GetRepoByID(rid,&repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	data :=map[string]interface{}{
		"vars": repo.Vars,
		"group": repo.Group,
		"tag": repo.Tag,
	}
	c.SetResult(nil,data,200)
}

func (c *RepoController) Health() {
	c.Data["json"] = map[string]interface{}{"health": "ok"}
	c.ServeJSON()
}