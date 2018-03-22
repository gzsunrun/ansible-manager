package role

import (
	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/orm"
)

func MasterSet() {
	node = func(n kv.Node, p bool) {
		if !p {
			for _, t := range kv.DefaultClient.GetStorage().Tasks {
				if t.NodeID == n.ID() {
					log.Error("node:", n.ID(), "timeout")
					CommitTask(t)
				}
			}
		} else {
			log.Info("node:", n.ID(), "Add")
		}

	}

	sche = func(t kv.Task, p bool) {
		log.Info("schedule task:", t.ID)
		if p {
			err := kv.DefaultClient.DeleteTask(t.ID)
			if err != nil {
				log.Error(err)
			}
			CommitTask(t)
		}
	}
	CleanTask()
}

func Scheduler(n map[string]kv.Node, t map[string]kv.Task) string {
	if len(n) == 0 {
		return ""
	}
	if len(t) == 0 {
		for k, node := range n {
			if node.Worker {
				return k
			}
		}
	}
	counts := make(map[string]int)
	nodesCount := len(n)
	nodes := make([]string, nodesCount)
	for _, v := range t {
		counts[v.NodeID]++
	}

	for k := range counts {
		flag := true
		for _, node := range n {
			if k == node.NodeID && node.Worker {
				flag = false
			}
			if _, ok := counts[node.NodeID]; !ok {
				counts[node.NodeID] = 0
			}
		}
		if flag {
			delete(counts, k)
		}
	}

	log.Info(counts)
	nc := len(counts)
	if nc == 0 {
		for k := range n {
			return k
		}
	}
	for i := 0; i < nc; i++ {
		for k, v := range counts {
			if i > 0 && k == nodes[i-1] {
				continue
			}
			if nodes[i] == "" {
				nodes[i] = k
				continue
			}
			if v < counts[nodes[i]] {
				nodes[i] = k
			}
		}
		delete(counts, nodes[i])
	}

	return nodes[0]
}

func CleanTask() {
	for taskID, t := range kv.DefaultClient.GetStorage().Tasks {
		flag := true
		for nid := range kv.DefaultClient.GetStorage().Nodes {
			if t.NodeID == nid {
				flag = false
			}
		}
		if flag {
			if t.Timer {
				log.Info("clean timer:", taskID)
				ormTimer := new(orm.Timer)
				ormTimer.ID = taskID
				ormTimer.Status = false
				err := orm.UpdateTimerStatus(ormTimer)
				if err != nil {
					log.Error(err)
				}
			} else {
				log.Info("clean task:", taskID)
				ormTask := new(orm.Task)
				ormTask.ID = taskID
				ormTask.Status = "error"
				err := orm.UpdateTask(ormTask)
				if err != nil {
					log.Error(err)
				}
			}

			kv.DefaultClient.DeleteTask(taskID)
		}
	}
}

func CommitTask(task kv.Task) {
	nodeID := Scheduler(kv.DefaultClient.GetStorage().Nodes, kv.DefaultClient.GetStorage().Tasks)
	if nodeID == "" {
		log.Error("no worker node")
		return
	}
	task.NodeID = nodeID
	log.Info("task:", task.ID, "->node:", nodeID)
	err := kv.DefaultClient.AddTask(task)
	if err != nil {
		log.Error(err)
	}
}
