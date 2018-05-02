package kv

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)
func Test_New(t *testing.T ){
	var index uint64
	for{
		client, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			t.Error(err)
		}
		opts := api.QueryOptions{
			WaitIndex:index,
		}
		kp, meta, err:=client.KV().List("a",&opts)
		if err != nil {
			t.Error(err)
		}
		if len(kp)>0{
			fmt.Println(kp[0].Key)
		}
		index=meta.LastIndex
		
		time.Sleep(time.Second*3)
	}
	
}