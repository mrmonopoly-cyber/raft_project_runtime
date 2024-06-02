package timeout

import "time"

type TimeoutPool interface{
    AddTimeout(name string, duration time.Duration)
    GetTimeoutNotifycationChan(name string) (<- chan time.Time, error)
    RestartTimeout(name string) error
    StopTimeout(name string) error
}

func NewTimeoutPool() TimeoutPool{
    return newTimeoutPoolImpl()
}
