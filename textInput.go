package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextInput struct {
	value     string
	isEditing bool
	maxSize   int
	label     string
}

func NewTextInput(value string, maxSize int, label string) TextInput {
	return TextInput{
		value:   value,
		maxSize: maxSize,
		label:   label,
	}
}

func (ti *TextInput) Value() string {
	return ti.value
}

func (ti *TextInput) SetValue(value string) {
	ti.value = value
}

func (ti *TextInput) Draw(bounds rl.Rectangle, fontSize float32) {

	gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, int64(fontSize))
	gui.Label(rl.Rectangle{
		X:      bounds.X,
		Y:      bounds.Y - fontSize,
		Width:  bounds.Width,
		Height: fontSize,
	}, ti.label)

	inputToggle := gui.TextBox(rl.Rectangle{
		X:      bounds.X,
		Y:      bounds.Y,
		Width:  bounds.Width,
		Height: bounds.Height,
	}, &ti.value, ti.maxSize, ti.isEditing)

	if inputToggle {
		ti.isEditing = !ti.isEditing
	}
}
