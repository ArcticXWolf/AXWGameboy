package ebitenprovider

import (
	"errors"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type AXWGameboyEbitenGame struct {
	Gameboy            *internal.Gameboy
	isPaused           bool
	isSpeedboostActive bool
	isTerminated       bool
	isOSBEnabled       bool
	osbData            []byte
	osbImg             image.Image
	osbMap             []*OnScreenButton
}

var Terminated = errors.New("terminated")

func NewAXWGameboyEbitenGame(gb *internal.Gameboy, enableOSB bool) *AXWGameboyEbitenGame {
	ag := &AXWGameboyEbitenGame{
		Gameboy:      gb,
		isOSBEnabled: enableOSB,
	}
	if ag.isOSBEnabled {
		err := ag.loadOSBBackground()
		if err != nil {
			log.Panic(err)
		}
		ag.initOSB()
	}
	return ag
}

func (a *AXWGameboyEbitenGame) Update() error {
	a.handleKeyboardInputs()
	events := a.handleOSBInputs()
	events = append(events, a.handleKeyboardInputsForMiscEvents()...)
	a.handleMiscEvents(events)

	if a.isTerminated {
		if a.Gameboy.Options.SavePath != "" {
			a.Gameboy.Memory.Cartridge.SaveRam(a.Gameboy.Options.SavePath)
		}
		return Terminated
	}

	if !a.isPaused {
		a.Gameboy.UpdateFrame(int(float32(internal.ClockSpeed) / 60.0 * a.Gameboy.Cpu.SpeedBoost))
	} else {
		a.Gameboy.Inputs.ClearButtonList()
	}
	return nil
}

func (a *AXWGameboyEbitenGame) Draw(screen *ebiten.Image) {
	bytes := a.Gameboy.GetReadyFramebufferAsBytearray()
	if a.isOSBEnabled {
		bytes = append(bytes, a.osbData...)
	}
	screen.ReplacePixels(bytes)
}

func (a *AXWGameboyEbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if a.isOSBEnabled {
		bounds := a.osbImg.Bounds()
		return internal.ScreenWidth, internal.ScreenHeight + bounds.Size().Y
	}
	return internal.ScreenWidth, internal.ScreenHeight
}
