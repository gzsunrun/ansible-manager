package output

import (
	"github.com/gzsunrun/ansible-manager/core/config"
)

type LogOutput interface {
	Read() ([]string, error)
	Write(msg string) error
	Close() error
	Clean() error
}

func NewLogOutput(id string) (LogOutput, error) {
	logger, err := NewFileLog(config.Cfg.FileLog.Path + id + ".log")
	return logger, err
}
