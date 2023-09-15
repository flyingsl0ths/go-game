package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	windowDimens := [2]float32{1280, 720}
	rl.InitWindow(int32(windowDimens[0]), int32(windowDimens[1]), "Catcher")

	defer rl.CloseWindow()

	image := rl.LoadImage("./assets/player.png")

	rl.SetTargetFPS(60)

	player := NewPlayer(image, rl.NewVector2(50., 50.), float32(spriteSize*2))

	rl.UnloadImage(image)

	textures := NewTextureAtlas("./assets/level.png", "./assets/bg.png", windowDimens)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		player = UpdatePlayer(player)

		rl.DrawTexture(textures.bg, 0, 0, rl.White)

		DrawPlatforms(&textures, windowDimens)

		DrawPlayer(player)

		rl.EndDrawing()
	}
}
