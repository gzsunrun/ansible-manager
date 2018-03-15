package controllers

import (
	"encoding/json"
	"html/template"

	log "github.com/astaxie/beego/logs"
	"github.com/satori/go.uuid"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/kv"
)

type TaskController struct{
	BaseController
}

func (c *TaskController)List(){
	defer c.ServeJSON()
	pid:=c.GetString("project_id")
	var tasks []orm.TaskList
	err:=orm.FindTasks(pid,&tasks)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	c.SetResult(nil,tasks,200)
}

func (c *TaskController)Create(){
	defer c.ServeJSON()
	task:=orm.Task{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &task); err != nil {
		log.Error(err)
		c.SetResult(err, nil, 400)
		return
	}
	if task.ID!=""{
		err:=orm.UpdateTask(&task)
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
		c.SetResult(nil,task.ID,200,"task_id")
		return
	}
	task.ID=uuid.Must(uuid.NewV4()).String()
	task.Status="created"
	err:=orm.CreateTask(&task)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil,task.ID,200,"task_id")
}

func (c *TaskController)Start(){
	defer c.ServeJSON()
	tid:=c.GetString("task_id")
	tag:=c.GetString("task_tag")
	task:=orm.Task{
		ID:tid,
		Status:"waiting",
		Tag:tag,
	}
	err:=orm.UpdateTask(&task)
	if err != nil {
		c.SetResult(err, nil, 400)
		return
	}
	master:=false
	worker:=false
	for _,node:=range kv.DefaultClient.GetStorage().Nodes{
		if node.Master{
			master=true
		}
		if node.Worker{
			worker=true
		}
	}
	if master&&worker{
		err=kv.DefaultClient.AddScheduler(kv.Task{
			ID:task.ID,
			Timer:false,
		})
		if err != nil {
			c.SetResult(err, nil, 400)
			return
		}
		c.SetResult(nil,tid,200,"task_id")
	}else{
		c.SetResult(err, nil, 400)
	}
	
}

func (c *TaskController)Stop(){
	defer c.ServeJSON()
	tid:=c.GetString("task_id")
	if _,ok:=kv.DefaultClient.GetStorage().Tasks[tid];ok{
		task:=new(orm.Task)
		task.ID=tid
		task.Status="stop"
		err :=orm.UpdateTask(task)
		if err!=nil{
			c.SetResult(err, nil, 400)
			return
		}
		c.SetResult(nil,nil,204)
	}
	err:=kv.DefaultClient.DeleteTask(tid)
	if err!=nil{
		c.SetResult(err, nil, 400)
		return
	}
	c.SetResult(nil,nil,204)
}

func (c *TaskController)Get(){
	defer c.ServeJSON()
	tid:=c.GetString("task_id")
	var task orm.Task
	_,err:=orm.GetTask(tid,&task)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	c.SetResult(nil,task,200)
}

func (c *TaskController)Del(){
	defer c.ServeJSON()
	tid:=c.GetString("task_id")
	err:=kv.DefaultClient.DeleteTask(tid)
	if err!=nil{
		c.SetResult(err, nil, 400)
		return
	}
	err=orm.DelTask(tid)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	c.SetResult(nil,nil,204)
}


type W struct{
	Str string
}
func(w *W)Write(p []byte)(int,error){
	w.Str+=string(p)
   	return 0,nil
}

type NewPlaybookParse struct {
	Hosts []orm.HostsList  						`json:"hosts"`
	Group []orm.Group 						`json:"group"`
	Vars  map[string]map[string]interface{} `json:"vars"`
}

func (c *TaskController)GetNotes(){
	defer c.ServeJSON()
	tid:=c.GetString("uuid")
	var task orm.Task
	res,err:=orm.GetTask(tid,&task)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	if !res{
		msg:=ErrorMsg{
			400,
			"应用未配置或配置不正确",
		}
		c.SetResult(nil,msg,400)
		return
	}
	var repo orm.Repository
	err=orm.GetRepoByID(task.RepoID,&repo)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	var hosts []orm.HostsList
	err = orm.FindHostFromProject(task.ProjectID,&hosts)
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	tmpl, err := template.New("tpl").Parse(repo.Note)
	if err!=nil{
		log.Error(err)
		msg:=ErrorMsg{
			400,
			"应用未配置或配置不正确",
		}
		c.SetResult(nil,msg,400)
		return
	}
	w:=&W{}
	np :=NewPlaybookParse{}
	np.Hosts=hosts
	for i,g:=range task.Group{
		for j,h:=range g.Hosts{
			for _,hh:=range hosts{
				if hh.ID==h.HostUUID{
					log.Info(hh.IP)
					task.Group[i].Hosts[j].IP=hh.IP
					task.Group[i].Hosts[j].HostName=hh.HostName
				}
			}
		}
	}
	np.Group=task.Group
	np.Vars=make(map[string]map[string]interface{})
	for _,v:=range task.Vars{
			np.Vars[v.Name]=v.Value.Vars	
	}
	err = tmpl.Execute(w, np)
	if err!=nil{
		log.Error(err)
		msg:=ErrorMsg{
			400,
			"应用未配置或配置不正确",
		}
		c.SetResult(nil,msg,400)
		return
	}
	c.SetResult(nil, w.Str, 200,"notes")
}


func (c *TaskController)GetTaskCount(){
	defer c.ServeJSON()
	counts,err:=orm.GetTaskCount()
	if err!=nil{
		c.SetResult(err,nil,400)
		return
	}
	c.SetResult(nil,counts,200)
}

func (c *TaskController)GetNodes(){
	defer c.ServeJSON()
	nodes :=make([]kv.Node,0)
	for _,node:=range kv.DefaultClient.GetStorage().Nodes{
		nodes=append(nodes,node)
	}
	c.SetResult(nil,nodes,200)
}