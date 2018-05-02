package tasks

import (
	"bufio"
	"encoding/json"
	"os/exec"

	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/output"
	"github.com/gzsunrun/ansible-manager/core/sockets"
)

// Task task struct
type Task struct {
	Desc orm.Task
	Path string
	LO   output.LogOutput
}

// log write log
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

// updateStatus update task status
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

// SendLog send log to ws
func SendLog(taskID, msg string) {
	sockets.Message(taskID, []byte(msg))
}

// logPipe log pipe
func (t *Task) logPipe(scanner *bufio.Scanner) {
	for scanner.Scan() {
		t.log(scanner.Text())
	}
}

// logCmd log cmd
func (t *Task) logCmd(cmd *exec.Cmd) {
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	go t.logPipe(bufio.NewScanner(stderr))
	go t.logPipe(bufio.NewScanner(stdout))
}
