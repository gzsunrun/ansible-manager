package log

// Hlog log libaray
type Hlog interface {
	Debug(f interface{},v ...interface{})
	Info(f interface{},v ...interface{})
	Warn(f interface{},v ...interface{})
	Error(f interface{},v ...interface{})
}

var globalLog Hlog

func init(){
	globalLog,_=NewBeegoLog("",7,true)
}

// SetHlogger reg a logger in global
func SetHlogger(logger Hlog){
	globalLog=logger
}


// Debug log debug
func Debug(f interface{},v ...interface{}){
	globalLog.Debug(f,v...)
}

// Info log info
func Info(f interface{},v ...interface{}){
	globalLog.Info(f,v...)
}

// Warn log warn
func Warn(f interface{},v ...interface{}){
	globalLog.Debug(f,v...)
}

// Error log error
func Error(f interface{},v ...interface{}){
	globalLog.Error(f,v...)
}

