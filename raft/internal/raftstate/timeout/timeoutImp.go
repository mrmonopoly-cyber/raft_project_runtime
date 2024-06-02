package timeout

import (
	"errors"
	"log"
	"sync"
	"time"
)

type timeout struct {
	timer             *time.Ticker
	duration          time.Duration
	timerNotification chan time.Time
}

type timeoutPool struct {
	timerMap sync.Map
}


// AddTimeout implements TimeoutPool.
func (t *timeoutPool) AddTimeout(name string, duration time.Duration) {
	var newTimer = timeout{
		timer:             time.NewTicker(duration),
		duration:          duration,
		timerNotification: make(chan time.Time),
	}

    newTimer.timer.Stop()
    go newTimer.notifyTimers()
	t.timerMap.Store(name, newTimer)
}

// GetTimeoutNotifycationChan implements TimeoutPool.
func (t *timeoutPool) GetTimeoutNotifycationChan(name string) (chan time.Time, error) {
    var timerInstance, err = t.findTimer(name)
    if err != nil{
        return nil,err
    }
    return timerInstance.timerNotification,nil
}

// RestartTimeout implements TimeoutPool.
func (t *timeoutPool) RestartTimeout(name string) error {
    var timerInstance, err = t.findTimer(name)
    if err != nil{
        return err
    }

	timerInstance.timer.Reset(timerInstance.duration)

	return nil
}

// StopTimeout implements TimeoutPool.
func (t *timeoutPool) StopTimeout(name string) error {
    var timerInstance, err = t.findTimer(name)
    if err != nil{
        return err
    }
    
    timerInstance.timer.Stop()
    return nil
}

//utility
func (t* timeoutPool) findTimer(name string) (timeout, error){
	var v, f = t.timerMap.Load(name)
    var timerInstace timeout

	if !f {
		return timerInstace,errors.New("timer not found")
	}

    timerInstace = v.(timeout)

    return timerInstace,nil
}

func (t *timeout) notifyTimers() {
	for {
		var ti = <-t.timer.C
        log.Println("debug: notification timer")
		t.timerNotification <- ti
	}
}

func newTimeoutPoolImpl() *timeoutPool {
	return &timeoutPool{
		timerMap: sync.Map{},
	}
}
