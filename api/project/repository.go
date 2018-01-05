package project

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/api/db"
	"github.com/gzsunrun/ansible-manager/api/s3"
	"github.com/gzsunrun/ansible-manager/config"
)

func CreateRepository(w http.ResponseWriter, r *http.Request) {
	f, h, err := r.FormFile("repo_path")
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	repoName := r.FormValue("repo_name")
	repoDesc := r.FormValue("repo_desc")
	if repoName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	repoPath := strconv.FormatInt(time.Now().UnixNano(), 10)
	fmt.Println(repoPath)
	t, err := os.Create(config.Cfg.AnsibleManager.WorkPath + "/repo/" + repoPath)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := io.Copy(t, f); err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		t.Close()
		return
	}
	t.Close()
	if config.Cfg.AnsibleManager.S3Status {
		defer os.Remove(config.Cfg.AnsibleManager.WorkPath + "/repo/" + repoPath)
	}
	vst, err := readVars(config.Cfg.AnsibleManager.WorkPath + "/repo/" + repoPath)
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
	err = s3.S3Put(c, h.Size, repoPath)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fmt.Println(repoPath)
	err = InsertVars(vst, repoName, repoPath, repoDesc)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetRepository(w http.ResponseWriter, r *http.Request) {
	var repos []db.Repository
	err := db.MysqlDB.Table("ansible_repository").Find(&repos)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, repos)
}

func GetRepositoryID(w http.ResponseWriter, r *http.Request) {
	repoID := r.FormValue("repo_id")
	var repo db.Repository
	if repoID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Table("ansible_repository").Where("repo_id=?", repoID).Get(&repo)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, repo)
}

func DeleteRepository(w http.ResponseWriter, r *http.Request) {
	repoID := r.FormValue("repo_id")
	if repoID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var object string
	_, err := db.MysqlDB.Table("ansible_repository").Where("repo_id=?", repoID).Cols("repo_path").Get(&object)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = db.MysqlDB.Exec("delete from ansible_repository where repo_id=?", repoID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s3.S3Delte(object)
	w.WriteHeader(http.StatusNoContent)
}
