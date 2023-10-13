package game

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	TITLE       uint32 = 0
	PAUSED      uint32 = 1
	GAME        uint32 = 2
	HIGH_SCORES uint32 = 3
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
	playerLives  Lives
	player       Player
	playerPoints Points
	powerUps     Spawner
	spriteSize   float32
	state        State
	textures     TextureAtlas
	windowDimens WindowDimens
}

func NewGameState(windowDimens [2]float32, spriteSize float32) GameState {
	player := NewPlayer("./assets/player.png", rl.NewVector2(50., (windowDimens[1]/2.)+spriteSize+20.), spriteSize+32.)

	textures := NewTextureAtlas("./assets/level.png", "./assets/bg.png", "./assets/hud.png", windowDimens)

	spawnBoundaries := Boundaries{
		bottom: windowDimens[1]/2. + spriteSize*2.0,
		top:    0.0 - textures.food[0].Width,
		width:  windowDimens[0]}

	mover := func(item *Item, delta float32) {
		item.position.Y += item.gravity * delta
		item.rotation += float32(EaseOutCirc(0.2))
	}

	return GameState{
		collectables: NewSpawner(1.25, 150., len(textures.food), len(textures.food), spawnBoundaries, mover),
		highScores:   []Score{},
		objects:      NewSpawner(rl.GetFrameTime()*10, 300., len(textures.objects), len(textures.objects), spawnBoundaries, mover),
		playerLives:  [2]rune{'0', '1'},
		player:       player,
		playerPoints: [9]rune{'0', '0', '0', '0', '0', '0', '0', '0', '0'},
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

	UpdateSpawner(&game.objects, delta)
}

func drawGameState(game *GameState) {
	rl.DrawTexture(game.textures.textureSheets.bg, 0, 0, rl.White)

	textures := &game.textures

	DrawHUD(textures, game.windowDimens, game.playerLives, game.playerPoints)

	drawSpawnedObjects(game)

	DrawPlayer(game.player)

	DrawPlatforms(textures, game.windowDimens, game.spriteSize)
}

func drawSpawnedObjects(game *GameState) {
	for _, item := range game.objects.items {
		rl.DrawTexturePro(game.textures.textureSheets.level,
			game.textures.objects[item.itemID],
			rl.NewRectangle(item.position.X, item.position.Y, game.spriteSize, game.spriteSize),
			rl.NewVector2(game.spriteSize/2., game.spriteSize/2.), item.rotation, rl.White)
	}
}
