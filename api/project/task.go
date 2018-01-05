package project

import (
	"net/http"
	"strconv"

	"github.com/gorilla/context"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/api/db"
	"github.com/gzsunrun/ansible-manager/api/tasks"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "Claims").(*MyCustomClaims)
	UID := claims.UserID
	taskName := r.FormValue("task_name")
	tplID := r.FormValue("tpl_id")
	pTag := r.FormValue("playbook_tag")
	if taskName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res, err := db.MysqlDB.Exec("insert into ansible_task (user_id,tpl_id,task_name,task_status,task_created,playbook_tag) values (?,?,?,'waitting',NOW(),?)  ", UID, tplID, taskName, pTag)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	taskID, err := res.LastInsertId()
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	go tasks.AddTask(int(taskID))
	obj := map[string]int{
		"task_id": int(taskID),
	}
	JsonWrite(w, http.StatusOK, obj)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "Claims").(*MyCustomClaims)
	UID := claims.UserID
	var tasks []db.Task
	err := db.MysqlDB.Table("ansible_task").Where("user_id =?", UID).Find(&tasks)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, tasks)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("task_id")
	var task db.Task
	if taskID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Table("ansible_task").Where("task_id=?", taskID).Get(&task)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("task_id")
	if taskID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("delete from ansible_task where task_id=?", taskID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("task_id")
	if taskID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var outputs []db.Output
	err := db.MysqlDB.Table("ansible_task_output").Where("task_id=?", taskID).Find(&outputs)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, outputs)
}

func StopTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.FormValue("task_id")
	if taskID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := strconv.Atoi(taskID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tasks.StopTask(id)
	w.WriteHeader(204)
}
