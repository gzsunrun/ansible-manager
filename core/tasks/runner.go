package tasks

import (
	//"encoding/json"
	"fmt"
	//"html/template"
	//"io/ioutil"
	"os"
	"os/exec"
	//"strconv"
	"time"

	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/config"
	//"github.com/gzsunrun/ansible-manager/core/function"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/output"
	"github.com/gzsunrun/ansible-manager/core/sockets"
	"github.com/gzsunrun/ansible-manager/core/storage"
	"github.com/gzsunrun/ansible-manager/core/template"
)

var workPath string
var taskList = make(chan string)

// AddTask add task
func AddTask(taskID string) {
	StopTask(taskID)
	taskList <- taskID
}

// RunTask run task
func RunTask() {
	workPath = config.Cfg.Common.WorkPath
	taskList = make(chan string)
	taskChan := make(chan bool, 5)
	go cmdst.Run()
	for {
		select {
		case id := <-taskList:
			go func() {
				taskChan <- true
				defer func() {
					<-taskChan
					kv.DefaultClient.DeleteTask(id)
				}()
				newTask(id)
			}()
		}
	}
}

// StopTask stop task
func StopTask(id string) {
	stopFunc := func(taskID string) {
		SendLog(taskID, "\nChanged:User Stop\n")
	}
	for k:=range cmdst.cmdTasks{
		if k==id{
			cmdst.StopCmd(id, stopFunc)
			return
		}
	}
	
}


// newTask new a task
func newTask(taskID string) error {
	task := orm.Task{}
	_, err := orm.GetTask(taskID, &task)
	if err != nil {
		log.Error(err)
		return err
	}
	lo, err := output.NewLogOutput(taskID)
	if err != nil {
		log.Error(err)
	}
	newTask := new(Task)
	newTask.Desc = task
	newTask.LO = lo
	time.Sleep(time.Second * 2)
	defer func() {
		newTask.Desc.End = time.Now()
		newTask.updateStatus()
		time.Sleep(2 * time.Second)
		sockets.CloseConn(newTask.Desc.ID)
		newTask.clcRepo()
	}()
	if task.Status != "waiting" {
		return nil
	}
	newTask.Desc.Status = "running"
	newTask.Desc.Start = time.Now()
	newTask.updateStatus()
	newTask.log("[DOWNLOAD REPOSITORY]")
	err = newTask.downloadRepository()
	if err != nil {
		newTask.log("fatal: error \n ")
		newTask.Desc.Status = "error"
		return err
	}
	newTask.log("ok: finish \n ")
	newTask.log("[INSTALL HOSTS AND VARS]")
	//err = newTask.installVars()
	wp,err := template.InstallVars(&(newTask.Desc), workPath+"/repo_"+newTask.Desc.ID)
	if err != nil {
		log.Error(err)
		newTask.log("fatal: error \n ")
		newTask.Desc.Status = "error"
		return err
	}
	newTask.Path=wp
	newTask.log("ok: finish \n ")
	err = newTask.runPlaybook()
	if err != nil {
		newTask.Desc.Status = "error"
		if err.Error() == "signal: killed" {
			newTask.Desc.Status = "stop"
			log.Info("stop task:", newTask.Desc.ID)
			return nil
		}
		return err
	}
	newTask.Desc.Status = "finish"
	log.Info("finish task:", newTask.Desc.ID)
	return nil
}

// runPlaybook run playbook
func (t *Task) runPlaybook() error {
	dir :=  workPath + "/repo_" + t.Desc.ID+"/"+t.Path
	log.Info("task %s start use dir: %s",t.Desc.ID,dir)
	args := make([]string, 0)
	args = append(args, "-i")
	args = append(args, dir+"/hosts")
	args = append(args, dir+"/index.yml")
	if t.Desc.Tag != "" {
		args = append(args, "--tags")
		args = append(args, `"`+t.Desc.Tag+`"`)
	}
	cmd := exec.Command("ansible-playbook", args...)
	cmd.Dir = dir
	cmd.Env = envVars(workPath, cmd.Dir, nil)
	t.logCmd(cmd)
	return cmdst.StartCmd(cmd, t.Desc.ID)
}

// envVars ...
func envVars(home string, pwd string, gitSSHCommand *string) []string {
	env := os.Environ()
	env = append(env, fmt.Sprintf("HOME=%s", home))
	env = append(env, fmt.Sprintf("PWD=%s", pwd))
	env = append(env, fmt.Sprintln("PYTHONUNBUFFERED=1"))

	if gitSSHCommand != nil {
		env = append(env, fmt.Sprintf("GIT_SSH_COMMAND=%s", *gitSSHCommand))
	}

	return env
}

// downloadRepository download repo
func (t *Task) downloadRepository() error {
	repo := orm.Repository{}
	err := orm.GetRepoByID(t.Desc.RepoID, &repo)
	if err != nil {
		log.Error(err)
		return err
	}

	repoParse := storage.StorageParse{
		LocalPath:  workPath + "/repo-tar-" + t.Desc.ID,
		RemotePath: repo.Path,
	}
	err = storage.Storage.Get(&repoParse)
	if err != nil {
		log.Error(err)
		return err
	}
	err = os.MkdirAll(workPath+"/repo_"+t.Desc.ID, 0664)
	if err != nil {
		log.Error(err)
		return err
	}
	cmd := exec.Command("tar", "zxvf", workPath+"/repo-tar-"+t.Desc.ID, "-C", workPath+"/repo_"+t.Desc.ID)
	err = cmdst.StartCmd(cmd, t.Desc.ID)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// clcRepo clean local repo
func (t *Task) clcRepo() error {

	args := make([]string, 0)
	args = append(args, "-rf")
	args = append(args, workPath+"/repo-tar-"+t.Desc.ID)
	args = append(args, workPath+"/repo_"+t.Desc.ID)
	cmd := exec.Command("rm", args...)
	err := cmd.Run()
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// PlaybookParse playbook parse struct
type PlaybookParse struct {
	Hosts []orm.HostsList `json:"hosts"`
	Group []orm.Group     `json:"group"`
	Vars  []orm.Vars      `json:"vars"`
}

// installVars create group vars
// func (t *Task) installVars() error {
// 	var hosts []orm.HostsList
// 	err := orm.FindHostFromProject(t.Desc.ProjectID, &hosts)
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	for i, g := range t.Desc.Group {
// 		for j, h := range g.Hosts {
// 			for _, hh := range hosts {
// 				if hh.ID == h.HostUUID {
// 					t.Desc.Group[i].Hosts[j].IP = hh.IP
// 					t.Desc.Group[i].Hosts[j].HostName = hh.HostName
// 				}
// 			}
// 		}
// 	}
// 	for i, val := range hosts {
// 		if val.Key != "" && val.Password != "" {
// 			if function.AuthKey(val) != "success" {
// 				hosts[i].Key = ""
// 			}
// 		}
// 		if val.HostName == "" {
// 			hosts[i].HostName = "host" + strconv.Itoa(i)
// 		}
// 		if val.Key != "" {
// 			if err := ioutil.WriteFile(workPath+"/repo_"+t.Desc.ID+"/key-"+val.IP, []byte(val.Key), 0600); err != nil {
// 				log.Error(err)
// 				return err
// 			}
// 		}
// 	}
// 	playbookParse := PlaybookParse{
// 		Hosts: hosts,
// 		Group: t.Desc.Group,
// 		Vars:  t.Desc.Vars,
// 	}

// 	for _, val := range t.Desc.Vars {
// 		v, err := json.Marshal(val.Value.Vars)
// 		if err != nil {
// 			log.Error(err)
// 			return err
// 		}
// 		err = ioutil.WriteFile(workPath+"/repo_"+t.Desc.ID+"/"+val.Path, v, 0600)
// 		if err != nil {
// 			log.Error(err)
// 			return err
// 		}
// 	}

// 	tmpl, err := template.ParseFiles(workPath + "/repo_" + t.Desc.ID + "/hosts")
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	fd, err := os.OpenFile(workPath+"/repo_"+t.Desc.ID+"/hosts", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	err = tmpl.Execute(fd, playbookParse)
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	return nil
// }

// installVars create group vars
