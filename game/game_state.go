package game

import (
	"os"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
	textures := NewTextureAtlas(windowDimens)

	const spriteSize float32 = 64.

	spawnBoundaries := Boundaries{
		bottom: windowDimens[1]/2. + spriteSize*2.0,
		top:    0.0 - textures.food[0].Width,
		width:  windowDimens[0]}

	gameFont := rl.LoadFont(mkAssetDir("main.ttf"))

	return GameState{
		collectables:               NewSpawner(rl.GetFrameTime()*20, 200., len(textures.food), len(textures.food)/3, spawnBoundaries),
		font:                       gameFont,
		gameOverTextAnimationTimer: 0.,
		gameOverTextPos:            rl.NewVector2((windowDimens[0]/2.)-(3*GAME_OVER_FONT_SIZE), 0-GAME_OVER_FONT_SIZE),
		highScores:                 []Score{},
		objects:                    NewSpawner(rl.GetFrameTime()*5, 300., len(textures.objects), len(textures.objects), spawnBoundaries),
		platformHitBoxes:           makePlatforms(windowDimens[1]/2.+spriteSize*2., spriteSize),
		player:                     NewPlayer(mkAssetDir("player.png"), rl.NewVector2(100, (windowDimens[1]/2.)+spriteSize+20.), windowDimens[1]+spriteSize+32., spriteSize+32.),
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
