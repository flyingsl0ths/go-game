package game

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawHUD(textureAtlas TextureAtlas, windowDimens WindowDimens, lives uint32, points uint32) {
	drawPlayerIcons(textureAtlas, lives)
	drawCollectableIcons(textureAtlas, windowDimens, points)
}

func drawPlayerIcons(textureAtlas TextureAtlas, lives uint32) {
	const PADDING float32 = 5.
	lifeIconPos := rl.NewVector2(PADDING, PADDING)

	lifeIcon := textureAtlas.hud[0]
	lifeIconSize := lifeIcon.Width * 2

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		lifeIcon,
		rl.NewRectangle(0., 0., lifeIconSize+2., lifeIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	lifeIconPos.X += lifeIconSize

	multiplierIconSize := textureAtlas.hud[2].Width * 2.

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[2],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y+(multiplierIconSize/2.), multiplierIconSize, multiplierIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	lifeIconPos.X += multiplierIconSize + PADDING
	numberIconSize := textureAtlas.hud[3].Width * 2
	lifeIconPos.Y += numberIconSize/2. + PADDING/2.

	playerLives := strconv.FormatUint(uint64(lives), 10)

	firstNumOffset := 0
	secondNumOffset := 0

	if lives >= 10 {
		firstNumOffset = int(playerLives[0] - '0')
		secondNumOffset = int(playerLives[1] - '0')
	} else {
		secondNumOffset = int(playerLives[0] - '0')
	}

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[3+firstNumOffset],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y, numberIconSize, numberIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	lifeIconPos.X += numberIconSize + PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[3+secondNumOffset],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y, numberIconSize, numberIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)
}

func drawCollectableIcons(textureAtlas TextureAtlas, windowDimens WindowDimens, points uint32) {
	const IMAGE_PADDING float32 = 6.

	pointsIcon := textureAtlas.hud[1]
	multiplierIcon := textureAtlas.hud[2]

	windowWidth := windowDimens.width
	multiplierIconSize := multiplierIcon.Width * 2
	pointsIconSize := pointsIcon.Width * 2

	offset := pointsIconSize + multiplierIconSize + (textureAtlas.hud[3].Width * 2 * 10)

	pointsIconPos := rl.NewVector2(windowWidth-offset, IMAGE_PADDING)

	yOffset := pointsIconPos.Y + IMAGE_PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		pointsIcon,
		rl.NewRectangle(pointsIconPos.X, pointsIconPos.Y, pointsIcon.Width*2., pointsIcon.Height*2),
		rl.NewVector2(0., 0.), 0., rl.White)

	pointsIconPos.X += pointsIconSize + IMAGE_PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		multiplierIcon,
		rl.NewRectangle(pointsIconPos.X, yOffset, multiplierIconSize, multiplierIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	pointsIconPos.X += multiplierIconSize

	const MAX_POINTS = 9
	const ZERO_IMAGE_INDEX = 3

	playerPoints := strconv.FormatUint(uint64(points), 10)

	numberIconSize := textureAtlas.hud[ZERO_IMAGE_INDEX].Width * 2
	padding := MAX_POINTS - len(playerPoints)
	playerPointsStride := 0

	for i := 0; i < padding; i++ {
		numberIcon := textureAtlas.hud[ZERO_IMAGE_INDEX]

		rl.DrawTexturePro(textureAtlas.textureSheets.hud,
			numberIcon,
			rl.NewRectangle(pointsIconPos.X+(numberIconSize*float32(i))+IMAGE_PADDING,
				yOffset, numberIconSize, numberIconSize),
			rl.NewVector2(0., 0.), 0., rl.White)

		playerPointsStride = i
	}

	if padding > 0 {
		playerPointsStride += 1
	}

	for _, c := range playerPoints {
		numberIcon := textureAtlas.hud[ZERO_IMAGE_INDEX+int(c-'0')]

		rl.DrawTexturePro(textureAtlas.textureSheets.hud,
			numberIcon,
			rl.NewRectangle(pointsIconPos.X+(numberIconSize*float32(playerPointsStride))+IMAGE_PADDING,
				yOffset, numberIconSize, numberIconSize),
			rl.NewVector2(0., 0.), 0., rl.White)

		playerPointsStride += 1
	}
}

func RunGame(game *GameState, delta float32) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	switch game.state {
	case State(TITLE):
		break
	case State(PAUSED):
		break
	case State(GAME):
		onGameState(game, delta)
		break
	case State(GAME_OVER):
		onGameOver(game, delta)
		break
	case State(HIGH_SCORES):
		break
	}

	rl.EndDrawing()
}

func onGameState(game *GameState, delta float32) {
	updateGameState(game, delta)

	drawGameState(game)
}

func updateGameState(game *GameState, delta float32) {
	updatePlayer(game, delta)

	updateSpawners(game, delta)
	updatePlatformHitBoxes(game)
}

func updatePlayer(game *GameState, delta float32) {
	game.player = UpdatePlayer(game.player, delta)

	confinePlayer(game)

	if game.player.position.Y > game.windowDimens.height+game.spriteSize+32. {
		game.state = State(GAME_OVER)
		return
	}

	if game.player.halt {
		return
	}

	handleBrokenPlatforms(game)
}

func confinePlayer(game *GameState) {
	const OFFSET = 3
	xMin := game.player.originalSize - OFFSET
	xMax := game.windowDimens.width - (game.player.originalSize - OFFSET)

	if game.player.position.X < xMin {
		game.player.position.X = xMin
	} else if game.player.position.X > xMax {
		game.player.position.X = xMax
	}
}

func handleBrokenPlatforms(game *GameState) {
	platformIndex := CellFrom(game.player.position.X, game.spriteSize, TOTAL_PLATFORMS)

	platform := game.platformHitBoxes[platformIndex]

	if !game.player.isJumping && MaxDamage(platform.damageCounter) {
		game.player.fell = true

		platform = game.platformHitBoxes[(platformIndex+1)%TOTAL_PLATFORMS]

		side := CollideWithSides(PlayerBoundingBox(&game.player), platform.box)

		if side == "bottom" {
			game.player.halt = true
		}
	}
}

func updateSpawners(game *GameState, delta float32) {
	collectablesMover := func(item *Item, delta float32) bool {
		spawnMover(game, item, delta)

		if item.collided && !game.player.wasHit {
			if game.playerPoints < MAX_POINTS {
				total := game.playerPoints + uint32(item.itemID)

				if total > MAX_POINTS {
					total = MAX_POINTS
				}

				game.playerPoints = total
				game.playerOneUpCounter = total

				if game.playerOneUpCounter >= ONE_UP {
					game.playerOneUpCounter = 0
					game.playerLives += 1
				}

			}
		}

		return true
	}

	objectMover := func(item *Item, delta float32) bool {
		spawnMover(game, item, delta)

		if !item.collided {
			return true
		}

		if game.playerPoints > 0 {
			itemDamage := uint32(item.itemID)
			game.playerPoints -= itemDamage
			game.playerOneUpCounter -= itemDamage
		}

		game.playerHitCounter = (game.playerHitCounter + 1) % PLAYER_HIT_MAX

		if game.playerHitCounter == 0 {
			game.playerLives -= 1
			if game.playerLives == 0 {
				game.state = State(GAME_OVER)
			}
			return false
		}

		game.player.wasHit = true
		return true
	}

	UpdateSpawner(&game.objects, objectMover, delta)

	UpdateSpawner(&game.collectables, collectablesMover, delta)
}

func spawnMover(game *GameState, item *Item, delta float32) {
	const ITEM_WIDTH float32 = 16.0
	const ITEM_HEIGHT float32 = 15.0

	item.position.Y += item.gravity * delta

	item.rotation += 1.5

	item.collided = rl.CheckCollisionRecs(rl.NewRectangle(item.position.X, item.position.Y, ITEM_WIDTH, ITEM_HEIGHT),
		PlayerBoundingBox(&game.player))
}

func updatePlatformHitBoxes(game *GameState) {
	const ITEM_WIDTH float32 = 16.0
	const ITEM_HEIGHT float32 = 15.0

	calculateDamage := func(item Item, game *GameState) {
		platformIndex := CellFrom(item.position.X, game.spriteSize, TOTAL_PLATFORMS)

		platform := game.platformHitBoxes[platformIndex]

		if !item.collided && rl.CheckCollisionRecs(rl.NewRectangle(item.position.X, item.position.Y, ITEM_WIDTH, ITEM_HEIGHT), platform.box) {
			platform.damageCounter = DamageCalc(platform.damageCounter, 0.5)

			game.platformHitBoxes[platformIndex] = platform
		}
	}

	for _, item := range game.objects.items {
		calculateDamage(item, game)
	}
}

func drawGameState(game *GameState) {
	rl.DrawTexture(game.textures.textureSheets.bg, 0, 0, rl.White)

	DrawHUD(game.textures, game.windowDimens, game.playerLives, game.playerPoints)

	drawSpawnedObjects(game)

	DrawPlayer(game.player)

	DrawPlatforms(game.textures, game.windowDimens, game.spriteSize, game.platformHitBoxes)
}

func drawSpawnedObjects(game *GameState) {
	for _, item := range game.collectables.items {
		rl.DrawTexturePro(game.textures.textureSheets.level,
			game.textures.food[item.itemID],
			rl.NewRectangle(item.position.X, item.position.Y, game.spriteSize, game.spriteSize),
			rl.NewVector2(game.spriteSize/2., game.spriteSize/2.), item.rotation, rl.White)
	}

	for _, item := range game.objects.items {
		rl.DrawTexturePro(game.textures.textureSheets.level,
			game.textures.objects[item.itemID],
			rl.NewRectangle(item.position.X, item.position.Y, game.spriteSize, game.spriteSize),
			rl.NewVector2(game.spriteSize/2., game.spriteSize/2.), item.rotation, rl.White)
	}

}

func onGameOver(game *GameState, delta float32) {
	updateGameOverState(game, delta)

	drawGameOverState(game, delta)
}

func updateGameOverState(game *GameState, delta float32) {
	game.gameOverTextAnimationTimer += 0.0039

	game.gameOverTextPos.Y += 10. * ElasticEaseOut(game.gameOverTextAnimationTimer, 0., 1., 0.5)

	if game.gameOverTextAnimationTimer >= 0.75 {
		game.state = State(HIGH_SCORES)
	}
}

func drawGameOverState(game *GameState, delta float32) {
	rl.DrawTexture(game.textures.textureSheets.bg, 0, 0, rl.White)

	DrawHUD(game.textures, game.windowDimens, game.playerLives, game.playerPoints)

	drawSpawnedObjects(game)

	rl.DrawTextEx(game.font, "GAME OVER", game.gameOverTextPos, GAME_OVER_FONT_SIZE, 0., rl.Red)

	DrawPlayer(game.player)

	DrawPlatforms(game.textures, game.windowDimens, game.spriteSize, game.platformHitBoxes)
}

func DrawPlatforms(textureAtlas TextureAtlas, windowDimens WindowDimens, textureSize float32, hitBoxes [20]HitBox) {
	const AT_EDGE = 2
	const REGULAR_TILE = 3

	startPosition := windowDimens.height/2. + textureSize*2
	toDraw := len(hitBoxes)
	remaining := toDraw - (AT_EDGE + 1)
	total := len(textureAtlas.platforms)

	for i := 0; i < toDraw; i++ {
		tile := REGULAR_TILE

		if i <= AT_EDGE {
			tile = i
		} else if i >= remaining {
			tile = total - (toDraw % i)
		}

		platformRect := rl.NewRectangle(textureSize*float32(i), startPosition, textureSize, textureSize)

		if !MaxDamage(hitBoxes[i].damageCounter) {
			rl.DrawTexturePro(textureAtlas.textureSheets.level,
				textureAtlas.platforms[tile],
				platformRect,
				rl.NewVector2(0., 0.), 0., rl.White)
		} else {
			continue
		}

		if hitBoxes[i].damageCounter.percentage >= 75. {
			rl.DrawTexturePro(textureAtlas.textureSheets.level,
				textureAtlas.overlays[1],
				platformRect,
				rl.NewVector2(0., 0.), 0., rl.White)
		} else if hitBoxes[i].damageCounter.percentage >= 45. {
			rl.DrawTexturePro(textureAtlas.textureSheets.level,
				textureAtlas.overlays[0],
				platformRect,
				rl.NewVector2(0., 0.), 0., rl.White)
		}
	}
}
