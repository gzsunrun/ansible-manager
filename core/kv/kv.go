package kv

import (
	"sync"

	"github.com/gzsunrun/ansible-manager/core/config"
)

type KVClient interface{
	LocalNode()*Node
	AddTask(task Task)error
	DelAllTask()error
	GetStorage()*LocalMenory
	AddScheduler(task Task)error
	DeleteTask(tid string)error
	SetCall(func(Node,bool),func(Task,bool),func(Task,bool))
}


type LocalMenory struct{
	Tasks		map[string]Task
	Nodes		map[string]Node
	Lock		*sync.Mutex
}


type Task struct{
	Timer 	bool
	ID 		string
	NodeID	string
}




var DefaultClient KVClient


func SetKVClient()error{
	var err error
	ep:=config.Cfg.Etcd.Endpoints
	port:=config.Cfg.Common.Port
	worker:=config.Cfg.Common.Worker
	master:=config.Cfg.Common.Master
	timeout:=config.Cfg.Common.Timeout
	c,err:=NewEtcd(ep,timeout,port,"/api/ansible/ws",worker,master)
	if err!=nil{
		return err
	}
	DefaultClient=c
	go c.KeepNode()
	c.RegNode()
	err=c.SyncData()
	if err!=nil{
		return err
	}
	return err
}