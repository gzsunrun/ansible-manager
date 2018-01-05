package project

import (
	"net/http"

	"github.com/gorilla/context"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/api/db"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "Claims").(*MyCustomClaims)
	UID := claims.UserID
	projectName := r.FormValue("project_name")
	if projectName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("insert into ansible_project (user_id,project_name,project_created) values (?,?,NOW())", UID, projectName)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetProject(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "Claims").(*MyCustomClaims)
	UID := claims.UserID
	var projects []db.Project
	err := db.MysqlDB.Table("ansible_project").Where("user_id =?", UID).Find(&projects)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, projects)
}

func GetProjectByID(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "Claims").(*MyCustomClaims)
	UID := claims.UserID
	projectID := r.FormValue("project_id")
	var project db.Project
	if projectID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Table("ansible_project").Where("user_id =? and project_id=?", UID, projectID).Get(&project)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, project)
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "Claims").(*MyCustomClaims)
	UID := claims.UserID
	projectID := r.FormValue("project_id")
	if projectID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("delete from ansible_project where project_id=? and user_id=?", projectID, UID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.FormValue("project_id")
	projectName := r.FormValue("project_name")
	if projectName == "" || projectID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("update ansible_project set project_name = ? where project_id =?  ", projectName, projectID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
