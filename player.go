package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	animation    LinearAnimation
	isMoving     bool
	originalSize float32
	playerSize   float32
	position     rl.Vector2
	texture      rl.Texture2D
	textureBox   rl.Rectangle
	textureSize  float32
	physics      Physics[rl.Vector2]
}

func NewPlayer(image *rl.Image, startPosition rl.Vector2, playerSize float32) Player {
	frameCount := 5

	return Player{
		animation: LinearAnimation{
			duration:    float32(1.),
			elapsedTime: float32(0.),
			frames:      uint32(frameCount),
		},
		isMoving:     false,
		originalSize: 32.,
		playerSize:   playerSize,
		position:     startPosition,
		texture:      rl.LoadTextureFromImage(image),
		textureBox:   rl.NewRectangle(0, 0, 32., 32.),
		textureSize:  float32(playerSize),
		physics: Physics[rl.Vector2]{
			gravity: -500, ground: startPosition.Y, jumpHeight: -300, velocity: rl.NewVector2(150., 0.)},
	}
}

func UpdatePlayer(player Player, delta float32) Player {
	return handlePhysics(handleMovement(handleSpriteAnimation(player, delta), delta), delta)
}

func DrawPlayer(player Player) {
	dest := rl.NewRectangle(player.position.X, player.position.Y, player.textureSize, player.textureSize)

	rl.DrawTexturePro(player.texture, player.textureBox, dest,
		rl.NewVector2(player.textureSize/2., player.textureSize/2.), 0., rl.White)
}

func handleMovement(player Player, delta float32) Player {
	player_ := player

	player_.isMoving = false

	difference := float32(0.)
	movement := player_.physics.velocity.X * delta

	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		difference = -movement
		player_.textureBox.Width = -player.originalSize
	}

	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		difference = movement
		player_.textureBox.Width = player.originalSize
	}

	if rl.IsKeyDown(rl.KeySpace) && player_.physics.velocity.Y == 0. {
		player_.physics.velocity.Y = player_.physics.jumpHeight
	}

	player_.position.X += difference

	if difference != 0. {
		player_.isMoving = true
	}

	return player_
}

func handlePhysics(player Player, delta float32) Player {
	player_ := player

	if player_.physics.velocity.Y != 0. {
		player_.position.Y += player_.physics.velocity.Y * delta
		player_.physics.velocity.Y = player_.physics.velocity.Y - player_.physics.gravity*delta
	}

	if player_.position.Y > player_.physics.ground {
		player_.physics.velocity.Y = 0.
		player_.position.Y = player_.physics.ground
	}

	return player_
}

func handleSpriteAnimation(player Player, delta float32) Player {
	player_ := player

	if !player.isMoving {
		player_.textureBox.X = 32. * 5.
		return player_
	}

	player_.animation = UpdateAnimation(player_.animation, delta)

	currentFrame := NextFrame(player_.animation)

	player_.textureBox.X = float32(currentFrame) * player_.originalSize

	return player_
}
