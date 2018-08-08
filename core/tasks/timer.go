package tasks

import (
	"time"

	"github.com/hashwing/log"
	"github.com/gzsunrun/ansible-manager/core/kv"
	"github.com/gzsunrun/ansible-manager/core/orm"
)

// TChan timer chan
var TChan = make(map[string]chan bool)

// SetTimer set timer
func SetTimer(timerID string) {
	res, timer, err := orm.GetTimer(timerID)
	if err != nil {
		log.Error(err)
	}
	if !res {
		return
	}
	if timer.Repeat == 0 {
		return
	}
	timer.Status = true
	var interval int64
	if timer.Start == 0 {
		timer.Start = int(time.Now().Unix())
		interval = int64(timer.Interval)
	} else {
		interval = (time.Now().Unix() - int64(timer.Start)) % int64(timer.Interval)
		timer.Start = int(time.Now().Unix() - interval)
	}
	err = orm.UpdateTimerStart(timer)
	if err != nil {
		log.Error(err)
	}

	t := time.NewTicker(time.Duration(interval) * time.Second)
	TChan[timerID] = make(chan bool, 1)
	done := make(chan bool, 1)
	log.Info("run timer:", timerID)
	go func(timer orm.Timer) {
		for {
			select {
			case <-t.C:
				if timer.Repeat == 0 {
					TChan[timer.ID] <- true
				} else {
					log.Info("timer run task:", timer.ID)
					kvTask := kv.Task{
						ID:    timer.TaskID,
						Timer: false,
					}
					err := kv.DefaultClient.AddScheduler(kvTask)
					if err != nil {
						log.Error(err)
					}
					timer.Start = int(time.Now().Unix())
					if timer.Repeat > 0 {
						timer.Repeat--
					}
					err = orm.UpdateTimerRun(&timer)
					if err != nil {
						log.Error(err)
					}
					if timer.Repeat == 0 {
						TChan[timer.ID] <- true
					} else {
						t = time.NewTicker(time.Duration(timer.Interval) * time.Second)
					}
				}

			case <-TChan[timerID]:
				close(done)
				log.Info("stop timer:", timer.ID)
				timer.Status = false
				err := orm.UpdateTimerStatus(&timer)
				if err != nil {
					log.Error(err)
				}
				err = kv.DefaultClient.DeleteTask(timer.ID)
				if err != nil {
					log.Error(err)
				}
				return
			}
		}
	}(*timer)
	<-done
	return
}

// StopTimer stop timer
func StopTimer(tid string) {
	TChan[tid] <- true
}
