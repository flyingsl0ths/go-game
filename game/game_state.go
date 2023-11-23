package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	TITLE       uint32 = 0
	PAUSED      uint32 = 1
	GAME        uint32 = 2
	GAME_OVER   uint32 = 3
	HIGH_SCORES uint32 = 4

	TOTAL_PLATFORMS             = 1280 / 64
	GAME_OVER_FONT_SIZE float32 = 100.
	ONE_UP              uint32  = 1000
	PLAYER_HIT_MAX      uint32  = 10
	MAX_POINTS          uint32  = 999999999
)

type State uint32

type Score struct {
	name  string
	score uint32
}

type GameState struct {
	collectables               Spawner
	font                       rl.Font
	gameOverTextAnimationTimer float32
	gameOverTextPos            rl.Vector2
	highScores                 []Score
	objects                    Spawner
	platformHitBoxes           [TOTAL_PLATFORMS]HitBox
	player                     Player
	playerHitCounter           uint32
	playerLives                uint32
	playerOneUpCounter         uint32
	playerPoints               uint32
	spriteSize                 float32
	state                      State
	textures                   TextureAtlas
	windowDimens               WindowDimens
}

func NewGameState(windowDimens [2]float32) GameState {
	textures := NewTextureAtlas("./assets/level.png", "./assets/bg.png", "./assets/hud.png", windowDimens)

	const spriteSize float32 = 64.

	spawnBoundaries := Boundaries{
		bottom: windowDimens[1]/2. + spriteSize*2.0,
		top:    0.0 - textures.food[0].Width,
		width:  windowDimens[0]}

	gameFont := rl.LoadFont("./assets/main.ttf")

	return GameState{
		collectables:               NewSpawner(rl.GetFrameTime()*20, 200., len(textures.food), len(textures.food)/3, spawnBoundaries),
		font:                       gameFont,
		gameOverTextAnimationTimer: 0.,
		gameOverTextPos:            rl.NewVector2((windowDimens[0]/2.)-(3*GAME_OVER_FONT_SIZE), 0-GAME_OVER_FONT_SIZE),
		highScores:                 []Score{},
		objects:                    NewSpawner(rl.GetFrameTime()*5, 300., len(textures.objects), len(textures.objects), spawnBoundaries),
		platformHitBoxes:           makePlatforms(windowDimens[1]/2.+spriteSize*2., spriteSize),
		player:                     NewPlayer("./assets/player.png", rl.NewVector2(100, (windowDimens[1]/2.)+spriteSize+20.), windowDimens[1]+spriteSize+32., spriteSize+32.),
		playerHitCounter:           0,
		playerLives:                1,
		playerOneUpCounter:         0,
		playerPoints:               0,
		spriteSize:                 spriteSize,
		state:                      State(GAME),
		textures:                   textures,
		windowDimens:               WindowDimens{width: windowDimens[0], height: windowDimens[1]},
	}
}

func makePlatforms(yPos float32, platformSize float32) [TOTAL_PLATFORMS]HitBox {
	// WINDOW WIDTH / platformSize
	total := int(1280 / platformSize)

	rs := [TOTAL_PLATFORMS]HitBox{}

	for i := 0; i < total; i++ {
		rs[i] = HitBox{
			box:           rl.NewRectangle(platformSize*float32(i), yPos, platformSize, platformSize),
			damageCounter: NewDamageCounter(25.),
		}
	}

	return rs
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
