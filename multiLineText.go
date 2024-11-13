package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MultiLineLine struct {
	text  string
	width float32
}

type MultiLineText struct {
	text     string
	maxWidth float32

	fontSize    float32
	fontSpacing float32

	Lines []MultiLineLine

	biggestWidth float32
}

func NewMultiLineText(text string, maxWidth float32, fontSize float32) MultiLineText {
	multi := MultiLineText{
		text:        text,
		maxWidth:    maxWidth,
		fontSize:    fontSize,
		fontSpacing: 1,
	}

	multi.calculateLayout()

	return multi
}

func (mt *MultiLineText) String() string {
  return mt.text
}

func (mt *MultiLineText) SetText(text string) {
	mt.text = text
	mt.calculateLayout()
}

func (mt *MultiLineText) SetMaxWidth(maxWidth float32) {
	mt.maxWidth = maxWidth
	mt.calculateLayout()
}

func (mt *MultiLineText) SetFontSize(fontSize float32) {
	mt.fontSize = fontSize
	mt.calculateLayout()
}

func (mt *MultiLineText) calculateLayout() {
	totalWidth := float32(rl.MeasureText(mt.text, int32(mt.fontSize)))

	if totalWidth <= mt.maxWidth {
		mt.Lines = []MultiLineLine{{text: mt.text, width: totalWidth}}
		mt.biggestWidth = totalWidth
		return
	}

	arr := []MultiLineLine{}

	tmpText := ""

	var prevWidth float32 = 0
	var width float32 = 0
	for _, char := range mt.text {
		width = rl.MeasureTextEx(rl.GetFontDefault(), tmpText+string(char), mt.fontSize, mt.fontSpacing).X

		if width > mt.maxWidth {
			arr = append(arr, MultiLineLine{text: tmpText, width: width - prevWidth})
			tmpText = ""

			if width-prevWidth > mt.biggestWidth {
				mt.biggestWidth = width - prevWidth
			}

			prevWidth = width
		}

		tmpText += string(char)
	}

	if len(tmpText) > 0 {
		if width-prevWidth > mt.biggestWidth {
			mt.biggestWidth = width - prevWidth
		}
		arr = append(arr, MultiLineLine{text: tmpText, width: width})
	}

	mt.Lines = arr
}

func (mt *MultiLineText) Clear() {
	mt.text = ""
	mt.biggestWidth = 0
	mt.Lines = []MultiLineLine{}
}

func (mt *MultiLineText) GetWidth() float32 {
	return mt.biggestWidth
}

func (mt *MultiLineText) GetHeight() float32 {
	return float32(len(mt.Lines)) * mt.fontSize
}

func (mt *MultiLineText) Draw(x, y float32, color rl.Color) {
	for i, line := range mt.Lines {
		rl.DrawTextEx(rl.GetFontDefault(), line.text, rl.Vector2{X: x, Y: y + float32(i)*mt.fontSize}, mt.fontSize, mt.fontSpacing, color)
	}
}
