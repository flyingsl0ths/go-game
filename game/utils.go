package game

import "math/rand"

type WindowDimens struct {
	width  float32
	height float32
}

func RandRange(min int, max int) int {
	if min > max {
		return 0
	} else {
		return rand.Intn(max-min) + min
	}
}
