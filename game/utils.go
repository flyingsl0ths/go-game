package game

import (
	"math"
	"math/rand"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type WindowDimens struct {
	width  float32
	height float32
}

func Center(dimensions WindowDimens) (x float32, y float32) {
	x = dimensions.width / 2.
	y = dimensions.height / 2.
	return x, y
}

func RandRange(min int, max int) int {
	if min > max {
		return 0
	} else {
		return rand.Intn(max-min) + min
	}
}

func CellFrom(x float32, cellSize float32, maxCells int) int {
	cell := int(math.Floor(float64(x / cellSize)))

	if cell >= maxCells {
		cell = maxCells - 1
	} else if cell < 0 {
		cell = 0
	}

	return int(cell)
}

func CollideWithSides(r1 rl.Rectangle, r2 rl.Rectangle) string {
	dx := (r1.X + r1.Width/2) - (r2.X + r2.Width/2)
	dy := (r1.Y + r1.Height/2) - (r2.Y + r2.Height/2)
	width := (r1.Width + r2.Width) / 2
	height := (r1.Height + r2.Height) / 2
	crossWidth := width * dy
	crossHeight := height * dx
	collision := "none"

	if math.Abs(float64(dx)) <= float64(width) && math.Abs(float64(dy)) <= float64(height) {
		if crossWidth > crossHeight {
			if crossWidth > -crossHeight {
				collision = "bottom"
			} else {
				collision = "left"
			}
		} else {

			if crossWidth > -crossHeight {
				collision = "right"
			} else {
				collision = "top"
			}
		}
	}

	return collision
}

func MkDir(rootDirectory string) func(string) string {
	return func(file string) string { return (rootDirectory + string(os.PathSeparator) + file) }
}

var assetDir string = "." + string(os.PathSeparator) + "assets"
var MkAssetDir func(string) string = MkDir(assetDir)
