package db

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type User struct {
	ID       int    `xorm:"user_id" json:"user_id"`
	Account  string `xorm:"user_account" json:"user_account"`
	Password string `xorm:"user_password" json:"-"`
}

type Project struct {
	ID      int       `xorm:"project_id" json:"project_id"`
	UID     int       `xorm:"user_id" json:"user_id"`
	Name    string    `xorm:"project_name" json:"project_name"`
	Created time.Time `xorm:"project_created" json:"project_created"`
}

type Repository struct {
	ID   int    `xorm:"repo_id" json:"repo_id"`
	Name string `xorm:"repo_name" json:"repo_name"`
	Path string `xorm:"repo_path" json:"repo_path"`
	Desc string `xorm:"repo_desc" json:"repo_desc"`
}

type Vars struct {
	ID     int    `xorm:"vars_id" json:"vars_id"`
	RepoID int    `xorm:"repo_id" json:"repo_id"`
	Type   string `xorm:"vars_type" json:"vars_type"`
	Name   string `xorm:"vars_name" json:"vars_name"`
	Path   string `xorm:"vars_path" json:"vars_path"`
	Parse  string `xorm:"vars_value" json:"vars_value"`
}

type Host struct {
	ID        int    `xorm:"host_id" json:"host_id"`
	ProjectID int    `xorm:"project_id" json:"project_id"`
	Alias     string `xorm:"host_alias" json:"host_alias"`
	Name      string `xorm:"host_name" json:"host_name"`
	IP        string `xorm:"host_ip" json:"host_ip"`
	User      string `xorm:"host_user" json:"host_user"`
	Password  string `xorm:"host_password" json:"-"`
	Key       string `xorm:"host_key" json:"-"`
}

type Task struct {
	ID      int       `xorm:"task_id" json:"task_id"`
	UID     int       `xorm:"user_id" json:"user_id"`
	TplID   int       `xorm:"tpl_id" json:"tpl_id"`
	Tag     string    `xorm:"playbook_tag" json:"playbook_tag"`
	Name    string    `xorm:"task_name" json:"task_name"`
	Status  string    `xorm:"task_status" json:"task_status"`
	Created time.Time `xorm:"task_created" json:"task_created"`
	Start   time.Time `xorm:"task_start" json:"task_start"`
	End     time.Time `xorm:"task_end" json:"task_end"`
}

type Template struct {
	ID            int    `xorm:"tpl_id" json:"tpl_id"`
	ProjectID     int    `xorm:"project_id" json:"project_id"`
	RepoID        int    `xorm:"repo_id" json:"repo_id"`
	Name          string `xorm:"tpl_name" json:"tpl_name"`
	Playbook      string `xorm:"playbook" json:"playbook"`
	PlaybookParse string `xorm:"playbook_parse" json:"playbook_parse"`
}

type Output struct {
	TaskID int       `xorm:"task_id" json:"task_id"`
	Output string    `xorm:"output" json:"output"`
	Time   time.Time `xorm:"time" json:"time"`
}

var MysqlDB *xorm.Engine

func NewDB(url string) error {
	var err error
	MysqlDB, err = xorm.NewEngine("mysql", url)
	return err
}
