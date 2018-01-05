package tasks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/api/db"
	"github.com/gzsunrun/ansible-manager/api/sockets"
)

func (t *Task) log(msg string) {
	now := time.Now()
	if strings.HasPrefix(msg, "ok") {
		msg = "\033[32m " + msg + "\033[0m"
	}
	if strings.HasPrefix(msg, "fatal") {
		msg = "\033[31m " + msg + "\033[0m"
	}
	if strings.HasPrefix(msg, "changed") {
		msg = "\033[33m " + msg + "\033[0m"
	}
	b, err := json.Marshal(&map[string]interface{}{
		"type":    "log",
		"output":  msg,
		"time":    now,
		"task_id": t.TaskContent.ID,
	})

	if err != nil {
		logs.Error(err)
		return
	}

	sockets.Message(t.TaskContent.ID, b)
	t.LogMsg += msg + "\n"
}

func logByID(taskID int, msg string) {
	now := time.Now()
	b, err := json.Marshal(&map[string]interface{}{
		"type":    "log",
		"output":  msg,
		"time":    now,
		"task_id": taskID,
	})
	if err != nil {
		logs.Error(err)
		return
	}
	sockets.Message(taskID, b)
}

func (t *Task) saveLog() {
	now := time.Now()
	_, err := db.MysqlDB.Exec("insert into ansible_task_output (task_id,time,output) values (?,?,?)  ", t.TaskContent.ID, now, t.LogMsg)
	if err != nil {
		logs.Error(err)
	}
}

func (t *Task) updateStatus() {
	status := "\033[32m " + t.TaskContent.Status + "\033[0m"
	if t.TaskContent.Status == "error" {
		status = "\033[31m " + t.TaskContent.Status + "\033[0m"
	}
	b, err := json.Marshal(&map[string]interface{}{
		"type":   "update",
		"start":  t.TaskContent.Start,
		"end":    t.TaskContent.End,
		"status": status,
	})

	if err != nil {
		panic(err)
	}
	t.LogMsg += "\nStatus:" + status + "\n"
	sockets.Message(t.TaskContent.ID, b)
	if _, err := db.MysqlDB.Exec("update ansible_task set task_status=?, task_start=?, task_end=? where task_id=?", t.TaskContent.Status, t.TaskContent.Start, t.TaskContent.End, t.TaskContent.ID); err != nil {
		fmt.Printf("Failed to update task status: %s\n", err.Error())
	}
}

func (t *Task) logPipe(scanner *bufio.Scanner) {
	for scanner.Scan() {
		t.log(scanner.Text())
	}
}

func (t *Task) logCmd(cmd *exec.Cmd) {
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	go t.logPipe(bufio.NewScanner(stderr))
	go t.logPipe(bufio.NewScanner(stdout))
}
