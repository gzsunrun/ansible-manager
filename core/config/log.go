package config

import (
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

// SetLog config log path
func SetLog(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0664)
	if err != nil {
		logs.Error("fail to create log dir")
		return err
	}
	err = logs.SetLogger(logs.AdapterMultiFile, `{"filename":"`+path+`","separate":["error"]}`)
	if err != nil {
		logs.Error("fail to config logrus")
	}
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger("console")
	return err
}
