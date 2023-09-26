package main

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
	food          Grid[rl.Rectangle]
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
		platforms: makeRectangles(7, 15, 0, 15, 16),
		overlays:  makeRectangles(2, 17, 105, 17, 16),
		scenery:   makeRectangles(5, 16, 139, 16, 16),
		food: Grid[rl.Rectangle]{
			objects: append(append(makeRectangles(15, 16, 16, 16, 16), makeRectangles(15, 16, 32, 16, 16)...), makeRectangles(15, 16, 48, 16, 16)...),
			columns: 15,
			rows:    3},
		objects: makeRectangles(15, 16, 64, 16, 16),
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

func makeRectangles(count int, stride float32, startPosition float32, width float32, height float32) []rl.Rectangle {
	bs := make([]rl.Rectangle, count)

	for i := 0; i < count; i++ {
		bs[i] = rl.NewRectangle(startPosition+(stride*float32(i)), 0, width, height)
	}

	return bs
}
