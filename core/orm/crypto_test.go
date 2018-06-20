package orm_test

import (
	"fmt"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"testing"
)

func Test_RsaEncrypt(t *testing.T) {
	a:=`12333fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff`
	data, err := orm.RsaEncrypt([]byte(a))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
	origData, err := orm.RsaDecrypt(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(origData))
}
