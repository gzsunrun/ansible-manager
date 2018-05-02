package orm

import (
	"time"

	log "github.com/astaxie/beego/logs"
)

// Task task table
type Task struct {
	ID        string    `xorm:"task_id" json:"task_id"`
	ProjectID string    `xorm:"project_id" json:"project_id"`
	RepoID    string    `xorm:"repo_id" json:"repo_id"`
	Name      string    `xorm:"task_name" json:"task_name"`
	Group     []Group   `xorm:"task_group" json:"task_group"`
	Vars      []Vars    `xorm:"task_vars" json:"task_vars"`
	Status    string    `xorm:"task_status" json:"task_status"`
	Start     time.Time `xorm:"task_start" json:"task_start"`
	End       time.Time `xorm:"task_end"  json:"task_end"`
	Tag       string    `xorm:"task_tag"  json:"task_tag"`
	Created   time.Time `xorm:"created" json:"created"`
}

// TaskList task output list
type TaskList struct {
	ID        string    `xorm:"task_id" json:"task_id"`
	ProjectID string    `xorm:"project_id" json:"-"`
	RepoID    string    `xorm:"repo_id" json:"repo_id"`
	Name      string    `xorm:"task_name" json:"task_name"`
	Group     []Group   `xorm:"task_group" json:"-"`
	Vars      []Vars    `xorm:"task_vars" json:"-"`
	Status    string    `xorm:"task_status" json:"task_status"`
	Start     time.Time `xorm:"task_start" json:"task_start"`
	End       time.Time `xorm:"task_end" json:"task_end"`
	Tag       string    `xorm:"task_tag"  json:"task_tag"`
	Created   time.Time `xorm:"created" json:"created"`
}

// GroupAttr group attr
type GroupAttr struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// GroupHosts group hosts
type GroupHosts struct {
	HostName string      `json:"host_name"`
	HostUUID string      `json:"host_uuid"`
	IP       string      `json:"-"`
	Alias    string      `json:"host_alias"`
	Attr     []GroupAttr `json:"attr"`
}

// Group group inventory
type Group struct {
	Name  string       `json:"group_name"`
	Hosts []GroupHosts `json:"hosts"`
	Attr  []GroupAttr  `json:"attr"`
}

// TaskCounts task count
type TaskCounts struct {
	Err     int64 `json:"error_total"`
	Success int64 `json:"success_total"`
	Total   int64 `json:"all_total"`
	Run     int64 `json:"run_total"`
}

// CreateTask create task
func CreateTask(task *Task) error {
	_, err := MysqlDB.Table("ansible_task").Insert(task)
	if err != nil {
		log.Error(err)
	}
	return err
}

// UpdateTask update task
func UpdateTask(task *Task) error {
	_, err := MysqlDB.Table("ansible_task").Where("task_id=?", task.ID).Update(task)
	if err != nil {
		log.Error(err)
	}
	return err
}

// UpdateTaskByProject update task by project
func UpdateTaskByProject(task *Task) error {
	_, err := MysqlDB.Table("ansible_task").Where("project_id=?", task.ProjectID).Update(task)
	if err != nil {
		log.Error(err)
	}
	return err
}

// GetTask get task
func GetTask(tid string, task interface{}) (bool, error) {
	res, err := MysqlDB.Table("ansible_task").Where("task_id=?", tid).Get(task)
	if err != nil {
		log.Error(err)
	}
	return res, err
}

// FindTasks find tasks
func FindTasks(pid string, task interface{}) error {
	err := MysqlDB.Table("ansible_task").Where("project_id=?", pid).Find(task)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// DelTask delete task
func DelTask(tid string) error {
	task := new(Task)
	_, err := MysqlDB.Table("ansible_task").Where("task_id=?", tid).Delete(task)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// GetTaskCount get task count
func GetTaskCount() (*TaskCounts, error) {
	task := new(Task)
	total, err := MysqlDB.Table("ansible_task").Count(task)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	successTotal, err := MysqlDB.Table("ansible_task").Where("task_status=?", "finish").Count(task)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	runTotal, err := MysqlDB.Table("ansible_task").Where("task_status=? or task_status=?", "running", "waiting").Count(task)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	errTotal, err := MysqlDB.Table("ansible_task").Where("task_status=?", "error").Count(task)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &TaskCounts{
		Err:     errTotal,
		Success: successTotal,
		Total:   total,
		Run:     runTotal,
	}, nil

}
