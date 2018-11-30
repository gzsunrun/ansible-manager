package helm_test

import (
	"testing"

	"github.com/gzsunrun/ansible-manager/core/helm"
)

func init() {
	helm.InitHarbor("http://192.168.1.100:5000", "helm-sunrun-charts")
}

func Test_GetCharts(t *testing.T) {
	res, err := helm.GetCharts()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}
