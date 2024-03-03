package timer

import "time"

func TimerDoor(doorOpenDuration_s float32, doorReset chan bool, doorOpen chan bool) {
	timer := time.NewTimer(time.Duration(doorOpenDuration_s) * time.Second)
	for {
		select {
		case <-timer.C: // Blocking until timer is done
			doorOpen <- false
		case <-doorReset:
			timer.Reset(time.Duration(doorOpenDuration_s) * time.Second)
		}
	}
}
