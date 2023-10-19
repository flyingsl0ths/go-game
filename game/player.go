package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	animation    LinearAnimation
	isMoving     bool
	isJumping    bool
	originalSize float32
	playerSize   float32
	position     rl.Vector2
	texture      rl.Texture2D
	textureBox   rl.Rectangle
	textureSize  float32
	physics      Physics[rl.Vector2]
}

func NewPlayer(spriteSheetPath string, startPosition rl.Vector2, playerSize float32) Player {
	spriteSheet := rl.LoadImage(spriteSheetPath)
	defer rl.UnloadImage(spriteSheet)

	return Player{
		animation: LinearAnimation{
			timer:  NewTimer(1.0, true),
			frames: 5,
		},
		isMoving:     false,
		isJumping:    false,
		originalSize: 32.,
		playerSize:   playerSize,
		position:     startPosition,
		texture:      rl.LoadTextureFromImage(spriteSheet),
		textureBox:   rl.NewRectangle(0, 0, 32., 32.),
		textureSize:  playerSize,
		physics: Physics[rl.Vector2]{
			gravity: -500, ground: startPosition.Y, jumpHeight: -300, velocity: rl.NewVector2(200., 0.)},
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
		player_.isJumping = true
	}

	player_.position.X += difference

	if difference != 0. {
		player_.isMoving = true
	}

	return player_
}

func handlePhysics(player Player, delta float32) Player {
	if !player.isJumping {
		return player
	}

	player_ := player

	if player_.physics.velocity.Y != 0. {
		player_.position.Y += player_.physics.velocity.Y * delta
		player_.physics.velocity.Y = player_.physics.velocity.Y - player_.physics.gravity*delta
	}

	if player_.position.Y > player_.physics.ground {
		player_.physics.velocity.Y = 0.
		player_.position.Y = player_.physics.ground
		player_.isJumping = false
	}

	return player_
}

func handleSpriteAnimation(player Player, delta float32) Player {
	player_ := player

	if player_.isJumping {
		player_.textureBox.X = player_.originalSize * 6.
		return player_
	}

	if !player_.isMoving {
		player_.textureBox.X = player_.originalSize * 5.
		return player_
	}

	player_.animation = UpdateAnimation(player_.animation, delta)

	currentFrame := NextFrame(player_.animation)

	player_.textureBox.X = float32(currentFrame) * player_.originalSize

	return player_
}
