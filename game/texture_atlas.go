package game

import rl "github.com/gen2brain/raylib-go/raylib"

type TextureSheets struct {
	bg               rl.Texture2D
	buttons          rl.Texture2D
	hud              rl.Texture2D
	level            rl.Texture2D
	titleScreenImage rl.Texture2D
}

type TextureAtlas struct {
	buttons       []rl.Rectangle
	food          []rl.Rectangle
	hud           []rl.Rectangle
	objects       []rl.Rectangle
	overlays      []rl.Rectangle
	platforms     []rl.Rectangle
	scenery       []rl.Rectangle
	textureSheets TextureSheets
}

func (ta *TextureAtlas) Release() {
	textureSheets := &ta.textureSheets
	rl.UnloadTexture(textureSheets.bg)
	rl.UnloadTexture(textureSheets.buttons)
	rl.UnloadTexture(textureSheets.hud)
	rl.UnloadTexture(textureSheets.level)
	rl.UnloadTexture(textureSheets.titleScreenImage)
}

func NewTextureAtlas(windowSize [2]float32) TextureAtlas {
	objects := rl.LoadImage(MkAssetDir("level.png"))
	bg := rl.LoadImage(MkAssetDir("bg.png"))
	hud := rl.LoadImage(MkAssetDir("hud.png"))
	buttons := rl.LoadImage(MkAssetDir("button.png"))
	titleScreenImage := rl.LoadImage(MkAssetDir("title.png"))

	defer rl.UnloadImage(hud)
	defer rl.UnloadImage(objects)
	defer rl.UnloadImage(bg)
	defer rl.UnloadImage(buttons)
	defer rl.UnloadImage(titleScreenImage)

	rl.ImageResizeNN(bg, int32(windowSize[0]), int32(windowSize[1]))

	atlas := TextureAtlas{
		textureSheets: TextureSheets{
			level:            rl.LoadTextureFromImage(objects),
			bg:               rl.LoadTextureFromImage(bg),
			hud:              rl.LoadTextureFromImage(hud),
			buttons:          rl.LoadTextureFromImage(buttons),
			titleScreenImage: rl.LoadTextureFromImage(titleScreenImage),
		},
		hud:       makeHUDRectangles(),
		platforms: makeRectangles(7, 15, rl.NewVector2(0., 0.), 15, 16),
		overlays:  makeRectangles(2, 17, rl.NewVector2(105., 0.), 17, 16),
		scenery:   makeRectangles(5, 16, rl.NewVector2(139., 0.), 16, 16),
		food:      makeCollectables(),
		objects:   makeRectangles(15, 16, rl.NewVector2(0., 64.), 16, 15),
		buttons:   makeButtonRectangles(),
	}

	return atlas
}

func makeHUDRectangles() []rl.Rectangle {
	rs := make([]rl.Rectangle, 15)

	const LIFE_ICON = 0
	const COLLECTABLE_ICON = 1
	const MULTIPLIER_ICON = 2

	rs[LIFE_ICON] = rl.NewRectangle(0, 0, 16., 11.)
	rs[COLLECTABLE_ICON] = rl.NewRectangle(16., 0, 25., 11.)
	rs[MULTIPLIER_ICON] = rl.NewRectangle(41., 3., 8., 8.)

	const REMAINING = 10
	const HALFWAY = REMAINING / 2
	const NUMS_STRIDE = 8
	const NUMS_WIDTH = 7.
	const NUMS_COL = 5

	j := 3
	for i := 0; i <= REMAINING; i++ {
		var y float32 = 11.

		if i > HALFWAY {
			y += NUMS_STRIDE
		}

		rs[j] = rl.NewRectangle(float32(i)*NUMS_STRIDE, y, NUMS_WIDTH, NUMS_WIDTH)

		j += 1
	}

	return rs
}

func makeCollectables() []rl.Rectangle {
	collectables := makeRectangles(15, 16, rl.NewVector2(0., 16.), 16, 15)
	collectables = append(collectables, makeRectangles(15, 16, rl.NewVector2(0., 32.), 16, 15)...)
	collectables = append(collectables, makeRectangles(15, 16, rl.NewVector2(0., 48.), 16, 15)...)
	return collectables
}

func makeRectangles(count int, stride float32, startPosition rl.Vector2, width float32, height float32) []rl.Rectangle {
	rects := make([]rl.Rectangle, count)

	for i := 0; i < count; i++ {
		rects[i] = rl.NewRectangle(startPosition.X+(stride*float32(i)), startPosition.Y, width, height)
	}

	return rects
}

func makeButtonRectangles() []rl.Rectangle {
	res := make([]rl.Rectangle, 3)

	const NORMAL = 0
	const HOVER = 1
	const PRESSED = 2

	res[NORMAL] = rl.NewRectangle(0, 0, 380, 203)
	res[HOVER] = rl.NewRectangle(0, 203, 380, 203)
	res[PRESSED] = rl.NewRectangle(0, 203, 380, 197)

	return res
}
