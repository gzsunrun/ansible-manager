package project

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/api/s3"
	"github.com/gzsunrun/ansible-manager/config"
)

func CloneGitRepo(w http.ResponseWriter, r *http.Request){
	repoName := r.FormValue("repo_name")
	repoDesc := r.FormValue("repo_desc")
	url := r.FormValue("repo_path")
	if repoName == ""||url=="" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	repoPath := strconv.FormatInt(time.Now().UnixNano(), 10)
	err:=s3.GitClone(url,repoPath)
	if err != nil {
		logs.Error("import",err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if config.Cfg.AnsibleManager.S3Status {
		defer os.Remove(config.Cfg.AnsibleManager.WorkPath + "/repo/" + repoPath)
	}
	vst, err := readVars(config.Cfg.AnsibleManager.WorkPath + "/repo/" + repoPath)
	if err != nil {
		logs.Error("read vars",err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fi,err:=os.Stat(config.Cfg.AnsibleManager.WorkPath + "/repo/" + repoPath)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	c, err := os.Open(config.Cfg.AnsibleManager.WorkPath + "/repo/" + repoPath)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer c.Close()
	err = s3.S3Put(c, fi.Size(), repoPath)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = InsertVars(vst, repoName, repoPath, repoDesc)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}