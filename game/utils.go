package game

import (
	"math"
	"math/rand"
)

type WindowDimens struct {
	width  float32
	height float32
}

type Lives [2]rune
type Points [9]rune

func removeFrom[T any](xs []T, i int) []T {
	return append(xs[:i], xs[i+1:]...)
}

func RandRange(min int, max int) int {
	if min > max {
		return 0
	} else {
		return rand.Intn(max-min) + min
	}
}

func EaseOutCirc(x float64) float64 {
	return math.Sqrt(1.0 - math.Pow(x-1.0, 2.0))
}
