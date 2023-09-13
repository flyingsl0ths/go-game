package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	animation   LinearAnimation
	position    rl.Vector2
	textureBox  rl.Rectangle
	texture     rl.Texture2D
	textureSize float32
}

func NewPlayer(image *rl.Image) Player {
	playerSize := 96
	frameCount := 5

	rl.ImageResizeNN(image, int32(playerSize)*int32(frameCount), int32(playerSize))

	return Player{
		textureSize: float32(playerSize),
		texture:     rl.LoadTextureFromImage(image),
		textureBox:  rl.NewRectangle(0, 0, float32(playerSize), float32(playerSize)),
		position:    rl.Vector2{X: 50, Y: 50},
		animation: LinearAnimation{
			duration:    float32(1.),
			elapsedTime: float32(0.),
			frames:      uint32(frameCount),
		},
	}
}

func UpdatePlayer(player Player) Player {
	player_ := player

	player_.animation = UpdateAnimation(player_.animation, rl.GetFrameTime())

	currentFrame := NextFrame(player_.animation)

	if currentFrame == 0 {
		player_.textureBox.X = 0
	} else {
		player_.textureBox.X = float32(currentFrame) * player_.textureSize
	}

	return player_
}

func DrawPlayer(player Player) {
	rl.DrawTextureRec(player.texture, player.textureBox, player.position, rl.White)
}