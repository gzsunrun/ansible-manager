package orm

import (
	"time"

	log "github.com/astaxie/beego/logs"
)


type Hosts struct {
	ID 			string 					`xorm:"host_id" json:"host_id"`
	UserID 		string 					`xorm:"user_id" json:"user_id"`
	Alias 		string 					`xorm:"host_alias" json:"host_alias"`
	HostName 	string 					`xorm:"host_name" json:"host_name"`
	IP  		string 					`xorm:"host_ip" json:"host_ip"`
	User 		string 					`xorm:"host_user" json:"host_user"`
	Password 	string 					`xorm:"host_password" json:"host_password"`
	Key 		string 					`xorm:"host_key" json:"host_key"`
	Status		bool					`xorm:"host_status" json:"host_status"`
	Created		time.Time  				`xorm:"created" json:"created"`
}

type HostsList struct {
	ID 			string 					`xorm:"host_id" json:"host_id"`
	UserID 		string 					`xorm:"user_id" json:"-"`
	Alias 		string 					`xorm:"host_alias" json:"host_alias"`
	Name 		string 					`xorm:"host_name" json:"host_name"`
	IP  		string 					`xorm:"host_ip" json:"host_ip"`
	User 		string 					`xorm:"host_user" json:"-"`
	Password 	string 					`xorm:"host_password" json:"-"`
	Key 		string 					`xorm:"host_key" json:"-"`
	Status		bool					`xorm:"host_status" json:"host_status"`
	Created		time.Time  				`xorm:"created" json:"created"`
}

func CreateHost(host *Hosts) error{
	_,err:=MysqlDB.Table("ansible_host").Insert(host)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func CreateHostList(host *HostsList) error{
	_,err:=MysqlDB.Table("ansible_host").Insert(host)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func UPdateHost(host *Hosts)error{
	_,err:=MysqlDB.Table("ansible_host").Where("host_id=?",host.ID).Update(host)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func UPdateNullHost(host *Hosts)error{
	_,err:=MysqlDB.Table("ansible_host").Where("host_id=?",host.ID).
	Cols("host_alias","host_name","host_ip","host_user","host_password","host_key","host_status").
	Update(host)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func UPdateAuthHost(host *Hosts)error{
	var err error
	if host.HostName==""{
		_,err=MysqlDB.Table("ansible_host").Where("host_id=?",host.ID).
		Cols("host_user","host_password","host_key").
		Update(host)
	}else{
		_,err=MysqlDB.Table("ansible_host").Where("host_id=?",host.ID).
		Cols("host_name","host_user","host_password","host_key").
		Update(host)
	}
	if err!=nil{
		log.Error(err)
	}
	return err
}

func FindHosts(uid string,hosts interface{})error{
	err:=MysqlDB.Table("ansible_host").Where("user_id=?",uid).Find(hosts)
	if err!=nil{
		log.Error(err)
		return err
	}
	return nil
}

func GetHost(hostID string,host interface{})(bool,error){
	res,err:=MysqlDB.Table("ansible_host").Where("host_id=?",hostID).Get(host)
	if err!=nil{
		log.Error(err)
	}
	return res,err
}

func DelHost(hid string)error{
	host:=new(Hosts)
	_,err:=MysqlDB.Table("ansible_host").Where("host_id=?",hid).Delete(host)
	if err!=nil{
		log.Error(err)
	}
	return err
}