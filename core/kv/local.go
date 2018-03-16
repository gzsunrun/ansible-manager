package kv

import (
	"sync"
)

// type LocalClient struct{
// 	Storage []map[string]string
// }

type LocalKV struct{
	Node		*Node
	Storage 	*LocalMenory
	nodeCall	func(Node,bool)	
	taskCall	func(Task,bool)
	scheCall	func(Task,bool) 
}

func NewLocalKV(port int,path string,worker,master bool)(*LocalKV,error){
	node,err:=NewNode(0,port,path,worker,master)
	if err!=nil{
		return nil,err
	}
	storage:=&LocalMenory{
		Tasks:make(map[string]Task),
		Nodes:make(map[string]Node),
		Lock:new(sync.Mutex),
	}
	storage.Nodes[node.NodeID]=*node
	return &LocalKV{
		Node:node,
		Storage:storage,
	},nil
}

// get this node info
func (lkv *LocalKV)LocalNode()*Node{
	return lkv.Node
}

//delete all tasks
func (lkv *LocalKV)DelAllTask()error{
	lkv.Storage.Tasks=make(map[string]Task)
	return nil
}

func (lkv *LocalKV)AddTask(task Task)error{
	lkv.Storage.Lock.Lock()
	lkv.Storage.Tasks[task.ID]=task
	lkv.Storage.Lock.Unlock()
	lkv.taskCall(task,true)
	return nil
}

func (lkv *LocalKV)DeleteTask(tid string)error{
	lkv.taskCall(lkv.Storage.Tasks[tid],false)
	lkv.Storage.Lock.Lock()
	delete(lkv.Storage.Tasks,tid)
	lkv.Storage.Lock.Unlock()
	return nil
}

func (lkv *LocalKV)GetStorage()*LocalMenory{
	return lkv.Storage
}

func (lkv *LocalKV)AddScheduler(task Task)error{
	lkv.scheCall(task,true)
	return nil
}

func (lkv *LocalKV)DeleteScheduler(tid string)error{
	return nil
}

func (lkv *LocalKV)SetCall(node func(Node,bool),task,sche func(Task,bool)){
	lkv.taskCall=task
	lkv.nodeCall=node
	lkv.scheCall=sche
}

