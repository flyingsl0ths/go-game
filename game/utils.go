package game

import (
	"math"
	"math/rand"
)

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

func EaseOutCirc(x float32) float32 {
	return float32(math.Sqrt(1.0 - math.Pow(float64(x)-1.0, 2.0)))
}
