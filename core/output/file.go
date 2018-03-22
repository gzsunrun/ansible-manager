package output

import (
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileLog struct {
	log      *logs.BeeLogger
	FilePath string
}

func NewFileLog(path string) (*FileLog, error) {
	dir, _ := filepath.Split(path)
	err := os.MkdirAll(dir, 0664)
	if err != nil {
		return nil, err
	}
	os.Remove(path)
	fg := new(FileLog)
	fg.log = logs.NewLogger()
	fg.FilePath = path
	return fg, fg.log.SetLogger("file", `{"filename":"`+path+`"}`)
}

func (fg *FileLog) Write(msg string) error {
	fg.log.Info(msg)
	return nil
}

func (fg *FileLog) Read() ([]byte, error) {
	return ioutil.ReadFile(fg.FilePath)
}
