package main

import rl "github.com/gen2brain/raylib-go/raylib"

type TextureAtlas struct {
	sheet     rl.Texture2D
	bg        rl.Texture2D
	platforms []rl.Rectangle
	overlays  []rl.Rectangle
	scenery   []rl.Rectangle
	food      Grid[rl.Rectangle]
	objects   []rl.Rectangle
}

func NewTextureAtlas(fileName string, bgFileName string, windowSize [2]float32) TextureAtlas {
	objects := rl.LoadImage(fileName)

	bgImage := rl.LoadImage(bgFileName)

	rl.ImageResizeNN(bgImage, int32(windowSize[0]), int32(windowSize[1]))

	atlas := TextureAtlas{
		sheet:     rl.LoadTextureFromImage(objects),
		bg:        rl.LoadTextureFromImage(bgImage),
		platforms: makeRectangles(7, 15, 0, 15, 16),
		overlays:  makeRectangles(2, 17, 105, 17, 16),
		scenery:   makeRectangles(5, 16, 139, 16, 16),
		food: Grid[rl.Rectangle]{
			objects: append(append(makeRectangles(15, 16, 16, 16, 16), makeRectangles(15, 16, 32, 16, 16)...), makeRectangles(15, 16, 48, 16, 16)...),
			columns: 15,
			rows:    3},
		objects: makeRectangles(15, 16, 64, 16, 16),
	}

	rl.UnloadImage(objects)
	rl.UnloadImage(bgImage)

	return atlas
}

func makeRectangles(count int, stride float32, startPosition float32, width float32, height float32) []rl.Rectangle {
	bs := make([]rl.Rectangle, count)

	for i := 0; i < count; i++ {
		bs[i] = rl.NewRectangle(startPosition+(stride*float32(i)), 0, width, height)
	}

	return bs
}
