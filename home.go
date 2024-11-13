package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Home struct {
}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) UpdateScene(sceneManager *SceneManager) {
	const (
		BTN_WIDTH  float32 = 180
		BTN_HEIGHT float32 = 80
		GAP        float32 = 20
	)

	rl.ClearBackground(rl.White)

	gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, 22)

	serverPressed := gui.Button(rl.Rectangle{
		X:      SCREEN_WIDTH/2 - BTN_WIDTH - GAP/2,
		Y:      SCREEN_HEIGHT/2 - BTN_HEIGHT/2,
		Width:  BTN_WIDTH,
		Height: BTN_HEIGHT,
	}, "As server")

	if serverPressed {
		sceneManager.SetScene(NewSetupScene(true))
	}

	clientPressed := gui.Button(rl.Rectangle{
		X:      SCREEN_WIDTH/2 + GAP/2,
		Y:      SCREEN_HEIGHT/2 - BTN_HEIGHT/2,
		Width:  BTN_WIDTH,
		Height: BTN_HEIGHT,
	}, "As client")

	if clientPressed {
		sceneManager.SetScene(NewSetupScene(false))
	}
}
