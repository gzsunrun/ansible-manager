package controllers

import (
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
