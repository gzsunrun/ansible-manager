package main

import (
	"os"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gzsunrun/ansible-manager/asset"
	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/helm"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/role"
	"github.com/gzsunrun/ansible-manager/core/sockets"
	"github.com/gzsunrun/ansible-manager/core/storage"
	"github.com/gzsunrun/ansible-manager/core/tasks"
	_ "github.com/gzsunrun/ansible-manager/routers"
	"github.com/hashwing/log"
)

const (
	// SERVICENAME service name
	SERVICENAME = "ansible-manager"
	// SERVICEDESC service desc
	SERVICEDESC = "ansible-manager"
	// LOGPATH log file path
	LOGPATH = "/var/log/ansible-manager/log.log"
	// CONFIGPATH config file path
	CONFIGPATH = "/etc/ansible-manager/ansible-manager.conf"
)

// run start process
func run() {
	config.SetLog(LOGPATH)
	err := config.NewConfig(CONFIGPATH)
	if err != nil {
		log.Error(err)
		return
	}
	helm.InitHarbor(config.Cfg.Harbor.URL, config.Cfg.Harbor.Repo)
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
		log.Error(err)
		return
	}
	role.Run()
	beego.BConfig.AppName = "ansible-manager"
	beego.BConfig.RunMode = beego.PROD
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.Log.FileLineNum = true
	beego.BConfig.WebConfig.DirectoryIndex = true
	if config.Cfg.Common.UAPI {
		err := asset.RestoreAssets("/var/lib/amgr/", "public")
		if err != nil {
			log.Error(err)
			return
		}
		beego.SetStaticPath("/ui", "/var/lib/amgr/public")
		//beego.SetStaticPath("/ui", "./public/")
	}
	beego.BConfig.Listen.HTTPPort = config.Cfg.Common.Port
	beego.Run()
}

func main() {
	config.BackGroundService(SERVICENAME, SERVICEDESC, nil, run)
}
