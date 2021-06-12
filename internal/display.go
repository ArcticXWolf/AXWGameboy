package internal

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const pixelSize int = 3

type Display struct {
	window  *pixelgl.Window
	picture *pixel.PictureData
}

func NewDisplay() *Display {
	cfg := pixelgl.WindowConfig{
		Title:  "AXWGameboy",
		Bounds: pixel.R(0, 0, float64(ScreenWidth*pixelSize), float64(ScreenHeight*pixelSize)),
	}
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	picture := &pixel.PictureData{
		Pix:    make([]color.RGBA, int(ScreenHeight*ScreenWidth)),
		Stride: int(ScreenWidth),
		Rect:   pixel.R(0, 0, float64(ScreenWidth), float64(ScreenHeight)),
	}

	return &Display{
		window:  window,
		picture: picture,
	}
}

func (d *Display) Render(gb *Gameboy) {
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			d.picture.Pix[(ScreenHeight-y-1)*ScreenWidth+x] = color.RGBA{
				R: gb.ReadyToRender[x][y][0],
				G: gb.ReadyToRender[x][y][1],
				B: gb.ReadyToRender[x][y][2],
				A: 255,
			}
		}
	}
	screenSprite := pixel.NewSprite(d.picture, pixel.R(0, 0, float64(ScreenWidth), float64(ScreenHeight)))

	d.window.Clear(color.White)
	screenSprite.Draw(d.window, pixel.IM.Moved(d.window.Bounds().Center()).Scaled(d.window.Bounds().Center(), float64(pixelSize)))

	d.window.Update()
}
