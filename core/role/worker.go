package role

import (
	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/tasks"
)

// WorkerSet worker config
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
			log.Debug("LocalNodeID:",kv.DefaultClient.LocalNode().ID(),"TaskNodeTD:",t.NodeID)
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
