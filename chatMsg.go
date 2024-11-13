package main

import (
	"encoding/binary"
	"io"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MSG_MAX_WIDTH       = SCREEN_WIDTH/2 - 50
	MSG_FONT_SIZE       = 18
	MSG_USERNAME_BYTES  = 24
	MSG_TIMESTEMP_BYTES = 8
	MSG_BYTES           = 250
)

type ChatMsg struct {
	username  string
	message   MultiLineText
	isSender  bool
	isSystem  bool
	isUnknown bool
	timestamp int64
}

func NewChatMsgFromBytes(bytes []byte) ChatMsg {
	const (
		P1    = MSG_USERNAME_BYTES
		P2    = MSG_TIMESTEMP_BYTES
		P3    = MSG_BYTES
		TOTAL = P1 + P2 + P3
	)

	if len(bytes) != TOTAL {
		return ChatMsg{
			isUnknown: true,
			username:  "Unknown",
			message:   NewMultiLineText(string(bytes), MSG_MAX_WIDTH, MSG_FONT_SIZE),
			timestamp: time.Now().UnixMilli(),
		}
	}

	return ChatMsg{
		username:  string(bytes[:P1]),
		timestamp: int64(binary.LittleEndian.Uint64(bytes[P1 : P1+P2])),
		message:   NewMultiLineText(string(bytes[P1+P2:]), MSG_MAX_WIDTH, MSG_FONT_SIZE),
	}
}

func (cm *ChatMsg) Bytes() []byte {
	const (
		P1 = MSG_USERNAME_BYTES
		P2 = MSG_TIMESTEMP_BYTES
		P3 = MSG_BYTES
	)

	bytesMsg := make([]byte, P1+P2+P3)

	copy(bytesMsg[:P1], cm.username)
	binary.LittleEndian.PutUint64(bytesMsg[P1:P1+P2], uint64(time.Now().UnixMilli()))
	copy(bytesMsg[P1+P2:], cm.message.String())

	return bytesMsg
}

func (cm *ChatMsg) WriteTo(writer io.Writer) (n int, err error) {
	return writer.Write(cm.Bytes())
}

func (cm *ChatMsg) Draw(chatBounds rl.Rectangle, yOffset float32) float32 {
	const (
		INNER_PAD          = 5
		MSG_GAP            = 10
		USERNAME_FONT_SIZE = 12
		TIME_FONT_SIZE     = 10
		TIMESTEMP_MT       = 4
	)

	msgWidth := cm.message.GetWidth() + INNER_PAD*2
	msgHeight := cm.message.GetHeight() + INNER_PAD*2
	yOffset += msgHeight

	msgX := chatBounds.X
	msgY := chatBounds.Y + chatBounds.Height - yOffset - USERNAME_FONT_SIZE - TIMESTEMP_MT
	usernameX := msgX
	timeX := msgX

	bgColor := rl.Beige
	username := cm.username
	timeStr := time.UnixMilli(cm.timestamp).Format(time.TimeOnly)

	if cm.isSystem {
		bgColor = rl.NewColor(255, 125, 125, 255)
	} else if cm.isSender {
		bgColor = rl.SkyBlue
		msgX = chatBounds.X + chatBounds.Width - msgWidth
		usernameX = chatBounds.X + chatBounds.Width - float32(rl.MeasureText(username, USERNAME_FONT_SIZE))
		timeX = chatBounds.X + chatBounds.Width - float32(rl.MeasureText(timeStr, TIME_FONT_SIZE))
	} else if cm.isUnknown {
		bgColor = rl.NewColor(242, 235, 27, 255)
	}

	rl.DrawRectangleV(rl.Vector2{X: msgX, Y: msgY}, rl.Vector2{X: msgWidth, Y: msgHeight}, bgColor)

	// username
	rl.DrawText(username, int32(usernameX), int32(msgY)-12, 12, rl.Black)

	// timestamp
	rl.DrawText(timeStr, int32(timeX), int32(msgY+msgHeight+4), 10, rl.Black)

	cm.message.Draw(msgX+INNER_PAD, msgY+INNER_PAD, rl.Black)

	return yOffset + MSG_GAP + TIME_FONT_SIZE + TIMESTEMP_MT + USERNAME_FONT_SIZE
}
