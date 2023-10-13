package game

type Physics[T any] struct {
	gravity    float32
	ground     float32
	jumpHeight float32
	velocity   T
}
