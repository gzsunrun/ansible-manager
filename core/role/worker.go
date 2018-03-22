package role

import (
	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/tasks"
)

func WorkerSet() {
	task = func(t kv.Task, p bool) {
		if p {
			if kv.DefaultClient.LocalNode().ID() == t.NodeID {
				if !t.Timer {
					ormTask := new(orm.Task)
					ormTask.ID = t.ID
					ormTask.Status = "waiting"
					err := orm.UpdateTask(ormTask)
					if err != nil {
						log.Error(err)
					}
					tasks.AddTask(t.ID)
				} else {
					go tasks.SetTimer(t.ID)
				}
			} else {
				if !t.Timer {
					tasks.StopTask(t.ID)
				} else {
					go tasks.StopTimer(t.ID)
				}
			}
		} else {
			if kv.DefaultClient.LocalNode().ID() == t.NodeID {
				if !t.Timer {
					tasks.StopTask(t.ID)
				} else {
					go tasks.StopTimer(t.ID)
				}
			}
		}

	}
}
