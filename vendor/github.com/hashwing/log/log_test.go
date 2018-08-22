package log_test

import (
	"testing"
	"github.com/hashwing/log"
)

func Test_Log(t *testing.T){
	// 实例化logger,使用beego log,参数分别是：日志文件路径（空为不输出）、是否开启Debug、是否输出控制台、是否把Error级日志输出单独文件
    logger,err:=log.NewBeegoLog("./test.log",7,false)
    if err!=nil{
        log.Warn(err)
    }
    // 将实例化的 logger 作为全局的logger
    log.SetHlogger(logger)

    // Debug
    log.Debug("This is a debug message")

    // Info
    log.Info("This","is","a","info","message")

    // Warn
    log.Warn("This is a %s message","warn")

    // Error
	log.Error("This is a error:",err)
}