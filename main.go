package main

import (
	gm "github.com/flyingsl0ths/go-game/game"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	windowDimens := [2]float32{1280, 720}

	rl.InitWindow(int32(windowDimens[0]), int32(windowDimens[1]), "Catcher")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	game := gm.NewGameState(windowDimens)

	for !rl.WindowShouldClose() {
		gm.RunGame(&game, rl.GetFrameTime())
	}
}
