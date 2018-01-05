package project

import (
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/api/db"
)

func CreateHost(w http.ResponseWriter, r *http.Request) {
	hostAlias := r.FormValue("host_alias")
	hostName := r.FormValue("host_name")
	hostIP := r.FormValue("host_ip")
	hostUser := r.FormValue("host_user")
	hostPassword := r.FormValue("host_password")
	projectID := r.FormValue("project_id")
	hostKey := r.FormValue("host_key")
	if hostName == "" || hostIP == "" || projectID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("insert into ansible_host (project_id,host_alias,host_name,host_ip,host_user,host_password,host_key) values (?,?,?,?,?,?,?)  ", projectID, hostAlias, hostName, hostIP, hostUser, hostPassword, hostKey)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetHost(w http.ResponseWriter, r *http.Request) {
	projectID := r.FormValue("project_id")
	if projectID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var hosts []db.Host
	err := db.MysqlDB.Table("ansible_host").Where("project_id =?", projectID).Find(&hosts)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, hosts)
}

func GetHostByID(w http.ResponseWriter, r *http.Request) {
	hostID := r.FormValue("host_id")
	var host db.Host
	if hostID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Table("ansible_host").Where("host_id=?", hostID).Get(&host)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, host)
}

func DeleteHost(w http.ResponseWriter, r *http.Request) {
	hostID := r.FormValue("host_id")
	if hostID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("delete from ansible_host where host_id=?", hostID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateHost(w http.ResponseWriter, r *http.Request) {
	hostID := r.FormValue("host_id")
	hostAlias := r.FormValue("host_alias")
	hostName := r.FormValue("host_name")
	hostIP := r.FormValue("host_ip")
	hostUser := r.FormValue("host_user")
	hostPassword := r.FormValue("host_password")
	projectID := r.FormValue("project_id")
	hostKey := r.FormValue("host_key")
	if hostName == "" || hostIP == "" || projectID == "" || hostID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("update ansible_host set host_name = ?,host_ip = ?,host_user = ?,host_password=?,host_key=?,project_id=?,host_alias=?  where host_id =?  ", hostName, hostIP, hostUser, hostPassword, hostKey, projectID, hostAlias, hostID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
