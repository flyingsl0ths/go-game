package game

import rl "github.com/gen2brain/raylib-go/raylib"

type TextureSheets struct {
	level rl.Texture2D
	bg    rl.Texture2D
	hud   rl.Texture2D
}

type TextureAtlas struct {
	textureSheets TextureSheets
	hud           []rl.Rectangle
	platforms     []rl.Rectangle
	overlays      []rl.Rectangle
	scenery       []rl.Rectangle
	food          []rl.Rectangle
	objects       []rl.Rectangle
}

func NewTextureAtlas(fileName string, bgFileName string, hudFileName string, windowSize [2]float32) TextureAtlas {
	objects := rl.LoadImage(fileName)
	bgImage := rl.LoadImage(bgFileName)
	hud := rl.LoadImage(hudFileName)

	defer rl.UnloadImage(hud)
	defer rl.UnloadImage(objects)
	defer rl.UnloadImage(bgImage)

	rl.ImageResizeNN(bgImage, int32(windowSize[0]), int32(windowSize[1]))

	atlas := TextureAtlas{
		textureSheets: TextureSheets{
			level: rl.LoadTextureFromImage(objects),
			bg:    rl.LoadTextureFromImage(bgImage),
			hud:   rl.LoadTextureFromImage(hud),
		},
		hud:       makeHUDRectangles(),
		platforms: makeRectangles(7, 15, rl.NewVector2(0., 0.), 15, 16),
		overlays:  makeRectangles(2, 17, rl.NewVector2(105., 0.), 17, 16),
		scenery:   makeRectangles(5, 16, rl.NewVector2(139., 0.), 16, 16),
		food:      makeCollectables(),
		objects:   makeRectangles(15, 16, rl.NewVector2(0., 64.), 16, 15),
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

	const NUMS = 12
	const NUMS_STRIDE = 8
	const NUMS_WIDTH = 7.
	const NUMS_COL = 5

	j := 3
	for i := 0; i < NUMS; i++ {
		var y float32 = 11.
		if i == NUMS_COL {
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
