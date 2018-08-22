## log

根据个人要求二次封装的log库


### 使用

```go

package main

import (
    "github.com/hashwing/log"
)

func main(){
    // 实例化logger,使用beego log,
    // 参数分别是：日志文件路径（空为不输出）、是否开启Debug、是否输出控制台、是否把Error级日志输出单独文件
    logger,err:=log.NewBeegoLog("/var/log/test/access.log",false,true,true)
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

```