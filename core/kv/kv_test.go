package kv_test

import (
	"testing"
	"time"
	"fmt"
	"github.com/gzsunrun/ansible-manager/core/kv"
)

func Test_SetKVClient(t *testing.T){
	err:=kv.SetKVClient()
	if err!=nil{
		t.Fatal(err)
	}
	time.Sleep(3*time.Second)
	fmt.Println("set:",kv.DefaultClient.GetStorage())
}

func Test_SetCall(t *testing.T){
	timer :=func(timer kv.Timer,put bool){
		fmt.Println(timer)
	}
	node :=func( node kv.Node,put bool){
		fmt.Println(node)
	}
	task :=func(task kv.Task,put bool){
		fmt.Println(task)
	}
	kv.DefaultClient.SetCall(timer,node,task)
}

func Test_AddConnect(t *testing.T){
	connect:=kv.Connect{
		ConnectID:"2",
		TaskID:"1234",
	}
	err:=kv.DefaultClient.AddConnect(connect)
	if err!=nil{
		t.Fatal(err)
	}
	time.Sleep(3*time.Second)
	fmt.Println("addconnect:",kv.DefaultClient.GetStorage())
}

func Test_AddTask(t *testing.T){
	task:=kv.Task{
		ID:"12345",
		NodeID:"10.21.1.1781",
	}
	err:=kv.DefaultClient.AddTask(task)
	if err!=nil{
		t.Fatal(err)
	}
	time.Sleep(3*time.Second)
	fmt.Println("addtask1:",kv.DefaultClient.GetStorage())
	task=kv.Task{
		ID:"12345",
		NodeID:"10.21.1.1782",
	}
	err=kv.DefaultClient.AddTask(task)
	if err!=nil{
		t.Fatal(err)
	}
	time.Sleep(3*time.Second)
	fmt.Println("addtask2:",kv.DefaultClient.GetStorage())
	task=kv.Task{
		ID:"12345",
		NodeID:"10.21.1.1783",
	}
	
	err=kv.DefaultClient.AddTask(task)
	if err!=nil{
		t.Fatal(err)
	}
	time.Sleep(3*time.Second)
	fmt.Println("addtask3:",kv.DefaultClient.GetStorage())
}

func Test_DeleteConnect(t *testing.T){
	err:=kv.DefaultClient.DeleteConnect("2")
	if err!=nil{
		t.Fatal(err)
	}
	time.Sleep(1*time.Second)
	fmt.Println("delconn:",kv.DefaultClient.GetStorage())
}

func Test_DeleteTask(t *testing.T){
	err:=kv.DefaultClient.DeleteTask("12345")
	if err!=nil{
		t.Fatal(err)
	}
	time.Sleep(1*time.Second)
	fmt.Println("deltask:",kv.DefaultClient.GetStorage())
}
