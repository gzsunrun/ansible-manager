package controllers

import (
	"strings"

	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/helm"
)

// HarborController repo controller
type HarborController struct {
	BaseController
}

// List get repo list
func (c *HarborController) List() {
	defer c.ServeJSON()
	data, err := helm.GetCharts()
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, data, 200)

}

func (c *HarborController) Sync() {
	defer c.ServeJSON()
	addr := c.GetString("addr")
	nameStr := c.GetString("names")
	names := strings.Split(nameStr, ",")

	if addr == "" {
		c.SetErrMsg(400, "地址不能为空")
		return
	}
	defer c.ServeJSON()
	data, err := helm.GetDstCharts(addr)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	repos := make([]helm.ChartVesion, 0)
	for _, charts := range data {
		if nameStr == "*" {
			for _, c := range charts.Cs {
				repo := helm.ChartVesion{
					Version: c.Version,
					Name:    charts.Name,
				}
				repos = append(repos, repo)
			}
		} else {
			for _, n := range names {
				if charts.Name == n {
					for _, c := range charts.Cs {
						repo := helm.ChartVesion{
							Version: c.Version,
							Name:    n,
						}
						repos = append(repos, repo)
					}
				}
			}
		}
	}

	err = helm.SyncCharts(addr, config.Cfg.Harbor.URL+"/chartrepo/"+config.Cfg.Harbor.Repo, repos)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
}
