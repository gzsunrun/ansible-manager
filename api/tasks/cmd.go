package tasks

import (
	"os/exec"
	"syscall"

	"github.com/astaxie/beego/logs"
)

type cmdStruct struct {
	startChan  chan *cmdTask
	stopChan   chan int
	finishChan chan int
	cmdTasks   map[int]int
}

type cmdTask struct {
	pid    int
	taskID int
}

var cmdst = cmdStruct{
	startChan:  make(chan *cmdTask),
	stopChan:   make(chan int),
	finishChan: make(chan int),
	cmdTasks:   make(map[int]int),
}

func (c *cmdStruct) run() {
	for {
		select {
		case ct := <-c.startChan:
			c.cmdTasks[ct.taskID] = ct.pid
		case taskID := <-c.stopChan:
			if c.cmdTasks[taskID] > 0 {
				err := killByPID(c.cmdTasks[taskID])
				if err == nil {
					logByID(taskID, "\033[33m changed: user stop \033[0m\n")
				}
			}
		case taskID := <-c.finishChan:
			delete(c.cmdTasks, taskID)
		}
	}
}

func killByPID(pid int) error {
	err := syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		logs.Error(err)
	}
	return err
}

func (cs *cmdStruct) startCmd(c *exec.Cmd, taskID int) error {
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
		logs.Error(err)
		return err
	}
	return nil
}

func (c *cmdStruct) stopCmd(taskID int) {
	c.stopChan <- taskID
}
