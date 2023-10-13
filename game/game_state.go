	spawnBoundaries := Boundaries{
		bottom: windowDimens[1]/2. + spriteSize*2.0,
		top:    0.0 - textures.food[0].Width,
		width:  windowDimens[0]}

	mover := func(item *Item, delta float32) {
		item.position.Y += item.gravity * delta
		item.rotation += float32(EaseOutCirc(0.2))
	}
		collectables: NewSpawner(1.25, 150., len(textures.food), len(textures.food), spawnBoundaries, mover),
		objects:      NewSpawner(rl.GetFrameTime()*10, 300., len(textures.objects), len(textures.objects), spawnBoundaries, mover),
	game.player = UpdatePlayer(game.player, delta)

	UpdateSpawner(&game.objects, delta)
	drawSpawnedObjects(game)
func drawSpawnedObjects(game *GameState) {
	for _, item := range game.objects.items {
		rl.DrawTexturePro(game.textures.textureSheets.level,
			game.textures.objects[item.itemID],
			rl.NewRectangle(item.position.X, item.position.Y, game.spriteSize, game.spriteSize),
			rl.NewVector2(game.spriteSize/2., game.spriteSize/2.), item.rotation, rl.White)
	}
