package controllers

import (
	"sync"
	"encoding/json"

	"github.com/satori/go.uuid"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/function"
)

type HostController struct{
	BaseController
}

func (c *HostController)List(){
	defer c.ServeJSON()
	uid:=c.GetUid()
	var hosts []orm.HostsList
	err:=orm.FindHosts(uid,&hosts)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	var wg sync.WaitGroup
	var l sync.Mutex
	for i,_:=range hosts{
		wg.Add(1)
		go func(i int,hosts *[]orm.HostsList){
			defer wg.Done()
			status:=function.SshDail((*hosts)[i])
			l.Lock()
			if status=="success"{
				(*hosts)[i].Status=true
			}else{
				(*hosts)[i].Status=false
			}
			l.Unlock()
		}(i,&hosts)
	}
	wg.Wait()
	c.SetResult(nil,hosts,200)
}
func (c *HostController)ListNO(){
	defer c.ServeJSON()
	uid:=c.GetUid()
	var hosts []orm.HostsList
	err:=orm.FindHosts(uid,&hosts)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	c.SetResult(nil,hosts,200)
}

func (c *HostController)Get(){
	defer c.ServeJSON()
	hid:=c.GetString("host_id")
	var host orm.HostsList
	_,err:=orm.GetHost(hid,&host)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	c.SetResult(nil,host,200)
}

func (c *HostController)Create(){
	defer c.ServeJSON()
	host:=orm.Hosts{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &host); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if host.ID !=""{
		var err error
		err=orm.UPdateHost(&host)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
		c.SetResult(nil,nil,204)
		return
	}
	host.ID=uuid.Must(uuid.NewV4()).String()
	host.UserID=c.GetUid()
	err:=orm.CreateHost(&host)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil,nil,204)
	defer c.ServeJSON()
}

func (c *HostController)Del(){
	defer c.ServeJSON()
	id:=c.GetString("host_id")
	err:=orm.DelHost(id)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil,nil,204)
}