package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/sockets"
	"github.com/gzsunrun/ansible-manager/core/function"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/core/tasks"
	_ "github.com/gzsunrun/ansible-manager/routers"
	"github.com/astaxie/beego"
)

const (
	SERVICE_NAME = "ansible-manager"
	SERVICE_DESC = "ansible-manager"
	LOG_PATH     = "/var/log/ansible-manager/log.log"
	CONFIG_PATH  = "/etc/ansible-manager/ansible-manager.conf"
)

func run() {
	os.MkdirAll(config.Cfg.Ansible.WorkPath, 0664)
	sockets.StartWS()
	function.NewS3Client()
	orm.NewDB()
	go tasks.RunTask()
	beego.BConfig.AppName = "sunruniaas-ansible"
	beego.BConfig.RunMode = beego.PROD
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.Log.FileLineNum = true
	beego.SetLogFuncCall(true)
	beego.SetStaticPath("/ui", "/root/go/src/github.com/gzsunrun/ansible-manager/public")
	beego.BConfig.Listen.HTTPPort = config.Cfg.Ansible.Port
	beego.Run()
}

func main() {
	err := logs.SetLogger(logs.AdapterMultiFile, `{"filename":"`+LOG_PATH+`","separate":["error"]}`)
	if err != nil {
		logs.Error("fail to config logrus")
		return
	}
	err = config.NewConfig(CONFIG_PATH)
	if err != nil {
		logs.Error(err)
		return
	}
	config.BackGroundService(SERVICE_NAME, SERVICE_DESC, nil, run)
}
