package main

import rl "github.com/gen2brain/raylib-go/raylib"

func main() {
	windowDimens := [2]float32{1280, 720}
	spriteSize := float32(64.)

	rl.InitWindow(int32(windowDimens[0]), int32(windowDimens[1]), "Catcher")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	player := NewPlayer("./assets/player.png", rl.NewVector2(50., (windowDimens[1]/2.)+20.), spriteSize+32.)

	textures := NewTextureAtlas("./assets/level.png", "./assets/bg.png", windowDimens)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		player = UpdatePlayer(player, rl.GetFrameTime())

		rl.DrawTexture(textures.bg, 0, 0, rl.White)

		DrawPlatforms(&textures, windowDimens, spriteSize)

		DrawPlayer(player)

		rl.EndDrawing()
	}
}
