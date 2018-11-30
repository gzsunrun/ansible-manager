package controllers

import (
	"encoding/json"
	"strings"

	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/helm"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/hashwing/log"
)

// HelmController helm controller
type HelmController struct {
	BaseController
}

// List get helm list
func (c *HelmController) List() {
	defer c.ServeJSON()
	pid := c.GetString("project_id")
	if pid == "" {
		c.SetErrMsg(400, "project_id 不能为空")
		return
	}
	var hosts []orm.HostsList
	err := orm.FindHostFromProject(pid, &hosts)
	if err != nil {
		c.SetErrMsg(500, "获取主机信息错误")
		return
	}
	var host orm.HostsList
	for i, h := range hosts {
		if strings.Contains(h.Alias, "@helm") {
			host = h
			break
		}
		if i == len(hosts)-1 {
			c.SetErrMsg(400, "所选集群没有可用的主机")
			return
		}
	}

	//config.Cfg.Harbor.URL
	r, err := helm.NewRunner(host.IP, host.User, host.Password, host.Key, host.Port)
	if err != nil {
		c.SetErrMsg(500, "连接主机错误")
		return
	}
	data, err := r.HelmList()
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, data, 200)

}

type installInput struct {
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
	Chart     string `json:"chart_name"`
	Version   string `json:"chart_version"`
	NameSpace string `json:"namespace"`
	Values    string `json:"values"`
	Update    bool   `json:"update"`
}

// Install install helm
func (c *HelmController) Install() {
	defer c.ServeJSON()
	iopt := installInput{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &iopt); err != nil {
		log.Error(err)
		c.SetResult(err, nil, 400)
		return
	}

	var hosts []orm.HostsList
	err := orm.FindHostFromProject(iopt.ProjectID, &hosts)
	if err != nil {
		c.SetErrMsg(500, "获取主机信息错误")
		return
	}
	var host orm.HostsList
	for i, h := range hosts {
		if strings.Contains(h.Alias, "@helm") {
			host = h
			break
		}
		if i == len(hosts)-1 {
			c.SetErrMsg(400, "所选集群没有可用的主机")
			return
		}
	}

	//config.Cfg.Harbor.URL
	r, err := helm.NewRunner(host.IP, host.User, host.Password, host.Key, host.Port)
	if err != nil {
		c.SetErrMsg(500, "连接主机错误")
		return
	}

	if !iopt.Update {
		data, err := r.HelmList()
		if err != nil {
			c.SetErrMsg(500, err.Error())
			return
		}
		for _, d := range data {
			if iopt.Name == d.Name {
				c.SetErrMsg(500, "发布名称已经存在")
				return
			}
		}
	}

	err = r.Install(iopt.Name, config.Cfg.Harbor.URL+"/chartrepo/"+config.Cfg.Harbor.Repo, iopt.Chart, iopt.Version, iopt.NameSpace, iopt.Values, iopt.Update)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, nil, 204)
}

type deleteOpt struct {
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
	Purge     bool   `json:"purge"`
}

func (c *HelmController) HelmDelete() {
	defer c.ServeJSON()
	dopt := deleteOpt{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dopt); err != nil {
		log.Error(err)
		c.SetResult(err, nil, 400)
		return
	}
	var hosts []orm.HostsList
	err := orm.FindHostFromProject(dopt.ProjectID, &hosts)
	if err != nil {
		c.SetErrMsg(500, "获取主机信息错误")
		return
	}
	var host orm.HostsList
	for i, h := range hosts {
		if strings.Contains(h.Alias, "@helm") {
			host = h
			break
		}
		if i == len(hosts)-1 {
			c.SetErrMsg(400, "所选集群没有可用的主机")
			return
		}
	}
	r, err := helm.NewRunner(host.IP, host.User, host.Password, host.Key, host.Port)
	if err != nil {
		c.SetErrMsg(500, "连接主机错误")
		return
	}
	err = r.HelmDelete(dopt.Name, dopt.Purge)
	if err != nil {
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, nil, 204)
}

// GetValues get values
func (c *HelmController) GetValues() {
	defer c.ServeJSON()
	name := c.GetString("chart_name")
	version := c.GetString("chart_version")
	if name == "" || version == "" {
		c.SetErrMsg(400, "请求参数错误")
		return
	}

	data, err := helm.GetValues(name, version)
	if err != nil {
		log.Error(err)
		c.SetResult(err, nil, 500)
		return
	}
	c.Ctx.Output.ContentType("application/json")
	c.Ctx.Output.Body(data)
}

// GetHistoryValues get values
func (c *HelmController) GetHistoryValues() {
	defer c.ServeJSON()
	name := c.GetString("name")
	projectID := c.GetString("project_id")
	if name == "" || projectID == "" {
		c.SetErrMsg(400, "请求参数错误")
		return
	}
	var hosts []orm.HostsList
	err := orm.FindHostFromProject(projectID, &hosts)
	if err != nil {
		c.SetErrMsg(500, "获取主机信息错误")
		return
	}
	var host orm.HostsList
	for i, h := range hosts {
		if strings.Contains(h.Alias, "@helm") {
			host = h
			break
		}
		if i == len(hosts)-1 {
			c.SetErrMsg(400, "所选集群没有可用的主机")
			return
		}
	}
	r, err := helm.NewRunner(host.IP, host.User, host.Password, host.Key, host.Port)
	if err != nil {
		c.SetErrMsg(500, "连接主机错误")
		return
	}

	data, err := r.GetValues(name)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, data, 200, "values")
}

// ReleaseStatus get status
func (c *HelmController) ReleaseStatus() {
	defer c.ServeJSON()
	name := c.GetString("name")
	projectID := c.GetString("project_id")
	if name == "" || projectID == "" {
		c.SetErrMsg(400, "请求参数错误")
		return
	}
	var hosts []orm.HostsList
	err := orm.FindHostFromProject(projectID, &hosts)
	if err != nil {
		c.SetErrMsg(500, "获取主机信息错误")
		return
	}
	var host orm.HostsList
	for i, h := range hosts {
		if strings.Contains(h.Alias, "@helm") {
			host = h
			break
		}
		if i == len(hosts)-1 {
			c.SetErrMsg(400, "所选集群没有可用的主机")
			return
		}
	}
	r, err := helm.NewRunner(host.IP, host.User, host.Password, host.Key, host.Port)
	if err != nil {
		c.SetErrMsg(500, "连接主机错误")
		return
	}

	data, err := r.HelmStatus(name)
	if err != nil {
		log.Error(err)
		c.SetErrMsg(500, err.Error())
		return
	}
	c.SetResult(nil, data, 200)
}
