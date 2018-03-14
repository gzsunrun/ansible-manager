package orm

import (
	"time"

	log "github.com/astaxie/beego/logs"
)
type Timer struct{
	ID				string		`xorm:"timer_id" json:"timer_id"`
	TaskID 			string		`xorm:"task_id" json:"task_id"`
	UserID			string		`xorm:"user_id" json:"user_id"`
	Name 			string		`xorm:"timer_name" json:"timer_name"`
	Start 			int			`xorm:"timer_start" json:"timer_start"`
	Interval 		int			`xorm:"timer_interval" json:"timer_interval"`
	Surplus 		int			`xorm:"-" json:"timer_surplus"`
	Repeat 			int			`xorm:"timer_repeat" json:"timer_repeat"`
	Status 			bool		`xorm:"timer_status" json:"timer_status"`
	Created 		time.Time	`xorm:"created" json:"created"`
}

func CreateTimer(t *Timer)error{
	_,err:=MysqlDB.Table("ansible_timer").Insert(t)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func FindTimers(uid string)(*[]Timer,error){
	var timers []Timer
	err:=MysqlDB.Table("ansible_timer").Where("user_id=?",uid).Find(&timers)
	if err!=nil{
		log.Error(err)
		return nil ,err
	}
	return &timers,err
}


func GetTimer(tid string)(bool,*Timer,error){
	var timer Timer
	res,err:=MysqlDB.Table("ansible_timer").Where("timer_id=?",tid).Get(&timer)
	if err!=nil{
		log.Error(err)
		return false,nil ,err
	}
	return res,&timer,err
}

func UpdateTimerStatus(t *Timer)error{
	_,err:=MysqlDB.Table("ansible_timer").Cols("timer_status").Where("timer_id=?",t.ID).Update(t)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func UpdateTimerRun(t *Timer)error{
	_,err:=MysqlDB.Table("ansible_timer").Cols("timer_repeat","timer_start").Where("timer_id=?",t.ID).Update(t)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func UpdateTimer(t *Timer)error{
	_,err:=MysqlDB.Table("ansible_timer").Where("timer_id=?",t.ID).Update(t)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func UpdateTimerStart(t *Timer)error{
	_,err:=MysqlDB.Table("ansible_timer").Cols("timer_status","timer_start").Where("timer_id=?",t.ID).Update(t)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func DelTimer(tid string)error{
	timer:=new(Timer)
	_,err:=MysqlDB.Table("ansible_timer").Where("timer_id=?",tid).Delete(timer)
	if err!=nil{
		log.Error(err)
		return err
	}
	return nil
}