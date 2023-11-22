package game

import rl "github.com/gen2brain/raylib-go/raylib"

type HitBox struct {
	box           rl.Rectangle
	damageCounter DamageCounter
}
