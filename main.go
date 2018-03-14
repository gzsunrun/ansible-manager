package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/sockets"
	"github.com/gzsunrun/ansible-manager/core/storage"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/role"
	"github.com/gzsunrun/ansible-manager/core/tasks"
	_ "github.com/gzsunrun/ansible-manager/routers"
	"github.com/astaxie/beego"
)

const (
	SERVICE_NAME = "ansible-manager"
	SERVICE_DESC = "ansible-manager"
	LOG_PATH     = "/var/log/ansible-manager/log.log"
	CONFIG_PATH  = "/etc/ansible-manager/ansible-manager.conf"
	HtmlPath 	="/usr/local/html/ansible-manager/public/"
)

func run() {
	config.SetLog(LOG_PATH)
	err := config.NewConfig(CONFIG_PATH)
	if err != nil {
		logs.Error(err)
		return
	}
	os.MkdirAll(config.Cfg.Common.WorkPath, 0664)
	sockets.StartWS()
	err=storage.SetStorage()
	if err != nil {
		return
	}
	orm.NewDB()
	go tasks.RunTask()
	err=kv.SetKVClient()
	if err != nil {
		logs.Error(err)
		return
	}
	role.Run()
	beego.BConfig.AppName = "ansible-manager"
	beego.BConfig.RunMode = beego.PROD
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.Log.FileLineNum = true
	beego.SetStaticPath("/ui", HtmlPath)
	beego.BConfig.Listen.HTTPPort = config.Cfg.Common.Port
	beego.Run()
}

func main() {
	config.BackGroundService(SERVICE_NAME, SERVICE_DESC, nil, run)
}
