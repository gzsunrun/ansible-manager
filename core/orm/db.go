package orm

import (
	"fmt"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/gzsunrun/ansible-manager/core/config"
)

// MysqlDB mysql engine
var MysqlDB *xorm.Engine

// NewDB new db
func NewDB() error {
	var err error
	if config.Cfg.DBDriver == "mysql" {
		dbURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
			config.Cfg.Mysql.MysqlUser,
			config.Cfg.Mysql.MysqlPassword,
			config.Cfg.Mysql.MysqlURL,
			config.Cfg.Mysql.MysqlName)
		MysqlDB, err = xorm.NewEngine("mysql", dbURL)
		return err
	}
	MysqlDB, err = xorm.NewEngine("sqlite3", config.Cfg.Sqlite3.Path)
	return err
}

func Import(path string) error {
	_, err := MysqlDB.ImportFile(path)
	return err
}
