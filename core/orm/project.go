package orm

import (
	"time"

	log "github.com/astaxie/beego/logs"
)

type Project struct {
	ID 			string 					`xorm:"project_id" json:"project_id"`
	UserID 		string 					`xorm:"user_id" json:"-"`
	Name 		string 					`xorm:"project_name" json:"project_name"`
	Created	    time.Time  				`xorm:"created" json:"created"`
}

type ProjectHost struct {
	ProjectID 	string `xorm:"project_id" json:"project_id"`
	HostID 		string `xorm:"host_id" json:"host_id"`
}

func GetProject(pid string,project interface{})(bool,error){
	res,err:=MysqlDB.Table("ansible_project").Where("project_id=?",pid).Get(project)
	if err!=nil{
		log.Error(err)
	}
	return res,err
}

func FindProject(uid string,project interface{})error{
	err:=MysqlDB.Table("ansible_project").Where("user_id=?",uid).Find(project)
	if err!=nil{
		log.Error(err)
		return err
	}
	return nil
}

func CreateProject(project *Project) error{
	_,err:=MysqlDB.Table("ansible_project").Insert(project)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func DelProject(pid string)error{
	project:=new(Project)
	_,err:=MysqlDB.Table("ansible_project").Where("project_id=?",pid).Delete(project)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func UPdateProject(project *Project)error{
	_,err:=MysqlDB.Table("ansible_project").Where("project_id=?",project.ID).Update(project)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func DelHostFormProject(pH *ProjectHost)error{
	projectHost:=new(ProjectHost)
	_,err:=MysqlDB.Table("ansible_project_host").Where("host_id=? and project_id=?",pH.HostID,pH.ProjectID).Delete(projectHost)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func AddHostToProject(projectHost *[]ProjectHost)error{
	_,err:=MysqlDB.Table("ansible_project_host").Insert(projectHost)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func FindProjectHost(pid string)(*[]ProjectHost,error){
	var phs []ProjectHost
	err:=MysqlDB.Table("ansible_project_host").Where("project_id=?",pid).Find(&phs)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	return &phs,nil
}


func DelAllHostsByPid(pid string)error{
	project:=new(ProjectHost)
	_,err:=MysqlDB.Table("ansible_project_host").Where("project_id=?",pid).Delete(project)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func FindHostFromProject(pid string,hosts interface{})error{
	err:=MysqlDB.Table("ansible_project_host").
	Join("INNER","ansible_host","ansible_host.host_id=ansible_project_host.host_id").
	Where("project_id=?",pid).Find(hosts)
	if err!=nil{
		log.Error(err)
	}
	return err
}