package game

import rl "github.com/gen2brain/raylib-go/raylib"

type ButtonState uint8

const (
	NORMAL  = ButtonState(0)
	HOVER   = ButtonState(1)
	PRESSED = ButtonState(2)
)

type Button struct {
	onClick      func(state *GameState)
	position     rl.Rectangle
	state        ButtonState
	textPosition rl.Vector2
}

type ButtonGroup struct {
	lastActive uint
	buttons    []Button
}

func NewButton(onClick func(state *GameState), position rl.Vector2) Button {
	const BUTTON_WIDTH float32 = 281.
	const BUTTON_HEIGHT float32 = 103.

	return Button{
		onClick:      onClick,
		position:     rl.NewRectangle(position.X-GAME_FONT_SIZE, position.Y-GAME_FONT_SIZE/2, BUTTON_WIDTH, BUTTON_HEIGHT),
		state:        NORMAL,
		textPosition: position,
	}
}
