package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/mux"
	"github.com/gzsunrun/ansible-manager/api/db"
	"github.com/gzsunrun/ansible-manager/api/router"
	"github.com/gzsunrun/ansible-manager/api/s3"
	"github.com/gzsunrun/ansible-manager/api/sockets"
	"github.com/gzsunrun/ansible-manager/api/tasks"
	"github.com/gzsunrun/ansible-manager/config"
)

const (
	SERVICE_NAME = "ansible-manager"
	SERVICE_DESC = "ansible-manager"
	LOG_PATH     = "/var/log/ansible-manager/log.log"
	CONFIG_PATH  = "/etc/ansible-manager/ansible-manager.conf"
)

func run() {
	os.MkdirAll(config.Cfg.AnsibleManager.WorkPath+"/repo", 0664)
	s3.NewClient()
	sockets.StartWS()
	db.NewDB(config.Cfg.AnsibleManager.MysqlURL)
	go tasks.RunTask()
	root := mux.NewRouter()
	router.NewRouter(root)
	err := http.ListenAndServe(":"+strconv.Itoa(config.Cfg.AnsibleManager.Port), root)
	if err != nil {
		logs.Error(err)
	}
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
