package controllers

import (
	"encoding/json"

	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/satori/go.uuid"
)

// ProjectController controller of project.
type ProjectController struct {
	BaseController
}

// Create create a project.
func (c *ProjectController) Create() {
	defer c.ServeJSON()
	project := orm.Project{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &project); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	project.ID = uuid.Must(uuid.NewV4()).String()
	project.UserID = c.GetUid()
	err := orm.CreateProject(&project)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

// ProjectAndHosts a struct that include project and host.
type ProjectAndHosts struct {
	Project      orm.Project       `json:"project"`
	ProjectHosts []orm.ProjectHost `json:"project_hosts"`
}

// CreateAndAddHosts create  or update a project and add hosts.
func (c *ProjectController) CreateAndAddHosts() {
	defer c.ServeJSON()
	projectHosts := ProjectAndHosts{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &projectHosts); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if projectHosts.Project.ID == "" {
		projectHosts.Project.ID = uuid.Must(uuid.NewV4()).String()
		projectHosts.Project.UserID = c.GetUid()
		err := orm.CreateProject(&projectHosts.Project)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
	} else {
		err := orm.UPdateProject(&projectHosts.Project)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
		err = orm.DelAllHostsByPid(projectHosts.Project.ID)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
	}
	for i := range projectHosts.ProjectHosts {
		projectHosts.ProjectHosts[i].ProjectID = projectHosts.Project.ID
	}
	err := orm.AddHostToProject(&projectHosts.ProjectHosts)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

// GetProject return the project info with project uuid.
func (c *ProjectController) GetProject() {
	defer c.ServeJSON()
	pid := c.GetString("project_id")
	project := orm.Project{}
	_, err := orm.GetProject(pid, &project)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	phs, err := orm.FindProjectHost(pid)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	data := ProjectAndHosts{
		Project:      project,
		ProjectHosts: *phs,
	}
	c.SetResult(err, data, 200)
}

// List return list of project by user uuid.
func (c *ProjectController) List() {
	defer c.ServeJSON()
	projects := []orm.Project{}
	uid := c.GetUid()
	err := orm.FindProject(uid, &projects)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, projects, 200)
}

// Del delete project by project uuid.
func (c *ProjectController) Del() {
	defer c.ServeJSON()
	pid := c.GetString("project_id")
	err := orm.DelProject(pid)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

// AddHost add a host to project
func (c *ProjectController) AddHost() {
	defer c.ServeJSON()
	projectHost := []orm.ProjectHost{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &projectHost); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	err := orm.AddHostToProject(&projectHost)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

// DelHost delete host from project
func (c ProjectController) DelHost() {
	defer c.ServeJSON()
	projectHost := orm.ProjectHost{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &projectHost); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	err := orm.DelHostFormProject(&projectHost)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, nil, 204)
}

// HostList get all hosts in project by project uuid.
func (c *ProjectController) HostList() {
	defer c.ServeJSON()
	pid := c.GetString("project_id")
	var hosts []orm.HostsList
	err := orm.FindHostFromProject(pid, &hosts)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil, hosts, 200)
}
