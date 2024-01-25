package game

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func RunGame(game *GameState, delta float32) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)

	switch game.state {
	case TITLE:
		onTitleState(game, delta)
	case PAUSED:
		onPauseState(game, delta)
	case GAME:
		onGameState(game, delta)
	case HIGH_SCORE_INPUT:
		onHighScoreInput(game, delta)
		break
	case GAME_OVER:
		onGameOver(game, delta)
		break
	case HIGH_SCORES:
		break
	}

	rl.EndDrawing()
}

func onTitleState(game *GameState, delta float32) {
	updateTitleState(game, delta)
	drawTitleState(game)
}

func updateTitleState(game *GameState, delta float32) {
	UpdateSpawner(&game.titleScreenCollectables, func(item *Item, delta float32) bool {
		item.position.Y += item.gravity * delta

		item.rotation += 1.5

		return true
	}, delta)

	updateTitleScreenButtons(game)
}

func updateTitleScreenButtons(game *GameState) {
	total := len(game.titleScreenButtons.buttons)

	buttons := game.titleScreenButtons.buttons
	for i := 0; i < total; i++ {
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) && rl.CheckCollisionPointRec(rl.GetMousePosition(), buttons[i].position) {
			buttons[i].state = PRESSED
			buttons[game.titleScreenButtons.lastActive].state = NORMAL
			buttons[i].onClick(game)
			game.titleScreenButtons.lastActive = uint(i)
			continue
		} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), buttons[i].position) {
			buttons[i].state = HOVER
			buttons[game.titleScreenButtons.lastActive].state = NORMAL
			game.titleScreenButtons.lastActive = uint(i)
			continue
		}

		buttons[i].state = NORMAL
		game.titleScreenButtons.lastActive = uint(i)
	}
}

func drawTitleState(game *GameState) {
	rl.DrawTexture(game.textures.textureSheets.bg, 0, 0, rl.RayWhite)

	for _, item := range game.titleScreenCollectables.items {
		rl.DrawTexturePro(game.textures.textureSheets.level,
			game.textures.food[item.itemID],
			rl.NewRectangle(item.position.X, item.position.Y, game.spriteSize, game.spriteSize),
			rl.NewVector2(game.spriteSize/2., game.spriteSize/2.), item.rotation, rl.RayWhite)
	}

	rl.DrawTexture(game.textures.textureSheets.titleScreenImage, int32(game.titleScreenImagePos.X), int32(game.titleScreenImagePos.Y), rl.RayWhite)

	drawTitleScreenButtons(game)
}

func drawTitleScreenButtons(game *GameState) {
	const BUTTON_WIDTH float32 = 281.
	const BUTTON_HEIGHT float32 = 103.

	buttonFont := rl.GetFontDefault()

	button := game.titleScreenButtons.buttons[0]
	buttonTextColor := rl.NewColor(255, 80, 88, 255)

	rl.DrawTexturePro(game.textures.textureSheets.buttons, game.textures.buttons[button.state], button.position,
		rl.NewVector2(0., 0), 0., rl.RayWhite)

	rl.DrawTextEx(buttonFont, "  GAME", button.textPosition, GAME_FONT_SIZE, 0., buttonTextColor)

	button = game.titleScreenButtons.buttons[1]

	rl.DrawTexturePro(game.textures.textureSheets.buttons, game.textures.buttons[button.state], button.position,
		rl.NewVector2(0, 0), 0., rl.RayWhite)

	rl.DrawTextEx(buttonFont, "SCORES", button.textPosition, GAME_FONT_SIZE, 0., buttonTextColor)

	button = game.titleScreenButtons.buttons[2]

	rl.DrawTexturePro(game.textures.textureSheets.buttons, game.textures.buttons[button.state], button.position,
		rl.NewVector2(0, 0), 0., rl.RayWhite)

	rl.DrawTextEx(buttonFont, "  EXIT", button.textPosition, GAME_FONT_SIZE, 0., buttonTextColor)
}

func onPauseState(game *GameState, delta float32) {
	updatePauseScreenButtons(game, delta)
	drawGameState(game)
	drawPauseScreenButtons(game, delta)
}

func updatePauseScreenButtons(game *GameState, delta float32) {
	total := len(game.pauseScreenButtons)

	buttons := game.pauseScreenButtons
	for i := 0; i < total; i++ {
		lastModified := (i + 1) % total

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) && rl.CheckCollisionPointRec(rl.GetMousePosition(), buttons[i].position) {
			buttons[i].state = PRESSED
			buttons[lastModified].state = NORMAL
			buttons[i].onClick(game)
			continue
		} else if rl.CheckCollisionPointRec(rl.GetMousePosition(), buttons[i].position) {
			buttons[i].state = HOVER
			buttons[lastModified].state = NORMAL
			continue
		}

		buttons[i].state = NORMAL
	}
}

func drawPauseScreenButtons(game *GameState, delta float32) {
	const BUTTON_WIDTH float32 = 281.
	const BUTTON_HEIGHT float32 = 103.

	buttonFont := rl.GetFontDefault()

	button := game.pauseScreenButtons[0]
	buttonTextColor := rl.NewColor(255, 80, 88, 255)

	rl.DrawTexturePro(game.textures.textureSheets.buttons, game.textures.buttons[button.state], button.position,
		rl.NewVector2(0., 0), 0., rl.RayWhite)

	rl.DrawTextEx(buttonFont, "RESUME", button.textPosition, GAME_FONT_SIZE, 0., buttonTextColor)

	button = game.pauseScreenButtons[1]

	rl.DrawTexturePro(game.textures.textureSheets.buttons, game.textures.buttons[button.state], button.position,
		rl.NewVector2(0, 0), 0., rl.RayWhite)

	rl.DrawTextEx(buttonFont, " TITLE", button.textPosition, GAME_FONT_SIZE, 0., buttonTextColor)
}

func onGameState(game *GameState, delta float32) {
	updateGameState(game, delta)

	drawGameState(game)
}

func onHighScoreInput(game *GameState, delta float32) {
	handleHighScoreInput(game)

	drawGameState(game)

	drawHighScoreInputLayer(game)
}

func handleHighScoreInput(game *GameState) {
	game.textInput = UpdateTextInput(game.textInput)

	if rl.IsKeyReleased(rl.KeyEnter) {
		handleHighScoreEntered(game)
	}
}

func handleHighScoreEntered(game *GameState) {
	encoder, err := NewEncoder[ScoreBoard](game.highScoresFilePath)

	defer encoder.Close()

	if err != nil {
		game.exitGame = true
	}

	postHighScore(game, &encoder)

	game.state = HIGH_SCORES

	game.lastState = HIGH_SCORE_INPUT
}

func postHighScore(game *GameState, encoder *Encoder[ScoreBoard]) {
	newScore := PlayerScore{PlayerName: string(game.textInput.text), Score: game.playerPoints}

	if !game.scoreboard.FirstPosted {
		game.scoreboard.FirstPosted = true
		game.scoreboard.Scores[0] = newScore
	} else {
		newScorePos := HighestScore(game.scoreboard.Scores, newScore)
		game.scoreboard.Scores[newScorePos] = newScore
	}

	if err := encoder.Encode(game.scoreboard); err != nil {
		game.exitGame = true
	}
}

func updateGameState(game *GameState, delta float32) {
	if (rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)) && rl.IsKeyDown(rl.KeySpace) {
		game.lastState = game.state
		game.state = PAUSED
		return
	}

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
	collectables := func(item *Item, delta float32) bool {
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

	objects := func(item *Item, delta float32) bool {
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
				game.state = HIGH_SCORE_INPUT
			}
			return false
		}

		game.player.wasHit = true
		return true
	}

	UpdateSpawner(&game.objects, objects, delta)

	UpdateSpawner(&game.collectables, collectables, delta)
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
	rl.DrawTexture(game.textures.textureSheets.bg, 0, 0, rl.RayWhite)

	DrawHUD(game.textures, game.windowDimens, game.playerLives, game.playerPoints)

	drawSpawnedObjects(game)

	DrawPlayer(game.player)

	DrawPlatforms(game.textures, game.windowDimens, game.spriteSize, game.platformHitBoxes)
}

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
		rl.NewVector2(0., 0.), 0., rl.RayWhite)

	lifeIconPos.X += lifeIconSize

	multiplierIconSize := textureAtlas.hud[2].Width * 2.

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[2],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y+(multiplierIconSize/2.), multiplierIconSize, multiplierIconSize),
		rl.NewVector2(0., 0.), 0., rl.RayWhite)

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
		rl.NewVector2(0., 0.), 0., rl.RayWhite)

	lifeIconPos.X += numberIconSize + PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[3+secondNumOffset],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y, numberIconSize, numberIconSize),
		rl.NewVector2(0., 0.), 0., rl.RayWhite)
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
		rl.NewVector2(0., 0.), 0., rl.RayWhite)

	pointsIconPos.X += pointsIconSize + IMAGE_PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		multiplierIcon,
		rl.NewRectangle(pointsIconPos.X, yOffset, multiplierIconSize, multiplierIconSize),
		rl.NewVector2(0., 0.), 0., rl.RayWhite)

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
			rl.NewVector2(0., 0.), 0., rl.RayWhite)

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
			rl.NewVector2(0., 0.), 0., rl.RayWhite)

		playerPointsStride += 1
	}
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
				rl.NewVector2(0., 0.), 0., rl.RayWhite)
		} else {
			continue
		}

		if hitBoxes[i].damageCounter.percentage >= 75. {
			rl.DrawTexturePro(textureAtlas.textureSheets.level,
				textureAtlas.overlays[1],
				platformRect,
				rl.NewVector2(0., 0.), 0., rl.RayWhite)
		} else if hitBoxes[i].damageCounter.percentage >= 45. {
			rl.DrawTexturePro(textureAtlas.textureSheets.level,
				textureAtlas.overlays[0],
				platformRect,
				rl.NewVector2(0., 0.), 0., rl.RayWhite)
		}
	}
}

func drawSpawnedObjects(game *GameState) {
	for _, item := range game.collectables.items {
		rl.DrawTexturePro(game.textures.textureSheets.level,
			game.textures.food[item.itemID],
			rl.NewRectangle(item.position.X, item.position.Y, game.spriteSize, game.spriteSize),
			rl.NewVector2(game.spriteSize/2., game.spriteSize/2.), item.rotation, rl.RayWhite)
	}

	for _, item := range game.objects.items {
		rl.DrawTexturePro(game.textures.textureSheets.level,
			game.textures.objects[item.itemID],
			rl.NewRectangle(item.position.X, item.position.Y, game.spriteSize, game.spriteSize),
			rl.NewVector2(game.spriteSize/2., game.spriteSize/2.), item.rotation, rl.RayWhite)
	}

}

func onGameOver(game *GameState, delta float32) {
	updateGameOverState(game, delta)

	drawGameOverState(game)
}

func updateGameOverState(game *GameState, delta float32) {
	game.gameOverTextAnimationTimer += 0.0040

	game.gameOverTextPos.Y += 10. * ElasticEaseOut(game.gameOverTextAnimationTimer, 0., 1., 0.5)

	if game.gameOverTextAnimationTimer >= 0.75 {
		game.state = HIGH_SCORE_INPUT
		game.lastState = GAME_OVER
	}
}

func drawGameOverState(game *GameState) {
	rl.DrawTexture(game.textures.textureSheets.bg, 0, 0, rl.RayWhite)

	DrawHUD(game.textures, game.windowDimens, game.playerLives, game.playerPoints)

	drawSpawnedObjects(game)

	rl.DrawTextEx(game.font, "GAME OVER", game.gameOverTextPos, GAME_OVER_FONT_SIZE, 0., rl.Red)

	DrawPlayer(game.player)

	DrawPlatforms(game.textures, game.windowDimens, game.spriteSize, game.platformHitBoxes)
}

func drawHighScoreInputLayer(game *GameState) {
	rl.DrawTextEx(game.font, "ENTER YOUR NAME", game.highScoreNameBannerPos, GAME_FONT_SIZE, 0., rl.RayWhite)

	DrawTextInput(game.textInput)
}
