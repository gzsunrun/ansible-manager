package orm_test

import (
	"fmt"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"testing"
)

func Test_RsaEncrypt(t *testing.T) {
	data, err := orm.RsaEncrypt([]byte("123456"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
	origData, err := orm.RsaDecrypt(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(origData))
}
