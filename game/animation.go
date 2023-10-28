package game

import "math"

type LinearFrameAnimation struct {
	timer  Timer
	frames uint32
}

func NewAnimation(duration float32, loops bool, frames uint32) LinearFrameAnimation {
	return LinearFrameAnimation{
		timer:  NewTimer(duration, loops),
		frames: frames,
	}
}

func UpdateAnimation(animation LinearFrameAnimation, dt float32) LinearFrameAnimation {
	anim := animation

	anim.timer = Tick(anim.timer, dt)

	return anim
}

func NextFrame(animation LinearFrameAnimation) uint32 {
	return uint32(math.Floor(float64(animation.timer.tick) / float64(animation.timer.duration) * float64(animation.frames)))
}

func ElasticEaseOut(time float32, begin float32, changeInValue float32, duration float32) float32 {
	if time == 0 {
		return begin
	}

	if (time / duration) == 1 {
		return begin + changeInValue
	}

	p := duration * 3.

	a := changeInValue

	s := p / 4

	return 1 - (a*float32(math.Pow(2., float64(-10.0*time)))*float32(math.Sin(float64((time*duration-s)*(2*math.Pi)/p))) + changeInValue + begin)
}
