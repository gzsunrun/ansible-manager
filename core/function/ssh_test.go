package function

import (
	"fmt"
	"testing"

	"github.com/gzsunrun/ansible-manager/core/orm"
)

func Test_SshDail(t *testing.T){
	h:=orm.HostsList{
		IP:"10.21.21.179",
		User:"root",
		Password:"sunrunvas",
	}
	fmt.Println(SshDail(h))
}