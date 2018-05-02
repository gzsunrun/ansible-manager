package kv

import (
	"sync"

	"github.com/gzsunrun/ansible-manager/core/config"
)

// KVClient kv client
type KVClient interface {
	LocalNode() *Node
	AddTask(task Task) error
	DelAllTask() error
	GetStorage() *LocalMenory
	AddScheduler(task Task) error
	DeleteTask(tid string) error
	SetCall(func(Node, bool), func(Task, bool), func(Task, bool))
}

// LocalMenory local data
type LocalMenory struct {
	Tasks map[string]Task
	Nodes map[string]Node
	Lock  *sync.Mutex
}

// Task task struct
type Task struct {
	Timer  bool
	ID     string
	NodeID string
}

// DefaultClient common kv client
var DefaultClient KVClient

// SetKVClient init common kv client
func SetKVClient() error {
	var err error
	ep := config.Cfg.Etcd.Endpoints
	port := config.Cfg.Common.Port
	worker := config.Cfg.Common.Worker
	master := config.Cfg.Common.Master
	timeout := config.Cfg.Common.Timeout
	if config.Cfg.Etcd.Enable {
		c, err := NewEtcd(ep, timeout, port, "/api/ansible/ws", worker, master)
		if err != nil {
			return err
		}
		DefaultClient = c
		go c.KeepNode()
		c.RegNode()
		err = c.SyncData()
		if err != nil {
			return err
		}
	} else {
		c, err := NewLocalKV(port, "/api/ansible/ws", worker, master)
		if err != nil {
			return err
		}
		DefaultClient = c
	}

	return err
}
