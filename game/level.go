package game

import rl "github.com/gen2brain/raylib-go/raylib"

func DrawHUD(textureAtlas *TextureAtlas, windowDimens WindowDimens, lives uint32, points uint32) {
	drawPlayerIcons(textureAtlas, lives)
	drawCollectableIcons(textureAtlas, windowDimens, points)
}

func drawPlayerIcons(textureAtlas *TextureAtlas, lives uint32) {
	const PADDING float32 = 5.
	lifeIconPos := rl.NewVector2(PADDING, PADDING)

	lifeIcon := textureAtlas.hud[0]
	lifeIconSize := lifeIcon.Width * 2

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		lifeIcon,
		rl.NewRectangle(0., 0., lifeIconSize+2., lifeIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	lifeIconPos.X += lifeIconSize

	multiplierIconSize := textureAtlas.hud[2].Width * 2.

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[2],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y+(multiplierIconSize/2.), multiplierIconSize, multiplierIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	lifeIconPos.X += multiplierIconSize + PADDING
	numberIconSize := textureAtlas.hud[3].Width * 2
	lifeIconPos.Y += numberIconSize/2. + PADDING/2.

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[3],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y, numberIconSize, numberIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	lifeIconPos.X += numberIconSize + PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[3],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y, numberIconSize, numberIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)
}

func drawCollectableIcons(textureAtlas *TextureAtlas, windowDimens WindowDimens, points uint32) {
	const PADDING float32 = 5.

	pointsIcon := textureAtlas.hud[1]
	multiplierIcon := textureAtlas.hud[2]

	windowWidth := windowDimens.width
	multiplierIconSize := multiplierIcon.Width * 2
	pointsIconSize := pointsIcon.Width * 2

	offset := pointsIconSize + multiplierIconSize + (textureAtlas.hud[3].Width * 2 * 10)

	pointsIconPos := rl.NewVector2(windowWidth-offset, PADDING)

	yOffset := pointsIconPos.Y + PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		pointsIcon,
		rl.NewRectangle(pointsIconPos.X, pointsIconPos.Y, pointsIcon.Width*2., pointsIcon.Height*2),
		rl.NewVector2(0., 0.), 0., rl.White)

	pointsIconPos.X += pointsIconSize + PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		multiplierIcon,
		rl.NewRectangle(pointsIconPos.X, yOffset, multiplierIconSize, multiplierIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	pointsIconPos.X += multiplierIconSize

	const MAX_POINTS = 9

	for i := 0; i < MAX_POINTS; i++ {
		numberIcon := textureAtlas.hud[3]
		numberIconSize := numberIcon.Width * 2

		rl.DrawTexturePro(textureAtlas.textureSheets.hud,
			numberIcon,
			rl.NewRectangle(pointsIconPos.X+(numberIconSize*float32(i))+PADDING,
				yOffset, numberIconSize, numberIconSize),
			rl.NewVector2(0., 0.), 0., rl.White)
	}
}

func DrawPlatforms(textureAtlas *TextureAtlas, windowDimens WindowDimens, textureSize float32) {
	startPosition := [2]float32{windowDimens.width / 2., windowDimens.height/2. + textureSize*2}

	rl.DrawTexturePro(textureAtlas.textureSheets.level,
		textureAtlas.platforms[0],
		rl.NewRectangle(0., startPosition[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	rl.DrawTexturePro(textureAtlas.textureSheets.level,
		textureAtlas.platforms[1],
		rl.NewRectangle(textureSize, startPosition[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	rl.DrawTexturePro(textureAtlas.textureSheets.level,
		textureAtlas.platforms[2],
		rl.NewRectangle(textureSize*2, startPosition[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	const currentlyDrawnAtEdge = 3
	toDraw := int(windowDimens.width/textureSize) - currentlyDrawnAtEdge

	for i := 3; i < toDraw; i++ {
		rl.DrawTexturePro(textureAtlas.textureSheets.level,
			textureAtlas.platforms[3],
			rl.NewRectangle(textureSize*float32(i), startPosition[1], textureSize, textureSize),
			rl.NewVector2(0., 0.), 0., rl.White)
	}

	total := len(textureAtlas.platforms)

	rl.DrawTexturePro(textureAtlas.textureSheets.level,
		textureAtlas.platforms[total-3],
		rl.NewRectangle(windowDimens.width-(textureSize*3), startPosition[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	rl.DrawTexturePro(textureAtlas.textureSheets.level,
		textureAtlas.platforms[total-2],
		rl.NewRectangle(windowDimens.width-(textureSize*2), startPosition[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	rl.DrawTexturePro(textureAtlas.textureSheets.level,
		textureAtlas.platforms[total-1],
		rl.NewRectangle(windowDimens.width-textureSize, startPosition[1], textureSize, textureSize),
		rl.NewVector2(0., 0.), 0., rl.White)
}
