package game

type Physics[T any] struct {
	bottom     float32
	gravity    float32
	ground     float32
	jumpHeight float32
	velocity   T
}
