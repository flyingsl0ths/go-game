package game

import rl "github.com/gen2brain/raylib-go/raylib"

type TextInput struct {
	blink          bool
	blinkTimer     float32
	blinkTimerStep float32
	cursorIndex    uint32
	fontSize       float32
	limit          uint32
	start          rl.Vector2
	text           []rune
	textColor      rl.Color
}

func NewTextInput(fontSize float32, limit uint32, start rl.Vector2, textColor rl.Color) TextInput {
	return TextInput{
		blink:          false,
		blinkTimer:     1.0,
		blinkTimerStep: 0.015,
		cursorIndex:    0,
		fontSize:       fontSize,
		limit:          limit,
		start:          start,
		text:           []rune(""),
		textColor:      textColor,
	}
}

func UpdateTextInput(input TextInput) TextInput {
	input_ := input

	input_.blinkTimer += input_.blinkTimerStep
	if input_.blinkTimer > 1.0 {
		input_.blinkTimer = 0.0
		input_.blink = !input_.blink
	}

	keyPressed := rl.GetKeyPressed()

	if isAplhaNum(keyPressed) && uint32(len(input_.text)) < input_.limit {
		input_.cursorIndex += 1
		input_.text = append(input_.text, rune(keyPressed))
	} else if keyPressed == rl.KeyBackspace && len(input_.text) > 0 {
		input_.cursorIndex -= 1
		input_.text = input_.text[0 : len(input_.text)-1]
	}

	return input_
}

func isAplhaNum(key int32) bool {
	return key >= rl.KeyApostrophe && key <= rl.KeyZ
}

func DrawTextInput(input TextInput) {
	cursorPos := input.start

	if input.cursorIndex > 0 {
		cursorPos = rl.NewVector2(input.start.X+textWidth(input.text, input.fontSize)+5, input.start.Y)
	}

	rl.DrawText(string(input.text), int32(input.start.X), int32(input.start.Y), int32(input.fontSize), input.textColor)
	if input.blink {
		rl.DrawRectangle(int32(cursorPos.X), int32(cursorPos.Y), 5, int32(input.fontSize), input.textColor)
	}
}

func textWidth(text []rune, fontSize float32) float32 {
	return float32(rl.MeasureText(string(text), int32(fontSize)))
}
