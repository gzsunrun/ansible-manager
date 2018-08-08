package tasks

import (
	"os/exec"
	"syscall"

	"github.com/hashwing/log"
)

type cmdStruct struct {
	startChan  chan *cmdTask
	stopChan   chan string
	finishChan chan string
	cmdTasks   map[string]int
	stopFunc   func(taskID string)
}

type cmdTask struct {
	pid    int
	taskID string
}

var cmdst = cmdStruct{
	startChan:  make(chan *cmdTask),
	stopChan:   make(chan string),
	finishChan: make(chan string),
	cmdTasks:   make(map[string]int),
}

func (c *cmdStruct) Run() {
	for {
		select {
		case ct := <-c.startChan:
			c.cmdTasks[ct.taskID] = ct.pid
		case taskID := <-c.stopChan:
			if c.cmdTasks[taskID] > 0 {
				err := killByPID(c.cmdTasks[taskID])
				if err == nil {
					c.stopFunc(taskID)
				}
			}
		case taskID := <-c.finishChan:
			delete(c.cmdTasks, taskID)
		}
	}
}

func killByPID(pid int) error {
	err := syscall.Kill(-pid, syscall.SIGKILL)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (cs *cmdStruct) StartCmd(c *exec.Cmd, taskID string) error {
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := c.Start()
	if err != nil {
		return err
	}
	ct := &cmdTask{
		pid:    c.Process.Pid,
		taskID: taskID,
	}
	cs.startChan <- ct
	err = c.Wait()
	cs.finishChan <- ct.taskID
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *cmdStruct) StopCmd(taskID string, f ...func(taskID string)) {
	c.stopFunc = f[0]
	c.stopChan <- taskID
}
