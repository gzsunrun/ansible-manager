package tasks

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gzsunrun/ansible-manager/api/db"
	"github.com/gzsunrun/ansible-manager/api/s3"
	"github.com/gzsunrun/ansible-manager/config"
	"github.com/astaxie/beego/logs"
)

var workPath string
var taskList chan int

type Task struct {
	TaskContent *db.Task
	TplContent  *db.Template
	LogMsg      string
}

func AddTask(taskID int) {
	taskList <- taskID
}

func RunTask() {
	workPath = config.Cfg.AnsibleManager.WorkPath
	taskList = make(chan int)
	taskChan := make(chan bool, 5)
	go cmdst.run()
	for {
		select {
		case id := <-taskList:
			go func() {
				taskChan <- true
				defer func() {
					<-taskChan
				}()
				newTask(id)
			}()
		}
	}
}

func StopTask(id int) {
	cmdst.stopCmd(id)
}

func newTask(taskID int) error {
	var task db.Task
	_, err := db.MysqlDB.Table("ansible_task").Where("task_id=?", taskID).Get(&task)
	if err != nil {
		logs.Error(err)
		return err
	}
	var tpl db.Template
	_, err = db.MysqlDB.Table("ansible_template").Where("tpl_id=?", task.TplID).Get(&tpl)
	if err != nil {
		logs.Error(err)
		return err
	}
	newTask := &Task{
		TaskContent: &task,
		TplContent:  &tpl,
	}
	time.Sleep(time.Second * 2)
	newTask.TaskContent.Status = "running"
	newTask.TaskContent.Start = time.Now()
	newTask.updateStatus()
	defer func() {
		newTask.TaskContent.End = time.Now()
		newTask.updateStatus()
		newTask.saveLog()
	}()
	newTask.log("[DOWNLOAD REPOSITORY]")
	err = newTask.downloadRepository()
	if err != nil {
		newTask.log("fatal: error \n ")
		newTask.TaskContent.Status = "error"
		return err
	}
	newTask.log("ok: finish \n ")
	newTask.log("[INSTALL VARS]")
	err = newTask.installVars()
	if err != nil {
		newTask.log("fatal: error \n ")
		newTask.TaskContent.Status = "error"
		return err
	}
	newTask.log("ok: finish \n ")
	err = newTask.runPlaybook()
	if err != nil {
		newTask.TaskContent.Status = "error"
		return err
	}
	newTask.clcRepo()
	newTask.TaskContent.Status = "finish"
	return nil
}

func (t *Task) runPlaybook() error {
	args := make([]string, 0)
	args = append(args, "-i")
	args = append(args, workPath+"/repo_"+strconv.Itoa(t.TaskContent.ID)+"/hosts")
	args = append(args, workPath+"/repo_"+strconv.Itoa(t.TaskContent.ID)+"/"+t.TplContent.Playbook)
	if t.TaskContent.Tag != "" {
		args = append(args, "--tags")
		args = append(args, `"`+t.TaskContent.Tag+`"`)
	}
	cmd := exec.Command("ansible-playbook", args...)
	cmd.Dir = workPath + "/repo_" + strconv.Itoa(t.TaskContent.ID)
	cmd.Env = envVars(workPath, cmd.Dir, nil)
	t.logCmd(cmd)
	return cmdst.startCmd(cmd, t.TaskContent.ID)
}

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

func (t *Task) downloadRepository() error {
	var value string
	_, err := db.MysqlDB.Table("ansible_repository").Where("repo_id=?", t.TplContent.RepoID).Cols("repo_path").Get(&value)
	if err != nil {
		logs.Error(err)
		return err
	}
	err = s3.S3Get(value, workPath+"/repo-tar-"+strconv.Itoa(t.TaskContent.ID))
	if err != nil {
		logs.Error(err)
		return err
	}
	err = os.MkdirAll(workPath+"/repo_"+strconv.Itoa(t.TaskContent.ID), 0664)
	if err != nil {
		logs.Error(err)
		return err
	}
	cmd := exec.Command("tar", "zxvf", workPath+"/repo-tar-"+strconv.Itoa(t.TaskContent.ID), "-C", workPath+"/repo_"+strconv.Itoa(t.TaskContent.ID))
	err = cmdst.startCmd(cmd, t.TaskContent.ID)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (t *Task) clcRepo() error {

	args := make([]string, 0)
	args = append(args, "-rf")
	args = append(args, workPath+"/repo-tar-"+strconv.Itoa(t.TaskContent.ID))
	args = append(args, workPath+"/repo_"+strconv.Itoa(t.TaskContent.ID))
	cmd := exec.Command("rm", args...)
	err := cmd.Run()
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (t *Task) installVars() error {
	if t.TplContent.PlaybookParse == "" {
		return nil
	}
	var data map[string][]map[string]interface{}
	err := json.Unmarshal([]byte(t.TplContent.PlaybookParse), &data)
	if err != nil {
		logs.Error(err)
		return err
	}
	for _, val := range data["hosts"] {
		hostKey, ok := val["host_key"].(string)
		hostIP, ok := val["host_ip"].(string)
		if ok && hostKey != "" {
			if err := ioutil.WriteFile(workPath+"/repo_"+strconv.Itoa(t.TaskContent.ID)+"/key-"+hostIP, []byte(hostKey), 0600); err != nil {
				logs.Error(err)
				return err
			}
		}
	}

	for _, val := range data["vars"] {
		var parseStr string
		varsValue, ok := val["vars_value"].(string)
		if ok {
			parseStr = strings.Replace(varsValue, "\\n", "\n", -1)
		}
		varsPath, ok := val["vars_path"].(string)
		if ok {
			err := ioutil.WriteFile(workPath+"/repo_"+strconv.Itoa(t.TaskContent.ID)+"/"+varsPath, []byte(parseStr), 0600)
			if err != nil {
				logs.Error(err)
				return err
			}
		}
	}
	tmpl, err := template.ParseFiles(workPath + "/repo_" + strconv.Itoa(t.TaskContent.ID) + "/hosts")
	if err != nil {
		logs.Error(err)
		return err
	}
	fd, err := os.OpenFile(workPath+"/repo_"+strconv.Itoa(t.TaskContent.ID)+"/hosts", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		logs.Error(err)
		return err
	}
	err = tmpl.Execute(fd, data)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil

}
