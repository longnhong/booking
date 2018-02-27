package common

import "time"

type dailyTimer struct {
	Timer
	at int64
}

var day int64 = 24 * 3600

func NewDailyTimer(at string, handler func()) (*dailyTimer, error) {
	var loc, _ = time.LoadLocation("Asia/Ho_Chi_Minh")
	var t, err = time.ParseInLocation("15:04", at, loc)
	if err != nil {
		return nil, err
	}
	var d = &dailyTimer{
		Timer: Timer{handler: handler},
	}
	d.at = t.Unix()

	return d, nil
}

func (t *dailyTimer) delay() int64 {
	// t.at is at year 0000
	return (t.at-time.Now().Unix())%day + day
}

func (t *dailyTimer) Start() {
	var delay = time.Second * time.Duration(t.delay())
	t.Schedule(delay, 24*time.Hour)
}

type Timer struct {
	timer   *time.Timer
	handler func()
}

func NewTimer(handler func()) *Timer {
	return &Timer{
		handler: handler,
	}
}

func (e *Timer) Schedule(delay time.Duration, loop time.Duration) {
	e.timer = time.AfterFunc(delay, func() {
		if loop > 0 {
			e.Schedule(loop, loop)
		}
		e.handler()
	})
}

func (e *Timer) Cancel() {
	e.timer.Stop()
}

func SleepTask(seconds string) {
	dur, _ := time.ParseDuration(seconds)
	time.Sleep(dur)
}

func Round2(x, unit float32) float32 {
	if x > 0 {
		return float32(int64(x/unit+0.5)) * unit
	}
	return float32(int64(x/unit-0.5)) * unit

}
