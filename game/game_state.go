package game

import (
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type State uint32

const (
	TITLE            State = 0
	PAUSED           State = 1
	GAME             State = 2
	GAME_OVER        State = 3
	HIGH_SCORE_INPUT State = 4
	HIGH_SCORES      State = 5

	TOTAL_PLATFORMS             = 1280 / 64
	GAME_OVER_FONT_SIZE float32 = 100.
	GAME_FONT_SIZE      float32 = GAME_OVER_FONT_SIZE / 2.
	ONE_UP              uint32  = 1000
	PLAYER_HIT_MAX      uint32  = 10
	MAX_POINTS          uint32  = 999999999
)

type GameState struct {
	collectables               Spawner
	exitGame                   bool
	font                       rl.Font
	gameOverTextAnimationTimer float32
	gameOverTextPos            rl.Vector2
	highScoreNameBannerPos     rl.Vector2
	highScoresFilePath         string
	lastState                  State
	levelSpawnBoundaries       Boundaries
	objects                    Spawner
	pauseScreenButtons         []Button
	platformHitBoxes           [TOTAL_PLATFORMS]HitBox
	player                     Player
	playerHitCounter           uint32
	playerLives                uint32
	playerOneUpCounter         uint32
	playerPoints               uint32
	pointerAnimation           LinearFrameAnimation
	scoreboard                 ScoreBoard
	spriteSize                 float32
	state                      State
	textInput                  TextInput
	textures                   TextureAtlas
	titleScreenButtons         ButtonGroup
	titleScreenCollectables    Spawner
	titleScreenImagePos        rl.Vector2
	windowDimens               WindowDimens
}

func ExitGame(game *GameState) bool {
	return game.exitGame
}

func NewGameState(windowDimens [2]float32) GameState {
	textures := NewTextureAtlas(windowDimens)

	const spriteSize float32 = 64.

	levelSpawnBoundaries := Boundaries{
		bottom: windowDimens[1]/2. + spriteSize*2.0,
		top:    0.0 - textures.food[0].Width,
		width:  windowDimens[0]}

	gameFont := rl.LoadFont(MkAssetDir("main.ttf"))

	window := WindowDimens{width: windowDimens[0], height: windowDimens[1]}

	cX, cY := Center(window)

	pathSep := string(os.PathSeparator)

	highScoresFilePath := "." + pathSep + "assets" + pathSep + "high_scores.json"

	scores, err := DecodeInto(highScoresFilePath, ScoreBoard{})

	if err != nil {
		panic("Unable to load high scores!! Exiting")
	}

	return GameState{
		collectables: NewSpawner(rl.GetFrameTime()*20, 200., len(textures.food), len(textures.food)/3, levelSpawnBoundaries),

		exitGame: false,

		font: gameFont,

		gameOverTextAnimationTimer: 0.,

		gameOverTextPos: rl.NewVector2(cX-(3*GAME_OVER_FONT_SIZE), 0-GAME_OVER_FONT_SIZE),

		highScoreNameBannerPos: rl.NewVector2(cX-(GAME_FONT_SIZE*5.5), 150.),

		highScoresFilePath: highScoresFilePath,

		lastState: TITLE,

		objects: NewSpawner(rl.GetFrameTime()*5, 300., len(textures.objects), len(textures.objects), levelSpawnBoundaries),

		pauseScreenButtons: []Button{
			NewButton(func(state *GameState) {
				state.lastState = PAUSED
				state.state = GAME
			},
				rl.NewVector2(cX-(GAME_FONT_SIZE*1.5), cY-GAME_FONT_SIZE*2.5)),
			NewButton(func(state *GameState) {
				state.lastState = PAUSED
				state.state = TITLE
			}, rl.NewVector2(cX-(GAME_FONT_SIZE*1.5), cY-GAME_FONT_SIZE*0.15))},

		platformHitBoxes:   makePlatforms(cY+spriteSize*2., spriteSize),
		player:             NewPlayer(MkAssetDir("player.png"), rl.NewVector2(100, cY+spriteSize+20.), cY+spriteSize+32., spriteSize+32.),
		playerHitCounter:   0,
		playerLives:        1,
		playerOneUpCounter: 0,
		playerPoints:       0,
		pointerAnimation:   NewAnimation(1.0, true, 10),
		scoreboard:         scores,
		spriteSize:         spriteSize,
		state:              TITLE,
		textInput:          NewTextInput(GAME_FONT_SIZE, 5, rl.NewVector2(cX-(GAME_FONT_SIZE*3.5), 250.), rl.RayWhite),
		textures:           textures,
		titleScreenButtons: ButtonGroup{
			lastActive: 0,
			buttons: []Button{
				NewButton(func(state *GameState) {
					state.state = GAME

					if state.lastState == PAUSED {
						ResetGameState(state)
					}
				}, rl.NewVector2(cX-(GAME_FONT_SIZE*1.5), 275.)),
				NewButton(func(state *GameState) { state.state = HIGH_SCORES }, rl.NewVector2(cX-(GAME_FONT_SIZE*1.5), 400.)),
				NewButton(func(state *GameState) { state.exitGame = true }, rl.NewVector2(cX-(GAME_FONT_SIZE*1.5), 525.))}},
		titleScreenCollectables: NewSpawner(rl.GetFrameTime()*20, 200., len(textures.food), len(textures.food)/3, Boundaries{
			bottom: windowDimens[1] + spriteSize*2.0,
			top:    0.0 - textures.food[0].Width,
			width:  windowDimens[0]},
		),
		titleScreenImagePos: rl.NewVector2(cX-float32(textures.textureSheets.titleScreenImage.Width/2), 100.),
		windowDimens:        window,
	}
}

func (state *GameState) Release() {
	rl.UnloadFont(state.font)
	state.textures.Release()
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

func ResetGameState(game *GameState) {
	levelSpawnBoundaries := Boundaries{
		bottom: game.windowDimens.height/2. + game.spriteSize*2.0,
		top:    0.0 - game.textures.food[0].Width,
		width:  game.windowDimens.width}

	game.collectables = NewSpawner(rl.GetFrameTime()*20, 200., len(game.textures.food), len(game.textures.food)/3, levelSpawnBoundaries)
	game.objects = NewSpawner(rl.GetFrameTime()*5, 300., len(game.textures.objects), len(game.textures.objects), levelSpawnBoundaries)
	game.gameOverTextAnimationTimer = 0.

	totalPauseButtons := len(game.pauseScreenButtons)
	for i := 0; i < totalPauseButtons; i++ {
		game.pauseScreenButtons[i].state = NORMAL
	}

	for i := 0; i < TOTAL_PLATFORMS; i++ {
		game.platformHitBoxes[i].damageCounter = NewDamageCounter(25.)
	}

	game.playerHitCounter = 0

	_, cY := Center(game.windowDimens)
	game.player = resettedPlayer(game.player, rl.NewVector2(100, cY+game.spriteSize+20.), cY+game.spriteSize+32.)

	game.playerLives = 1
	game.playerPoints = 0
	game.playerOneUpCounter = 0
}

func resettedPlayer(player Player, startPosition rl.Vector2, ground float32) Player {
	jumpHeight := float32(300.0)

	return Player{
		animation: LinearFrameAnimation{
			timer:  NewTimer(1.0, true),
			frames: 5,
		},
		fell:         false,
		halt:         false,
		stunTimer:    stunTimer(),
		isMoving:     false,
		isJumping:    false,
		originalSize: 32.,
		position:     startPosition,
		texture:      player.texture,
		textureBox:   rl.NewRectangle(0, 0, 32., 32.),
		textureSize:  player.textureSize,
		physics: Physics[rl.Vector2]{
			bottom:     ground,
			gravity:    -500,
			ground:     startPosition.Y,
			jumpHeight: -jumpHeight,
			velocity:   rl.NewVector2(200., 0.)},
		wasHit: false,
	}
}
