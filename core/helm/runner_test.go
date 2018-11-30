package helm_test

import (
	"testing"

	"github.com/gzsunrun/ansible-manager/core/helm"
)

func Test_HelmList(t *testing.T) {
	r, err := helm.NewRunner("10.21.21.45", "root", "sunrunvas", "", "22")
	if err != nil {
		t.Error(err)
		return
	}
	//r.HelmList()
	//r.Install("am-test", "http://192.168.1.101:5000/chartrepo/helm-sunrun-charts", "grafana", "1.12.8", "hashwing", "")
	r.HelmStatus("prometheus")

}
