package template

import (
	"testing"
)

func Test_RefVars(t *testing.T){
	RefVars("/root/go/src/github.com/gzsunrun/ansible-manager/tools/amcreate/testdata/vars.yml")
}