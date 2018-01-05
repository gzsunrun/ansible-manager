package project

import (
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/api/db"
)

func InsertVars(vst *varsStruct, repoName, repoPath, repoDesc string) error {
	session := db.MysqlDB.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	res, err := session.Exec("insert into ansible_repository (repo_name,repo_path,repo_desc) values (?,?,?)  ", repoName, repoPath, repoDesc)
	if err != nil {
		logs.Error(err)
		return err
	}
	repoID, _ := res.LastInsertId()
	_, err = session.Exec("insert into ansible_vars (repo_id,vars_type,vars_name,vars_path,vars_value) values (?,?,?,?,?)  ", repoID, "tag", "tag", "tag.json", vst.tag)
	if err != nil {
		logs.Error(err)
		return err
	}
	_, err = session.Exec("insert into ansible_vars (repo_id,vars_type,vars_name,vars_path,vars_value) values (?,?,?,?,?)  ", repoID, "group", "group", "group.json", vst.group)
	if err != nil {
		logs.Error(err)
		return err
	}
	for _, v := range vst.vars {
		_, err = session.Exec("insert into ansible_vars (repo_id,vars_type,vars_name,vars_path,vars_value) values (?,?,?,?,?)  ", repoID, "vars", v["name"], v["path"], v["value"])
		if err != nil {
			logs.Error(err)
			return err
		}
	}
	return session.Commit()
}

func CreateVars(w http.ResponseWriter, r *http.Request) {
	varsType := r.FormValue("vars_type")
	varsName := r.FormValue("vars_name")
	varsValue := r.FormValue("vars_value")
	varsPath := r.FormValue("vars_path")
	repoID := r.FormValue("repo_id")
	if varsName == "" || varsValue == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("insert into ansible_vars (repo_id,vars_type,vars_name,vars_path,vars_value) values (?,?,?,?,?)  ", repoID, varsType, varsName, varsPath, varsValue)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetVars(w http.ResponseWriter, r *http.Request) {
	repoID := r.FormValue("repo_id")
	if repoID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var vars []db.Vars
	err := db.MysqlDB.Table("ansible_vars").Where("repo_id =?", repoID).Find(&vars)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, vars)
}

func GetVarsByID(w http.ResponseWriter, r *http.Request) {
	varsID := r.FormValue("vars_id")
	var vars db.Vars
	if varsID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Table("ansible_vars").Where("vars_id=?", varsID).Get(&vars)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, vars)
}

func DeleteVars(w http.ResponseWriter, r *http.Request) {
	varsID := r.FormValue("vars_id")
	if varsID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("delete from ansible_vars where vars_id=?", varsID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func UpdateVars(w http.ResponseWriter, r *http.Request) {
	varsName := r.FormValue("vars_name")
	varsValue := r.FormValue("vars_value")
	varsPath := r.FormValue("vars_path")
	varsID := r.FormValue("vars_id")
	if varsName == "" || varsID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err := db.MysqlDB.Exec("update ansible_vars set vars_name = ?,vars_path = ?,vars_value = ?  where vars_id =?  ", varsName, varsPath, varsValue, varsID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetTagByTplID(w http.ResponseWriter, r *http.Request) {
	tplID := r.FormValue("tpl_id")
	if tplID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var repoID int
	_, err := db.MysqlDB.Table("ansible_template").Where("tpl_id=?", tplID).Cols("repo_id").Get(&repoID)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var tag db.Vars
	_, err = db.MysqlDB.Table("ansible_vars").Where("repo_id=?", repoID).Where("vars_type=?", "tag").Get(&tag)
	if err != nil {
		logs.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	JsonWrite(w, http.StatusOK, tag)
}
