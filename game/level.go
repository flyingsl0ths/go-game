package game

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawHUD(textureAtlas TextureAtlas, windowDimens WindowDimens, lives uint32, points uint32) {
	drawPlayerIcons(textureAtlas, lives)
	drawCollectableIcons(textureAtlas, windowDimens, points)
}

func drawPlayerIcons(textureAtlas TextureAtlas, lives uint32) {
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

	playerLives := strconv.FormatUint(uint64(lives), 10)

	firstNumOffset := 0
	secondNumOffset := 0

	if lives >= 10 {
		firstNumOffset = int(playerLives[0] - '0')
		secondNumOffset = int(playerLives[1] - '0')
	} else {
		secondNumOffset = int(playerLives[0] - '0')
	}

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[3+firstNumOffset],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y, numberIconSize, numberIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	lifeIconPos.X += numberIconSize + PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		textureAtlas.hud[3+secondNumOffset],
		rl.NewRectangle(lifeIconPos.X, lifeIconPos.Y, numberIconSize, numberIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)
}

func drawCollectableIcons(textureAtlas TextureAtlas, windowDimens WindowDimens, points uint32) {
	const IMAGE_PADDING float32 = 6.

	pointsIcon := textureAtlas.hud[1]
	multiplierIcon := textureAtlas.hud[2]

	windowWidth := windowDimens.width
	multiplierIconSize := multiplierIcon.Width * 2
	pointsIconSize := pointsIcon.Width * 2

	offset := pointsIconSize + multiplierIconSize + (textureAtlas.hud[3].Width * 2 * 10)

	pointsIconPos := rl.NewVector2(windowWidth-offset, IMAGE_PADDING)

	yOffset := pointsIconPos.Y + IMAGE_PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		pointsIcon,
		rl.NewRectangle(pointsIconPos.X, pointsIconPos.Y, pointsIcon.Width*2., pointsIcon.Height*2),
		rl.NewVector2(0., 0.), 0., rl.White)

	pointsIconPos.X += pointsIconSize + IMAGE_PADDING

	rl.DrawTexturePro(textureAtlas.textureSheets.hud,
		multiplierIcon,
		rl.NewRectangle(pointsIconPos.X, yOffset, multiplierIconSize, multiplierIconSize),
		rl.NewVector2(0., 0.), 0., rl.White)

	pointsIconPos.X += multiplierIconSize

	const MAX_POINTS = 9
	const ZERO_IMAGE_INDEX = 3

	playerPoints := strconv.FormatUint(uint64(points), 10)

	numberIconSize := textureAtlas.hud[ZERO_IMAGE_INDEX].Width * 2
	padding := MAX_POINTS - len(playerPoints)
	playerPointsStride := 0

	for i := 0; i < padding; i++ {
		numberIcon := textureAtlas.hud[ZERO_IMAGE_INDEX]

		rl.DrawTexturePro(textureAtlas.textureSheets.hud,
			numberIcon,
			rl.NewRectangle(pointsIconPos.X+(numberIconSize*float32(i))+IMAGE_PADDING,
				yOffset, numberIconSize, numberIconSize),
			rl.NewVector2(0., 0.), 0., rl.White)

		playerPointsStride = i
	}

	if padding > 0 {
		playerPointsStride += 1
	}

	for _, c := range playerPoints {
		numberIcon := textureAtlas.hud[ZERO_IMAGE_INDEX+int(c-'0')]

		rl.DrawTexturePro(textureAtlas.textureSheets.hud,
			numberIcon,
			rl.NewRectangle(pointsIconPos.X+(numberIconSize*float32(playerPointsStride))+IMAGE_PADDING,
				yOffset, numberIconSize, numberIconSize),
			rl.NewVector2(0., 0.), 0., rl.White)

		playerPointsStride += 1
	}
}

func DrawPlatforms(textureAtlas TextureAtlas, windowDimens WindowDimens, textureSize float32, hitboxes [26]HitBox) {
	const AT_EDGE = 2
	const REGULAR_TILE = 3

	startPosition := windowDimens.height/2. + textureSize*2
	toDraw := int(windowDimens.width / textureSize)
	remaining := toDraw - (AT_EDGE + 1)
	total := len(textureAtlas.platforms)

	for i := 0; i < toDraw; i++ {
		tile := REGULAR_TILE

		if i <= AT_EDGE {
			tile = i
		} else if i >= remaining {
			tile = total - (toDraw % i)
			println(tile)
		}

		rl.DrawTexturePro(textureAtlas.textureSheets.level,
			textureAtlas.platforms[tile],
			rl.NewRectangle(textureSize*float32(i), startPosition, textureSize, textureSize),
			rl.NewVector2(0., 0.), 0., rl.White)
	}
}
