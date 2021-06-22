package internal

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const pixelSize int = 5

type Display struct {
	window  *pixelgl.Window
	picture *pixel.PictureData
}

func NewDisplay() *Display {
	cfg := pixelgl.WindowConfig{
		Title:  "AXWGameboy",
		Bounds: pixel.R(0, 0, float64(ScreenWidth*pixelSize), float64(ScreenHeight*pixelSize)),
		VSync:  true,
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

var keys = map[pixelgl.Button]Button{
	pixelgl.KeyZ:       ButtonA,
	pixelgl.KeyX:       ButtonB,
	pixelgl.KeyLeftAlt: ButtonSelect,
	pixelgl.KeySpace:   ButtonStart,
	pixelgl.KeyRight:   ButtonRight,
	pixelgl.KeyLeft:    ButtonLeft,
	pixelgl.KeyUp:      ButtonUp,
	pixelgl.KeyDown:    ButtonDown,
}

var gamepadKeys = map[pixelgl.GamepadButton]Button{
	pixelgl.ButtonA:         ButtonA,
	pixelgl.ButtonB:         ButtonB,
	pixelgl.ButtonBack:      ButtonSelect,
	pixelgl.ButtonStart:     ButtonStart,
	pixelgl.ButtonDpadRight: ButtonRight,
	pixelgl.ButtonDpadLeft:  ButtonLeft,
	pixelgl.ButtonDpadUp:    ButtonUp,
	pixelgl.ButtonDpadDown:  ButtonDown,
}

func (d *Display) HandleInput(gb *Gameboy) {
	if d.window.JustPressed(pixelgl.KeyEscape) {
		gb.Quit = true
	}
	if d.window.JustPressed(pixelgl.KeyD) {
		gb.Debugger.triggerBreakpoint(gb)
	}
	if d.window.JustPressed(pixelgl.KeyLeftShift) {
		gb.Cpu.SpeedBoost = 4.0
	} else if d.window.JustReleased(pixelgl.KeyLeftShift) {
		gb.Cpu.SpeedBoost = 1.0
	}

	for key, button := range keys {
		if d.window.JustPressed(key) {
			gb.Inputs.buttonsPressed = append(gb.Inputs.buttonsPressed, button)
		}
		if d.window.JustReleased(key) {
			gb.Inputs.buttonsReleased = append(gb.Inputs.buttonsReleased, button)
		}
	}

	for joystick := pixelgl.Joystick1; joystick <= pixelgl.JoystickLast; joystick++ {
		if present := d.window.JoystickPresent(joystick); present {
			for key, button := range gamepadKeys {
				if d.window.JoystickJustPressed(joystick, key) {
					gb.Inputs.buttonsPressed = append(gb.Inputs.buttonsPressed, button)
				}
				if d.window.JoystickJustReleased(joystick, key) {
					gb.Inputs.buttonsReleased = append(gb.Inputs.buttonsReleased, button)
				}
			}
		}
	}
}
