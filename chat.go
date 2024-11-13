package main

import (
	"fmt"
	"math"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Chat struct {
	IsHost   bool
	address  string
	username string

	tcpServer *TcpServer
	tcpClient *TcpClient
	messages  []ChatMsg

	message  string
	isTyping bool

	scrollY      int32
	sceneManager *SceneManager
}

func NewChat() *Chat {
	chat := &Chat{
		messages: []ChatMsg{},
	}

	return chat
}

func NewHostChat(sceneManager *SceneManager, host, port, username string) (*Chat, error) {
	chat := &Chat{
		sceneManager: sceneManager,
		IsHost:       true,
		address:      fmt.Sprintf("%v:%v", host, port),
		username:     username,
	}

	chat.tcpServer = NewTcpServer(chat)

	if err := chat.tcpServer.Start(); err != nil {
		return nil, err
	}

	return chat, nil
}

func NewClientChat(sceneManager *SceneManager, host, port, username string) (*Chat, error) {
	chat := &Chat{
		sceneManager: sceneManager,
		address:      fmt.Sprintf("%v:%v", host, port),
		username:     username,
	}

	chat.tcpClient = NewTcpClient(chat)

	if err := chat.tcpClient.Start(); err != nil {
		return nil, err
	}

	return chat, nil
}

func (c *Chat) ReciveChatMsg(chatMsg ChatMsg) {
	c.messages = append([]ChatMsg{chatMsg}, c.messages...)
}

func (c *Chat) ReciveSystemMsg(msg string) {
	c.messages = append([]ChatMsg{{
		username:  "System",
		message:   NewMultiLineText(msg, MSG_MAX_WIDTH, MSG_FONT_SIZE),
		isSystem:  true,
		timestamp: time.Now().UnixMilli(),
	}}, c.messages...)
}

func (c *Chat) sendMsg(msg string) {
	chatMsg := ChatMsg{
		username:  c.username,
		message:   NewMultiLineText(msg, MSG_MAX_WIDTH, MSG_FONT_SIZE),
		isSender:  true,
		timestamp: time.Now().UnixMilli(),
	}

	if c.IsHost {
		c.tcpServer.BroadcastMessage(chatMsg, nil)
		c.messages = append([]ChatMsg{chatMsg}, c.messages...)
		return
	}

	c.tcpClient.SendMsg(chatMsg)

}

func (c *Chat) quit() {
	if c.tcpServer != nil {
		c.tcpServer.Close()
		c.tcpServer = nil
		return
	}

	if c.tcpClient != nil {
		c.tcpClient.Close()
		c.tcpClient = nil
		return
	}
}

func (c *Chat) UpdateScene(sceneManager *SceneManager) {
	const (
		SCRENN_PADDING  float32 = 10
		GAP                     = 10
		INPUT_WIDTH             = SCREEN_WIDTH - SCRENN_PADDING*2
		INPUT_HEIGHT            = 50
		BTN_WIDTH               = 90
		SCROLLBAR_WIDTH         = 15
		CHAT_WIDTH              = SCREEN_WIDTH - GAP*2 - SCROLLBAR_WIDTH
		CHAT_HEIGHT             = SCREEN_HEIGHT - SCRENN_PADDING - INPUT_HEIGHT - GAP
	)
	rl.ClearBackground(rl.White)

	c.scrollY -= int32(rl.GetMouseWheelMove())

	var chatYOffset float32 = -float32(c.scrollY) + GAP
	for _, chatMsg := range c.messages {
		chatYOffset = chatMsg.Draw(rl.Rectangle{
			X:      GAP,
			Y:      0,
			Width:  CHAT_WIDTH,
			Height: CHAT_HEIGHT,
		}, chatYOffset)
	}

	var totalMsgsHeight float32 = chatYOffset + float32(c.scrollY) - GAP
	var scrollArea int32 = 0
	if totalMsgsHeight > CHAT_HEIGHT {
		scrollArea = int32(math.Abs(float64(CHAT_HEIGHT-totalMsgsHeight))) + 75
	}

	c.scrollY = gui.ScrollBar(rl.Rectangle{X: SCREEN_WIDTH - SCROLLBAR_WIDTH, Y: 0, Width: SCROLLBAR_WIDTH, Height: CHAT_HEIGHT}, c.scrollY, 0, scrollArea)

	// bottom panel bg
	rl.DrawRectangleV(rl.Vector2{X: 0, Y: CHAT_HEIGHT}, rl.Vector2{X: SCREEN_WIDTH, Y: INPUT_HEIGHT + GAP + SCRENN_PADDING}, rl.White)
	// bottom panel border top
	rl.DrawRectangleV(rl.Vector2{X: 0, Y: CHAT_HEIGHT}, rl.Vector2{X: SCREEN_WIDTH, Y: 1}, rl.LightGray)

	gui.SetStyle(gui.DEFAULT, gui.TEXT_SIZE, 20)
	inputToggle := gui.TextBox(rl.Rectangle{
		X:      SCREEN_WIDTH/2 - INPUT_WIDTH/2,
		Y:      SCREEN_HEIGHT - INPUT_HEIGHT - SCRENN_PADDING,
		Width:  INPUT_WIDTH - BTN_WIDTH - GAP,
		Height: INPUT_HEIGHT,
	}, &c.message, MSG_BYTES, c.isTyping)

	sentMsgPressed := gui.Button(rl.Rectangle{
		X:      SCREEN_WIDTH - BTN_WIDTH - SCRENN_PADDING,
		Y:      SCREEN_HEIGHT - INPUT_HEIGHT - SCRENN_PADDING,
		Width:  BTN_WIDTH,
		Height: INPUT_HEIGHT,
	}, "Send")

	quitPressed := gui.Button(rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  60,
		Height: 40,
	}, "Quit")

	if quitPressed {
		c.quit()
		sceneManager.SetScene(NewHome())
		return
	}

	if (sentMsgPressed || inputToggle && c.isTyping) && len(c.message) > 0 {
		c.sendMsg(c.message)
		c.message = ""
		return
	}

	if inputToggle {
		c.isTyping = !c.isTyping
	}
}
