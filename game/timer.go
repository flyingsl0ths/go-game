package game

type Timer struct {
	tick     float32
	duration float32
	finished bool
	loops    bool
}

func NewTimer(duration float32, loops bool) Timer {
	return Timer{
		tick:     0.0,
		duration: duration,
		finished: false,
		loops:    loops,
	}
}

func Tick(timer Timer, dt float32) Timer {
	if !timer.finished {
		timer.tick += dt
	}

	if timer.tick >= timer.duration {
		timer.finished = true
		if timer.loops {
			timer.tick -= timer.duration
			timer.finished = false
		}
	}

	return timer
}
