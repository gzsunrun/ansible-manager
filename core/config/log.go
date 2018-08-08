package config

import (
	"github.com/hashwing/log"
)

// SetLog config logger
func SetLog(path string) error {
	logger,err:=log.NewBeegoLog(path,7,true)
	if err!=nil{
		log.Error("new a logger error:",err)
		return err
	}
	log.SetHlogger(logger)
	return nil
}

// SetLogger set logger
func SetLogger(logger log.Hlog){
	log.SetHlogger(logger)
}
