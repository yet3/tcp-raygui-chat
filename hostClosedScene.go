package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type HostClosedScene struct{}

func NewHostClosedScene() *HostClosedScene {
	return &HostClosedScene{}
}

func (hcs *HostClosedScene) UpdateScene(sceneManager *SceneManager) {
	const (
		INFO       = "Host has closed the connection"
		TEXT_SIZE  = 28
		BTN_WIDTH  = 160
		BTN_HEIGHT = 50
	)
	rl.ClearBackground(rl.White)

	gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, TEXT_SIZE)
	textWidth := float32(rl.MeasureText(INFO, TEXT_SIZE))
	gui.Label(rl.Rectangle{
		X:      SCREEN_WIDTH/2 - textWidth/2,
		Y:      SCREEN_HEIGHT/2 - TEXT_SIZE,
		Width:  textWidth,
		Height: TEXT_SIZE,
	}, INFO)

	gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, 24)
	homePressed := gui.Button(rl.Rectangle{
		X:      SCREEN_WIDTH/2 - BTN_WIDTH/2,
		Y:      SCREEN_HEIGHT/2 + 10,
		Width:  BTN_WIDTH,
		Height: BTN_HEIGHT,
	}, "Home")

	if homePressed {
		sceneManager.SetScene(NewHome())
	}
}
