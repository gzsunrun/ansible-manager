package config

import (
	"github.com/astaxie/beego/logs"
)

func SetLog(path string)error{
	err := logs.SetLogger(logs.AdapterMultiFile, `{"filename":"`+path+`","separate":["error"]}`)
	if err != nil {
		logs.Error("fail to config logrus")
	}
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger("console")
	return err
}