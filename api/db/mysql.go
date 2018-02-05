package db

import (
	"time"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gzsunrun/ansible-manager/config"
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
	HostName  string `xorm:"host_name" json:"host_name"`
	IP        string `xorm:"host_ip" json:"host_ip"`
	User      string `xorm:"host_user" json:"host_user"`
	Password  string `xorm:"host_password" json:"host_password"`
	Key       string `xorm:"host_key" json:"host_key"`
}

type GroupAttr struct {
	Key   string `json:"key"`
	Type string  `json:"type"`
	Value string `json:"value"`
}

type GroupHosts struct {
	//Name string          `json:"host_name"`
	HostUUID string      `json:"host_uuid"`
	Attr     []GroupAttr `json:"attr"`
	HostName string 	 `json:"host_name"`
	IP       string 	 `json:"-"`
}

type GroupRole struct {
	Name  string       `json:"group_name"`
	Hosts []GroupHosts `json:"hosts"`
	Attr  []GroupAttr  `json:"attr"`
}

type RoleVars struct {
	Path  string 				 `json:"vars_path"`
	Value map[string]interface{} `json:"vars_value"`
	Name  string 				 `json:"vars_name"`
}

type Parse struct{
	Hosts []Host 		`json:"hosts"`
	Group []GroupRole 	`json:"group"`
	Vars  []RoleVars  	`json:"vars"`
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

func NewDB() error {
	var err error
	dbURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		config.Cfg.AnsibleManager.MysqlUser,
		config.Cfg.AnsibleManager.MysqlPassword,
		config.Cfg.AnsibleManager.MysqlURL,
		config.Cfg.AnsibleManager.MysqlName)

	MysqlDB, err = xorm.NewEngine("mysql", dbURL)
	return err
}
