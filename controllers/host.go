package controllers

import (
	"encoding/json"
	"sync"

	"github.com/gzsunrun/ansible-manager/core/function"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/satori/go.uuid"
)

//HostController host controller
type HostController struct {
	BaseController
}

// List get host list with status
func (c *HostController) List() {
	defer c.ServeJSON()
	uid := c.GetUid()
	var hosts []orm.HostsList
	err := orm.FindHosts(uid, &hosts)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	var wg sync.WaitGroup
	var l sync.Mutex
	for i:= range hosts {
		wg.Add(1)
		go func(i int, hosts *[]orm.HostsList) {
			defer wg.Done()
			status := function.SshDail((*hosts)[i])
			l.Lock()
			if status == "success" {
				(*hosts)[i].Status = true
			} else {
				(*hosts)[i].Status = false
			}
			l.Unlock()
		}(i, &hosts)
	}
	wg.Wait()
	c.SetResult(nil, hosts, 200)
}

// ListNO get host list without status
func (c *HostController) ListNO() {
	defer c.ServeJSON()
	uid := c.GetUid()
	var hosts []orm.HostsList
	err := orm.FindHosts(uid, &hosts)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, hosts, 200)
}

// Get get host by host uuid
func (c *HostController) Get() {
	defer c.ServeJSON()
	hid := c.GetString("host_id")
	var host orm.HostsList
	_, err := orm.GetHost(hid, &host)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, host, 200)
}

// Create create  or update a host
func (c *HostController) Create() {
	defer c.ServeJSON()
	host := orm.Hosts{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &host); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if host.ID != "" {
		var err error
		err = orm.UPdateHost(&host)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
		c.SetResult(nil, nil, 204)
		return
	}
	host.ID = uuid.Must(uuid.NewV4()).String()
	host.UserID = c.GetUid()
	err := orm.CreateHost(&host)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
	defer c.ServeJSON()
}

// Del delete a host
func (c *HostController) Del() {
	defer c.ServeJSON()
	id := c.GetString("host_id")
	err := orm.DelHost(id)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}
