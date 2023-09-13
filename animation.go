package main

import "math"

type LinearAnimation struct {
	duration    float32
	elapsedTime float32
	frames      uint32
}

func UpdateAnimation(animation LinearAnimation, dt float32) LinearAnimation {
	anim := animation

	anim.elapsedTime += dt

	if anim.elapsedTime >= anim.duration {
		anim.elapsedTime -= anim.duration
	}

	return anim
}

func NextFrame(animation LinearAnimation) uint32 {
	return uint32(math.Floor(float64(animation.elapsedTime) / float64(animation.duration) * float64(animation.frames)))
}
