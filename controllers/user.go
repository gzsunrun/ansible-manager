package controllers


import (
	"encoding/json"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/auth"
	//log "github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
)

type UserController struct {
	BaseController
}

func (c *UserController)Login(){
	defer c.ServeJSON()
	user:=c.GetString("account")
	psw:=c.GetString("password")
	res,uid:=orm.AuthUser(user,psw)
	if !res{
		c.SetResult(nil,nil,403)
		return
	}
	token,err:=auth.IssueTokenUsingDgrijalva(uid,nil)
	if err!=nil{
		c.SetResult(err,nil,403)
		return
	}
	c.SetResult(nil,token,200,"token")
}

func (c *UserController)Info(){
	defer c.ServeJSON()
	user,err:=orm.GetUser(c.GetUid())
	if err!=nil{
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil,user,200)
}

func (c *UserController)Create(){
	defer c.ServeJSON()
	user:=orm.User{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &user); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if user.ID!=""{
		err:=orm.UpdateUser(&user)
		if err!=nil{
			c.SetResult(err, nil, 400)
			return
		}
	}else{
		user.ID=uuid.Must(uuid.NewV4()).String()
		err:=orm.AddUser(&user)
		if err!=nil{
			c.SetResult(err, nil, 400)
			return
		}
	}
	c.SetResult(nil, nil, 204)
}

func (c *UserController)Del(){
	defer c.ServeJSON()
	uid:=c.GetString("uid")
	err:=orm.DelUser(uid)
	if err!=nil{
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

func (c *UserController)List(){
	defer c.ServeJSON()
	users,err:=orm.FindUsers()
	if err!=nil{
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, users, 200)
}