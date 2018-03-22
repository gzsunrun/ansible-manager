package main

import (
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/role"
	"github.com/gzsunrun/ansible-manager/core/sockets"
	"github.com/gzsunrun/ansible-manager/core/storage"
	"github.com/gzsunrun/ansible-manager/core/tasks"
	_ "github.com/gzsunrun/ansible-manager/routers"
)

const (
	// SERVICENAME service name
	SERVICENAME = "ansible-manager"
	// SERVICEDESC service desc
	SERVICEDESC = "ansible-manager"
	// LOGPATH log file path
	LOGPATH     = "/var/log/ansible-manager/log.log"
	// CONFIGPATH config file path
	CONFIGPATH  = "/etc/ansible-manager/ansible-manager.conf"
)

// run start process
func run() {
	config.SetLog(LOGPATH)
	err := config.NewConfig(CONFIGPATH)
	if err != nil {
		logs.Error(err)
		return
	}
	os.MkdirAll(config.Cfg.Common.WorkPath, 0664)
	sockets.StartWS()
	err = storage.SetStorage()
	if err != nil {
		return
	}
	orm.NewDB()
	go tasks.RunTask()
	err = kv.SetKVClient()
	if err != nil {
		logs.Error(err)
		return
	}
	role.Run()
	beego.BConfig.AppName = "ansible-manager"
	beego.BConfig.RunMode = beego.PROD
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.Log.FileLineNum = true
	if config.Cfg.Common.UAPI {
		beego.SetStaticPath("/ui", "public/")
	}
	beego.BConfig.Listen.HTTPPort = config.Cfg.Common.Port
	beego.Run()
}

func main() {
	config.BackGroundService(SERVICENAME, SERVICEDESC, nil, run)
}
