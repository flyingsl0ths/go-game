package game

import "math"

type LinearAnimation struct {
	timer  Timer
	frames uint32
}

func NewAnimation(duration float32, loops bool, frames uint32) LinearAnimation {
	return LinearAnimation{
		timer:  NewTimer(duration, loops),
		frames: frames,
	}
}

func UpdateAnimation(animation LinearAnimation, dt float32) LinearAnimation {
	anim := animation

	anim.timer = Tick(anim.timer, dt)

	return anim
}

func NextFrame(animation LinearAnimation) uint32 {
	return uint32(math.Floor(float64(animation.timer.tick) / float64(animation.timer.duration) * float64(animation.frames)))
}
