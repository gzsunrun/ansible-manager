package orm

import (
	"fmt"

	"github.com/gzsunrun/ansible-manager/core/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)



var MysqlDB *xorm.Engine

func NewDB() error {
	var err error
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		config.Cfg.Ansible.MysqlUser, 
		config.Cfg.Ansible.MysqlPassword,
		config.Cfg.Ansible.MysqlURL, 
		config.Cfg.Ansible.MysqlName)
	MysqlDB, err = xorm.NewEngine("mysql", dbUrl)
	return err
}