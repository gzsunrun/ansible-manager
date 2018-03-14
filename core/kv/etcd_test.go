package kv

import (
	"testing"
	"time"
	"fmt"

	//"github.com/coreos/etcd/clientv3"
)

func Test_NewEtcd(t *testing.T){
	_,err:=NewEtcd([]string{"http://127.0.0.1:2379"},5*time.Second)
	if err!=nil{
		t.Fatal(err)
	}
}

func Test_RegNode(t *testing.T){
	c,err:=NewEtcd([]string{"http://127.0.0.1:2379"},5*time.Second)
	if err!=nil{
		t.Error(err)
	}
	// _,err=c.Delete("ansible/",clientv3.WithPrefix())
	// if err!=nil{
	// 	t.Error(err)
	// }
	n,err:=NewNode(5)
	if err!=nil{
		t.Error(err)
	}
	c.Node=n
	err=c.RegNode()
	if err!=nil{
		t.Error(err)
	}
}

func Test_FindNodes(t *testing.T){
	c,err:=NewEtcd([]string{"http://127.0.0.1:2379"},5*time.Second)
	if err!=nil{
		t.Error(err)
	}
	var nodes []Node
	err =c.FindNodes(&nodes)
	if err!=nil{
		t.Error(err)
	}
	fmt.Println(nodes)
	for _,v:=range nodes{
		var node Node
		err =c.GetNode(&node,v.ID())
		if err!=nil{
			t.Error(err)
		}
		fmt.Println(node)
			
	}
}


func Test_WatchNode(t *testing.T){
	c,err:=NewEtcd([]string{"http://127.0.0.1:2379"},5*time.Second)
	if err!=nil{
		t.Error(err)
	}
	pre:=c.Storage.Nodes
	go c.WatchData()
	err=c.SyncData()
	if err!=nil{
		t.Error(err)
	}
	time.Sleep(7*time.Second)
	now:=c.Storage.Nodes
	if len(pre)==len(now){
		t.Error("watchdata err")
	}
	fmt.Println(pre,now)
}

