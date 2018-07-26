package controllers

import (
	"encoding/json"
	"time"

	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/satori/go.uuid"
)

// TimerController timer controller
type TimerController struct {
	BaseController
}

// List find timer list
func (c *TimerController) List() {
	defer c.ServeJSON()
	uid := c.GetUid()
	timers, err := orm.FindTimers(uid)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	for i, v := range *timers {
		if v.Status {
			(*timers)[i].Surplus = v.Interval - (int(time.Now().Unix()) - v.Start)
		} else {
			(*timers)[i].Surplus = -999
		}
	}
	c.SetResult(nil, timers, 200)
}

// Create create  or update a timer
func (c *TimerController) Create() {
	defer c.ServeJSON()
	timer := orm.Timer{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &timer); err != nil {
		log.Error(err)
		c.SetResult(err, nil, 400)
		return
	}
	timer.UserID = c.GetUid()
	if timer.ID != "" {
		err := orm.UpdateTimer(&timer)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
		c.SetResult(nil, nil, 204)
		return
	}
	timer.ID = uuid.NewV4().String()
	err := orm.CreateTimer(&timer)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if timer.Status {
		err := kv.DefaultClient.AddScheduler(kv.Task{
			ID:    timer.ID,
			Timer: true,
		})
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
	}
	c.SetResult(nil, nil, 204)
}

// Start start a timer
func (c *TimerController) Start() {
	defer c.ServeJSON()
	tid := c.GetString("timer_id")
	if tid == "" {
		c.SetResult(nil, nil, 400)
		return
	}
	master := false
	worker := false
	for _, node := range kv.DefaultClient.GetStorage().Nodes {
		if node.Master {
			master = true
		}
		if node.Worker {
			worker = true
		}
	}
	if master && worker {
		err := kv.DefaultClient.AddScheduler(kv.Task{
			ID:    tid,
			Timer: true,
		})
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
		c.SetResult(nil, nil, 204)
	} else {
		c.SetResult(nil, nil, 400)
	}
}

// Stop stop a timer
func (c *TimerController) Stop() {
	defer c.ServeJSON()
	tid := c.GetString("timer_id")
	err := kv.DefaultClient.DeleteTask(tid)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

// Del delete a timer
func (c *TimerController) Del() {
	defer c.ServeJSON()
	tid := c.GetString("timer_id")
	err := kv.DefaultClient.DeleteTask(tid)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	err = orm.DelTimer(tid)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

// Get get a timer info
func (c *TimerController) Get() {
	defer c.ServeJSON()
	tid := c.GetString("timer_id")
	res, timer, err := orm.GetTimer(tid)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if !res {
		c.SetResult(nil, nil, 204)
		return
	}
	c.SetResult(nil, timer, 200)
}
