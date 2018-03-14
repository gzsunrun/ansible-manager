package role_test

import (
	"strconv"
	"testing"
	"github.com/gzsunrun/ansible-manager/core/master"
	"github.com/gzsunrun/ansible-manager/core/kv"
)

func Test_Scheduler(t *testing.T){
	n:=make(map[string]kv.Node)
	task:=make(map[string]kv.Task)
	for i:=0;i<4;i++{
		n[strconv.Itoa(i)]=kv.Node{
			IP:strconv.Itoa(i),
		}
	}
	for i:=0;i<20;i++{
		if i<2{
			task[strconv.Itoa(i)]=kv.Task{
				ID:strconv.Itoa(i),
				NodeID:"0",
			}
			continue
		}
		if i<7{
			task[strconv.Itoa(i)]=kv.Task{
				ID:strconv.Itoa(i),
				NodeID:"1",
			}
			continue
		}
		if i<15{
			task[strconv.Itoa(i)]=kv.Task{
				ID:strconv.Itoa(i),
				NodeID:"4",
			}
			continue
		}
		if i<20{
			task[strconv.Itoa(i)]=kv.Task{
				ID:strconv.Itoa(i),
				NodeID:"3",
			}
			continue
		}
	}

	if master.Scheduler(n,task)!="0"{
		t.Error("res is error")
	}
}