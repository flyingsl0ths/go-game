package main

import rl "github.com/gen2brain/raylib-go/raylib"

func DrawPlatforms(textureAtlas *TextureAtlas, windowDimens [2]float32, textureSize float32) {
	windowCenter := [2]float32{windowDimens[0] / 2., windowDimens[1]/2. + textureSize}

	rl.DrawTexturePro(textureAtlas.sheet,
		textureAtlas.platforms[0],
		rl.NewRectangle(0., windowCenter[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	rl.DrawTexturePro(textureAtlas.sheet,
		textureAtlas.platforms[1],
		rl.NewRectangle(textureSize, windowCenter[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	rl.DrawTexturePro(textureAtlas.sheet,
		textureAtlas.platforms[2],
		rl.NewRectangle(textureSize*2, windowCenter[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	const currentlyDrawnAtEdge = 3
	toDraw := int(windowDimens[0]/textureSize) - currentlyDrawnAtEdge

	for i := 3; i < toDraw; i++ {
		rl.DrawTexturePro(textureAtlas.sheet,
			textureAtlas.platforms[3],
			rl.NewRectangle(textureSize*float32(i), windowCenter[1], textureSize, textureSize),
			rl.NewVector2(0., 0.), 0., rl.White)
	}

	total := len(textureAtlas.platforms)

	rl.DrawTexturePro(textureAtlas.sheet,
		textureAtlas.platforms[total-3],
		rl.NewRectangle(windowDimens[0]-(textureSize*3), windowCenter[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	rl.DrawTexturePro(textureAtlas.sheet,
		textureAtlas.platforms[total-2],
		rl.NewRectangle(windowDimens[0]-(textureSize*2), windowCenter[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	rl.DrawTexturePro(textureAtlas.sheet,
		textureAtlas.platforms[total-1],
		rl.NewRectangle(windowDimens[0]-textureSize, windowCenter[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)
}
