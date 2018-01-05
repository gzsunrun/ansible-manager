package project

import (
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/api/db"
)

func CreateTemplate(w http.ResponseWriter, r *http.Request) {
	projectID := r.FormValue("project_id")
	repoID := r.FormValue("repo_id")
	tplName := r.FormValue("tpl_name")
	playbook := r.FormValue("playbook")
	playbookParse := r.FormValue("playbook_parse")
	if tplName == "" || repoID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("insert into ansible_template (project_id,repo_id,tpl_name,playbook,playbook_parse) values (?,?,?,?,?)  ", projectID, repoID, tplName, playbook, playbookParse)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetTemplate(w http.ResponseWriter, r *http.Request) {
	projectID := r.FormValue("project_id")
	var tpls []db.Template
	err := db.MysqlDB.Table("ansible_template").Where("project_id =?", projectID).Find(&tpls)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, tpls)
}

func GetTemplateByID(w http.ResponseWriter, r *http.Request) {
	tplID := r.FormValue("tpl_id")
	var tpl db.Template
	if tplID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Table("ansible_template").Where("tpl_id=?", tplID).Get(&tpl)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, tpl)
}

func DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	tplID := r.FormValue("tpl_id")
	if tplID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("delete from ansible_template where tpl_id=?", tplID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	tplID := r.FormValue("tpl_id")
	repoID := r.FormValue("repo_id")
	tplName := r.FormValue("tpl_name")
	playbook := r.FormValue("playbook")
	playbookParse := r.FormValue("playbook_parse")
	if tplName == "" || tplID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("update ansible_template set  repo_id=? ,tpl_name=? ,playbook=?,playbook_parse=? where tpl_id=?", repoID, tplName, playbook, playbookParse, tplID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
