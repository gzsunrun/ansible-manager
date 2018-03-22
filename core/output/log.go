package output

import (
	"github.com/gzsunrun/ansible-manager/core/config"
)

type LogOutput interface {
	Read() ([]byte, error)
	Write(msg string) error
}

func NewLogOutput(id string) (LogOutput, error) {
	return NewFileLog(config.Cfg.FileLog.Path + id + ".log")
}
