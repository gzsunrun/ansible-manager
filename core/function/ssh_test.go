package function

import (
	"fmt"
	"testing"

	"github.com/gzsunrun/ansible-manager/core/orm"
)

func Test_SshDail(t *testing.T){
	h:=orm.HostsList{
		IP:"10.21.1.161",
		Password:"sunrunvas",
	}
	fmt.Println(SshDail(h))
}