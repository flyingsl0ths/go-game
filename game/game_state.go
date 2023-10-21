package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	TITLE       uint32 = 0
	PAUSED      uint32 = 1
	GAME        uint32 = 2
	HIGH_SCORES uint32 = 3
	MAX_POINTS  uint32 = 999999999
)

type State uint32

type Score struct {
	name  string
	score uint32
}

type GameState struct {
	collectables Spawner
	highScores   []Score
	objects      Spawner
	playerLives  uint32
	player       Player
	playerPoints uint32
	powerUps     Spawner
	spriteSize   float32
	state        State
	textures     TextureAtlas
	windowDimens WindowDimens
}

func NewGameState(windowDimens [2]float32) GameState {
	textures := NewTextureAtlas("./assets/level.png", "./assets/bg.png", "./assets/hud.png", windowDimens)

	const spriteSize float32 = 64.

	spawnBoundaries := Boundaries{
		bottom: windowDimens[1]/2. + spriteSize*2.0,
		top:    0.0 - textures.food[0].Width,
		width:  windowDimens[0]}

	return GameState{
		collectables: NewSpawner(rl.GetFrameTime()*20, 200., len(textures.food), len(textures.food)/3, spawnBoundaries),
		highScores:   []Score{},
		objects:      NewSpawner(rl.GetFrameTime()*5, 300., len(textures.objects), len(textures.objects), spawnBoundaries),
		playerLives:  1,
		player:       NewPlayer("./assets/player.png", rl.NewVector2(50., (windowDimens[1]/2.)+spriteSize+20.), spriteSize+32.),
		playerPoints: 0,
		spriteSize:   spriteSize,
		state:        State(GAME),
		textures:     textures,
		windowDimens: WindowDimens{width: windowDimens[0], height: windowDimens[1]},
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
	game.player = UpdatePlayer(game.player, delta)

	updateSpawners(game, delta)
}

func updateSpawners(game *GameState, delta float32) {
	collectablesMover := func(item *Item, delta float32) {
		spawnMover(game, item, delta)

		if item.collided && !game.player.wasHit {
			if game.playerPoints < MAX_POINTS {
				total := game.playerPoints + uint32(item.itemID)
				if total > MAX_POINTS {
					total = MAX_POINTS
				}
				game.playerPoints = total

			}
		}
	}

	objectMover := func(item *Item, delta float32) {
		spawnMover(game, item, delta)

		if item.collided {
			if game.playerPoints > 0 {
				game.playerPoints -= uint32(item.itemID)
			}

			game.player.wasHit = true
		}
	}

	UpdateSpawner(&game.objects, objectMover, delta)

	UpdateSpawner(&game.collectables, collectablesMover, delta)
}

func spawnMover(game *GameState, item *Item, delta float32) {
	const ITEM_WIDTH float32 = 16.0
	const ITEM_HEIGHT float32 = 15.0

	item.position.Y += item.gravity * delta

	item.rotation += EaseOutCirc(0.2)

	rl.CheckCollisionRecs(rl.NewRectangle(item.position.X, item.position.Y, ITEM_WIDTH, ITEM_HEIGHT),
		PlayerBoundingBox(&game.player))

	item.collided = rl.CheckCollisionRecs(rl.NewRectangle(item.position.X, item.position.Y, ITEM_WIDTH, ITEM_HEIGHT),
		PlayerBoundingBox(&game.player))
}

func drawGameState(game *GameState) {
	rl.DrawTexture(game.textures.textureSheets.bg, 0, 0, rl.White)

	DrawHUD(game.textures, game.windowDimens, game.playerLives, game.playerPoints)

	drawSpawnedObjects(game)

	DrawPlayer(game.player)

	DrawPlatforms(game.textures, game.windowDimens, game.spriteSize)
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
