package orm

import (
	"time"
	"encoding/hex"
	"crypto/md5"
	"io"

	log "github.com/astaxie/beego/logs"
)

type User struct {
	ID 			string 		`xorm:"user_id" json:"user_id"` 
	Account 	string 		`xorm:"user_account" json:"user_account"`
	Password 	string 		`xorm:"user_password" json:"user_password"`
	Created		time.Time  	`xorm:"created" json:"created"`
}

type UserList struct {
	ID 			string 		`xorm:"user_id" json:"user_id"` 
	Account 	string 		`xorm:"user_account" json:"user_account"`
	Password 	string 		`xorm:"user_password" json:"-"`
	Created		time.Time  	`xorm:"created" json:"created"`
}


func AuthUser(a,p string)(bool,string){
	h := md5.New()
	io.WriteString(h, p)
	passwdMD5 := hex.EncodeToString(h.Sum(nil))
	var user User
	res,err:=MysqlDB.Table("ansible_user").Where("user_account=? and user_password=?",a,passwdMD5).Get(&user)
	if err!=nil{
		log.Error(err)
	}
	return res,user.ID
}

func AddUser(user *User)error{
	h := md5.New()
	io.WriteString(h, user.Password)
	passwdMD5 := hex.EncodeToString(h.Sum(nil))
	user.Password=passwdMD5
	_,err:=MysqlDB.Table("ansible_user").Insert(user)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func UpdateUser(user *User)error{
	h := md5.New()
	io.WriteString(h, user.Password)
	passwdMD5 := hex.EncodeToString(h.Sum(nil))
	user.Password=passwdMD5
	_,err:=MysqlDB.Table("ansible_user").Where("user_id=?",user.ID).Update(user)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func FindUsers()(*[]UserList,error){
	var user []UserList
	err:=MysqlDB.Table("ansible_user").Find(&user)
	if err!=nil{
		log.Error(err)
	}
	return &user,err
}

func DelUser(uid string)error{
	user:=new(User)
	_,err:=MysqlDB.Table("ansible_user").Where("user_id=?",uid).Delete(user)
	if err!=nil{
		log.Error(err)
	}
	return err
}
func GetUser(uid string)(*UserList,error){
	var user UserList
	res,err:=MysqlDB.Table("ansible_user").Where("user_id=?",uid).Get(&user)
	if err!=nil{
		log.Error(err)
	}
	if !res{
		return nil,err
	}
	return &user,err
}