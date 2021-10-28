package gameview

import (
	"bytes"
	_ "embed"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
)

//go:embed osb.png
var osb []byte

func (ag *AXWGameboyEbitenGameView) loadOSBBackground() (err error) {
	ag.osbImg, err = png.Decode(bytes.NewReader(osb))
	if err != nil {
		return err
	}
	ag.osbData = make([]byte, 0)
	bounds := ag.osbImg.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := ag.osbImg.At(x, y).RGBA()
			ag.osbData = append(ag.osbData, byte(r), byte(g), byte(b), byte(a))
		}
	}
	return nil
}

type OnScreenButton struct {
	xMin int
	xMax int
	yMin int
	yMax int

	isMiscButton bool
	gameButton   internal.Button
	event        MiscEvent

	touched bool

	risingEdgeDelay  bool
	fallingEdgeDelay bool
}

func (a *AXWGameboyEbitenGameView) initOSB() {
	a.osbMap = []*OnScreenButton{
		{
			xMin:       125,
			yMin:       70,
			xMax:       150,
			yMax:       90,
			gameButton: internal.ButtonA,
		},
		{
			xMin:       100,
			yMin:       105,
			xMax:       125,
			yMax:       125,
			gameButton: internal.ButtonB,
		},
		{
			xMin:       85,
			yMin:       5,
			xMax:       125,
			yMax:       20,
			gameButton: internal.ButtonStart,
		},
		{
			xMin:       30,
			yMin:       5,
			xMax:       70,
			yMax:       20,
			gameButton: internal.ButtonSelect,
		},
		{
			xMin:       30,
			yMin:       55,
			xMax:       60,
			yMax:       85,
			gameButton: internal.ButtonUp,
		},
		{
			xMin:       60,
			yMin:       80,
			xMax:       90,
			yMax:       110,
			gameButton: internal.ButtonRight,
		},
		{
			xMin:       30,
			yMin:       110,
			xMax:       60,
			yMax:       140,
			gameButton: internal.ButtonDown,
		},
		{
			xMin:       5,
			yMin:       80,
			xMax:       35,
			yMax:       110,
			gameButton: internal.ButtonLeft,
		},
		{
			xMin:         90,
			yMin:         50,
			xMax:         110,
			yMax:         70,
			isMiscButton: true,
			event:        SpeedboostToggle,
		},
		{
			xMin:         0,
			yMin:         0,
			xMax:         20,
			yMax:         20,
			isMiscButton: true,
			event:        PauseToggle,
		},
		// {
		// 	xMin:         140,
		// 	yMin:         0,
		// 	xMax:         160,
		// 	yMax:         20,
		// 	isMiscButton: true,
		// 	event:        ShutdownGame,
		// },
	}
}

func (a *AXWGameboyEbitenGameView) handleOSBInputs() []MiscEvent {
	var events []MiscEvent
	tids := ebiten.TouchIDs()
	for _, osb := range a.osbMap {
		osb.touched = false
		for _, tid := range tids {
			x, y := ebiten.TouchPosition(tid)
			if (x >= osb.xMin && x <= osb.xMax) && (y >= osb.yMin+internal.ScreenHeight && y <= osb.yMax+internal.ScreenHeight) {
				osb.touched = true
			}
		}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if (x >= osb.xMin && x <= osb.xMax) && (y >= osb.yMin+internal.ScreenHeight && y <= osb.yMax+internal.ScreenHeight) {
				osb.touched = true
			}
		}

		if osb.detectRisingEdge(osb.touched) {
			if osb.isMiscButton {
				events = append(events, osb.event)
			} else {
				a.Gameboy.Inputs.ButtonsPressed = append(a.Gameboy.Inputs.ButtonsPressed, osb.gameButton)
			}
		}
		if osb.detectFallingEdge(osb.touched) {
			// Do not handle MiscButton release
			if !osb.isMiscButton {
				a.Gameboy.Inputs.ButtonsReleased = append(a.Gameboy.Inputs.ButtonsReleased, osb.gameButton)
			}
		}
	}
	return events
}

func (osb *OnScreenButton) detectFallingEdge(signal bool) bool {
	result := !signal && osb.fallingEdgeDelay
	osb.fallingEdgeDelay = signal
	return result
}

func (osb *OnScreenButton) detectRisingEdge(signal bool) bool {
	result := signal && !osb.risingEdgeDelay
	osb.risingEdgeDelay = signal
	return result
}
