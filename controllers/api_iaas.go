package controllers

import (
	"sync"
	"encoding/json"

	"github.com/satori/go.uuid"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/iaas"
	"github.com/gzsunrun/ansible-manager/core/tasks"
	"github.com/gzsunrun/ansible-manager/core/function"
)

type IaaSController struct{
	BaseController
}


type TaskJSON struct {
	ID            string    			`json:"uuid"`
	RepoName      string      			`json:"repo_name"`
	RepoID        string       			`json:"repo_id"`
	Tag           string    			`json:"playbook_tag"`
	PlaybookParse PlaybookParse     	`json:"playbook_parse"`
}


type PlaybookParse struct {
	Group []orm.Group `json:"group"`
	Vars  []orm.Vars  `json:"vars"`
}


func (c *IaaSController) ProjectHosts(){
	defer c.ServeJSON()
	projectID:=c.GetString("uuid")
	project:=orm.Project{}
	res,err:=orm.GetProject(projectID,&project)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if !res{
		newProject:=orm.Project{
			ID:projectID,
			UserID:c.GetUid(),
			Name:projectID,
		}
		err:=orm.CreateProject(&newProject)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
	}

	hosts,err:=iaas.GetProjectHosts(projectID)
	if err!=nil{
		c.SetResult(err, nil, 400)
		return
	}
	err=orm.DelAllHostsByPid(projectID)
	if err!=nil{
		c.SetResult(err, nil, 400)
		return
	}
	var wg sync.WaitGroup
	var l sync.Mutex
	for i,v:=range *hosts{
		host:=orm.HostsList{}
		res,err:=orm.GetHost(v.ID,&host)
		if err!=nil{
			c.SetResult(err, nil, 400)
			return
		}
		wg.Add(1)
		if !res{
			v.UserID=c.GetUid()
			err=orm.CreateHostList(&v)
			if err!=nil{
				c.SetResult(err, nil, 400)
				return
			}
			go func(i int,h orm.HostsList){
				defer wg.Done()
				status:=function.SshDail(h)
				l.Lock()
				if status=="success"{
					(*hosts)[i].Status=true
				}else{
					(*hosts)[i].Status=false
				}
				l.Unlock()
			}(i,v)
		}else{
			go func(i int,h orm.HostsList){
				defer wg.Done()
				status:=function.SshDail(h)
				l.Lock()
				if status=="success"{
					(*hosts)[i].Status=true
				}else{
					(*hosts)[i].Status=false
				}
				l.Unlock()
			}(i,host)
		}
		ph:=orm.ProjectHost{
			ProjectID:projectID,
			HostID:v.ID,
		}
		phs:=[]orm.ProjectHost{ph}
		err=orm.AddHostToProject(&phs)
		if err!=nil{
			c.SetResult(err, nil, 400)
			return
		}
	}
	wg.Wait()
	c.SetResult(nil,hosts,200)
}

func (c *IaaSController)CreateHost(){
	defer c.ServeJSON()
	host:=orm.Hosts{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &host); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if host.ID !=""{
		var err error
		if host.Password!=""||host.Key!=""{
			err=orm.UPdateAuthHost(&host)
		}else{
			err=orm.UPdateHost(&host)
		}
		if err != nil {
			msg:=ErrorMsg{
				Code:400,
				Message:"服务错误，更新失败",
			}
			c.SetResult(err, msg, 400)
			return
		}
	}else{
		host.ID=uuid.NewV4().String()
		host.UserID=c.GetUid()
		err:=orm.CreateHost(&host)
		if err != nil {
			msg:=ErrorMsg{
				Code:400,
				Message:"服务错误，插入失败",
			}
			c.SetResult(err, msg, 400)
			return
		}
	}
	var hostL orm.HostsList
	res,err:=orm.GetHost(host.ID,&hostL)
	if err != nil||!res {
		msg:=ErrorMsg{
			Code:400,
			Message:"服务错误，获取状态失败",
		}
		c.SetResult(err, msg, 400)
		return
	}
	status:=function.SshDail(hostL)
	if status=="fail"{
		msg:=ErrorMsg{
			Code:400,
			Message:"主机网络异常，请确认主机是否开机或网络是否正常",
		}
		c.SetResult(nil,msg,400)
		return
	}
	if status=="auth"{
		msg:=ErrorMsg{
			Code:400,
			Message:"主机认证失败，请确认主机用户密码或私钥是否正确",
		}
		c.SetResult(nil,msg,400)
		return
	}
	
	c.SetResult(nil,nil,204)
}

func (c *IaaSController) RepoVars() {
	defer c.ServeJSON()
	rid:=c.GetString("repo_id")
	var repo orm.Repository
	err:=orm.GetRepoByID(rid,&repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	group:=map[string]interface{}{
		"vars_value":repo.Group,
	}
	tag:=map[string]interface{}{
		"vars_value":repo.Tag,
	}
	data :=map[string]interface{}{
		"vars": repo.Vars,
		"group": group,
		"tag": tag,
	}
	c.SetResult(nil,data,200)
}
func (c *IaaSController)Stop(){
	defer c.ServeJSON()
	tid:=c.GetString("uuid")
	tasks.StopTask(tid)
	c.SetResult(nil,nil,204)
}

func (c *IaaSController) CreateTask(){
	defer c.ServeJSON()
	task:=TaskJSON{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &task); err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	oldTask:=orm.Task{}
	res,err:=orm.GetTask(task.ID,&oldTask)
	if err!=nil{
			c.SetResult(err, nil, 400)
			return
	}
	newTask :=orm.Task{
		ID:task.ID,
		ProjectID:task.ID,
		RepoID:task.RepoID,
		Name:"iaas",
		Group:task.PlaybookParse.Group,
		Vars:task.PlaybookParse.Vars,
		Status:"waiting",
	}
	if oldTask.ID!=""&&res{
		err:=orm.UpdateTask(&newTask)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
		tasks.AddTask(newTask.ID)
		c.SetResult(nil,newTask.ID,200,"task_id")
		return
	}

	err=orm.CreateTask(&newTask)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	tasks.AddTask(newTask.ID)
	c.SetResult(nil,newTask.ID,200,"task_id")
}

func (c *IaaSController)GetTask(){
	defer c.ServeJSON()
	tid :=c.GetString("uuid")
	task:=orm.Task{}
	res,err:=orm.GetTask(tid,&task)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	if !res{
		msg:=ErrorMsg{
			400,
			"项目还未编排",
		}
		c.SetResult(nil,msg,400)
		return
	}
	repo:=orm.Repository{}
	err=orm.GetRepoByID(task.RepoID,&repo)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	PlaybookParse:=PlaybookParse{
		Group:task.Group,
		Vars:task.Vars,
	}
	data:=TaskJSON{
		ID:task.ID,
		RepoName:repo.Name,
		RepoID:task.RepoID,
		PlaybookParse:PlaybookParse, 
	}
	c.SetResult(nil,data,200)
}
