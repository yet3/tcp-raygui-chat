package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	HOME_SCENE SceneName = "HOME"
	SETUP_SCENE SceneName= "SETUP"
	CHAT_SCENE SceneName= "CHAT"
)

const (
	SCREEN_WIDTH  = 550
	SCREEN_HEIGHT = 720
)

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "tcp raygui chat")
	defer rl.CloseWindow()

	sceneManager := NewSceneManager()
  sceneManager.SetScene(NewHome())

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		sceneManager.Update()

		rl.EndDrawing()
	}
}
