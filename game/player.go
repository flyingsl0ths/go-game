package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Player struct {
	animation    LinearFrameAnimation
	fell         bool
	halt         bool
	isJumping    bool
	isMoving     bool
	originalSize float32
	physics      Physics[rl.Vector2]
	position     rl.Vector2
	stunTimer    Timer
	texture      rl.Texture2D
	textureBox   rl.Rectangle
	textureSize  float32
	wasHit       bool
}

func NewPlayer(spriteSheetPath string, startPosition rl.Vector2, ground float32, playerSize float32) Player {
	spriteSheet := rl.LoadImage(spriteSheetPath)

	defer rl.UnloadImage(spriteSheet)

	jumpHeight := float32(300.0)

	return Player{
		animation:    NewAnimation(1.0, true, 5),
		fell:         false,
		halt:         false,
		stunTimer:    stunTimer(),
		isMoving:     false,
		isJumping:    false,
		originalSize: 32.,
		position:     startPosition,
		texture:      rl.LoadTextureFromImage(spriteSheet),
		textureBox:   rl.NewRectangle(0, 0, 32., 32.),
		textureSize:  playerSize,
		physics: Physics[rl.Vector2]{
			bottom:     ground,
			gravity:    -500,
			ground:     startPosition.Y,
			jumpHeight: -jumpHeight,
			velocity:   rl.NewVector2(200., 0.)},
		wasHit: false,
	}
}

func stunTimer() Timer {
	return NewTimer(0.75, false)
}

func PlayerBoundingBox(player *Player) rl.Rectangle {
	return rl.NewRectangle(player.position.X, player.position.Y, player.textureSize, player.textureSize)
}

func UpdatePlayer(player Player, delta float32) Player {
	return handlePhysics(handleCollision(handleSpriteChange(player, delta), delta), delta)
}

func DrawPlayer(player Player) {
	dest := rl.NewRectangle(player.position.X, player.position.Y, player.textureSize, player.textureSize)

	rl.DrawTexturePro(player.texture, player.textureBox, dest,
		rl.NewVector2(player.textureSize/2., player.textureSize/2.), 0., rl.White)
}

func handleSpriteChange(player Player, delta float32) Player {
	player_ := player

	if player_.fell {
		player_.textureBox.X = player_.originalSize * 8.
		return player_
	}

	if player_.wasHit {
		player_.textureBox.X = player_.originalSize * 7.
		return player_
	}

	if player_.isJumping {
		player_.textureBox.X = player_.originalSize * 6.
		return player_
	}

	if !player_.isMoving {
		player_.textureBox.X = player_.originalSize * 5.
		return player_
	}

	return handleSpriteAnimation(player_, delta)
}

func handleSpriteAnimation(player Player, delta float32) Player {
	player_ := player

	player_.animation = UpdateAnimation(player_.animation, delta)

	currentFrame := NextFrame(player_.animation)

	player_.textureBox.X = float32(currentFrame) * player_.originalSize

	return player_
}

func handleCollision(player Player, delta float32) Player {
	player_ := player

	if player_.wasHit {
		return onTick(player_, delta)
	} else {
		return handleMovement(player_, delta)
	}
}

func onTick(player Player, delta float32) Player {
	player_ := player

	player_.stunTimer = Tick(player_.stunTimer, delta)

	if player_.stunTimer.finished {
		player_.stunTimer = stunTimer()
		player_.wasHit = false
	}

	return player_
}

func handleMovement(player Player, delta float32) Player {
	if player.halt {
		return player
	}

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
	player_ := player

	if player_.fell {
		player_.position.Y += -player_.physics.jumpHeight * delta
		return player_
	}

	if !player.isJumping {
		return player_
	}

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
