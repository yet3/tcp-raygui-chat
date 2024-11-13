package main

import (
	"math/rand"
	"strconv"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SetupScene struct {
	isHost   bool
	errorMsg MultiLineText

	hostInput     TextInput
	portInput     TextInput
	usernameInput TextInput

	quich chan bool
}

func NewSetupScene(isHost bool) *SetupScene {
	return &SetupScene{
		isHost:        isHost,
		errorMsg:      NewMultiLineText("", SCREEN_WIDTH-40, 18),
		hostInput:     NewTextInput("", 100, "Host"),
		portInput:     NewTextInput("", 100, "Port"),
		usernameInput: NewTextInput("", 100, "Username"),
		quich:         make(chan bool, 1),
	}
}

func drawInput(input *TextInput, yIdx float32) float32 {
	const (
		INPUT_WIDTH  float32 = 300
		INPUT_HEIGHT         = 45
		FONT_SIZE            = 20
		GAP                  = 20
		TOTAL_HEIGHT         = INPUT_HEIGHT + FONT_SIZE + GAP
	)

	y := yIdx*TOTAL_HEIGHT + SCREEN_HEIGHT/3.5

	input.Draw(rl.Rectangle{
		X:      SCREEN_WIDTH/2 - INPUT_WIDTH/2,
		Y:      y - INPUT_HEIGHT/2,
		Width:  INPUT_WIDTH,
		Height: INPUT_HEIGHT,
	}, FONT_SIZE)

	return y + INPUT_HEIGHT/2
}

func (ss *SetupScene) resetValues() {
	ss.errorMsg.Clear()
	ss.hostInput.SetValue("")
	ss.portInput.SetValue("")
	ss.usernameInput.SetValue("")
}

func (ss *SetupScene) defaultValues() {
	ss.resetValues()
	ss.hostInput.SetValue("localhost")
	ss.portInput.SetValue("3000")

	if !ss.isHost {
		ss.usernameInput.SetValue("Guest" + strconv.Itoa(rand.Intn(1000)))
	}
}

func (ss *SetupScene) checkInputsForErrors() bool {
	if len(ss.hostInput.Value()) == 0 {
		ss.errorMsg.SetText("Field 'Host' cannot be empty")
		return true
	}

	if len(ss.portInput.Value()) == 0 {
		ss.errorMsg.SetText("Field 'Port' cannot be empty")
		return true
	}

	if len(ss.usernameInput.value) == 0 && !ss.isHost {
		ss.errorMsg.SetText("Field 'Username' cannot be empty")
		return true
	}

	return false
}

func (ss *SetupScene) UpdateScene(sceneManager *SceneManager) {
	const (
		BTN_WIDTH  float32 = 120
		BTN_HEIGHT float32 = 50
	)

	rl.ClearBackground(rl.White)

	gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, 19)
	goBackPressed := gui.Button(rl.Rectangle{
		X:      16,
		Y:      16,
		Width:  110,
		Height: 45,
	}, "Go back")

	if goBackPressed {
		sceneManager.SetScene(NewHome())
		ss.resetValues()
		return
	}

	drawInput(&ss.hostInput, 0)
	currentY := drawInput(&ss.portInput, 1)
	if !ss.isHost {
		currentY = drawInput(&ss.usernameInput, 2)
	}

	currentY += 40

	resetPressed := gui.Button(rl.Rectangle{
		X:      SCREEN_WIDTH/2 - BTN_WIDTH/2,
		Y:      currentY,
		Width:  BTN_WIDTH,
		Height: BTN_HEIGHT,
	}, "Reset")

	if resetPressed {
		ss.resetValues()
	}

	defaultsPressed := gui.Button(rl.Rectangle{
		X:      SCREEN_WIDTH/2 - BTN_WIDTH*1.5 - 20,
		Y:      currentY,
		Width:  BTN_WIDTH,
		Height: BTN_HEIGHT,
	}, "Defaults")

	if defaultsPressed {
		ss.defaultValues()
	}

	startPressed := gui.Button(rl.Rectangle{
		X:      SCREEN_WIDTH/2 + BTN_WIDTH/2 + 20,
		Y:      currentY,
		Width:  BTN_WIDTH,
		Height: BTN_HEIGHT,
	}, "Start")

	if startPressed {
		if !ss.checkInputsForErrors() {

			if ss.isHost {
				chat, err := NewHostChat(sceneManager, ss.hostInput.Value(), ss.portInput.Value(), "Host")
				if err != nil {
					ss.errorMsg.SetText(err.Error())
					return
				}
				sceneManager.SetScene(chat)
			} else {
				chat, err := NewClientChat(sceneManager, ss.hostInput.Value(), ss.portInput.Value(), ss.usernameInput.Value())
				if err != nil {
					ss.errorMsg.SetText(err.Error())
					return
				}
				sceneManager.SetScene(chat)
			}
		}
	}
	currentY += BTN_HEIGHT + 40

	ss.errorMsg.Draw(SCREEN_WIDTH/2-ss.errorMsg.GetWidth()/2, currentY, rl.Red)
}
