package internal

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var pixelSize float64 = 3

type Display struct {
	window  *pixelgl.Window
	picture *pixel.PictureData
}

func NewDisplay() *Display {
	cfg := pixelgl.WindowConfig{
		Title:  "AXWGameboy",
		Bounds: pixel.R(0, 0, ScreenWidth*pixelSize, ScreenHeight*pixelSize),
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	picture := &pixel.PictureData{
		Pix:    make([]color.RGBA, int(ScreenHeight*ScreenWidth)),
		Stride: int(ScreenWidth),
		Rect:   pixel.R(0, 0, ScreenWidth, ScreenHeight),
	}

	return &Display{
		window:  window,
		picture: picture,
	}
}
