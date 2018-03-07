package tasks

import (
	"bufio"
	"encoding/json"
	"os/exec"

	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/sockets"
	"github.com/gzsunrun/ansible-manager/core/output"
	log "github.com/astaxie/beego/logs"
)

type Task struct {
	Desc   	orm.Task
	LO 		output.LogOutput
}

func (t *Task) log(msg string) {
	t.LO.Write(msg)
	b, err := json.Marshal(&map[string]interface{}{
		"type":    "log",
		"output":  msg,
		"task_id": t.Desc.ID,
	})

	if err != nil {
		log.Error(err)
		return
	}
	sockets.Message(t.Desc.ID, b)
}

func (t *Task) updateStatus() {
	b, err := json.Marshal(&map[string]interface{}{
		"type":   "update",
		"start":  t.Desc.Start,
		"end":    t.Desc.End,
		"status": t.Desc.Status,
	})

	if err != nil {
		panic(err)
	}
	orm.UpdateTask(&t.Desc)
	sockets.Message(t.Desc.ID, b)
}

func SendLog(taskID, msg string) {
	sockets.Message(taskID, []byte(msg))
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

