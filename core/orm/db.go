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
	dbURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		config.Cfg.Mysql.MysqlUser,
		config.Cfg.Mysql.MysqlPassword,
		config.Cfg.Mysql.MysqlURL,
		config.Cfg.Mysql.MysqlName)
	MysqlDB, err = xorm.NewEngine("mysql", dbURL)
	return err
}
