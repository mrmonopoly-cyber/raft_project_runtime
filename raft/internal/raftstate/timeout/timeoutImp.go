package timeout

import (
	"errors"
	"sync"
	"time"
)

type timeout struct {
	timer             *time.Timer
	duration          time.Duration
	timerNotification chan time.Time
}

type timeoutPool struct {
	timerMap sync.Map
}


// AddTimeout implements TimeoutPool.
func (t *timeoutPool) AddTimeout(name string, duration time.Duration) {
	var newTimer = timeout{
		timer:             nil,
		duration:          duration,
		timerNotification: make(chan time.Time),
	}

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
	switch timerInstance.timer.C{
	case nil:
		timerInstance.timer = time.NewTimer(timerInstance.duration)
	default:
		timerInstance.timer.Reset(timerInstance.duration)
	}

	go timerInstance.notifyTimers()
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
		t.timerNotification <- ti
	}
}

func newTimeoutPoolImpl() *timeoutPool {
	return &timeoutPool{
		timerMap: sync.Map{},
	}
}
