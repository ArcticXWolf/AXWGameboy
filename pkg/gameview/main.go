package gameview

import (
	"errors"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
	"go.janniklasrichter.de/axwgameboy/pkg/gameview/debugscreens"
)

type AXWGameboyEbitenGameView struct {
	Gameboy            *internal.Gameboy
	isPaused           bool
	isSpeedboostActive bool
	isTerminated       bool
	isOSBRequested     bool
	isOSBEnabled       bool
	isDebugRequested   bool
	isDebugEnabled     bool
	debugScreens       []debugscreens.DebugScreen
	osbData            []byte
	osbImg             image.Image
	osbMap             []*OnScreenButton
}

var Terminated = errors.New("terminated")

func NewAXWGameboyEbitenGameView(options *internal.GameboyOptions) *AXWGameboyEbitenGameView {
	gb, err := internal.NewGameboy(options)
	if err != nil {
		log.Panicf("Error loading rom: %s", err)
	}
	log.Printf("Loaded %s", gb.Memory.Cartridge.CartridgeInfo())

	ag := &AXWGameboyEbitenGameView{
		Gameboy:        gb,
		isOSBRequested: options.OSBEnabled,
		debugScreens: []debugscreens.DebugScreen{
			&debugscreens.Tilemap{
				X: 170,
				Y: 0,
			},
			&debugscreens.PaletteList{
				X: 300,
				Y: 0,
			},
			&debugscreens.BgMap{
				X:     350,
				Y:     0,
				MapId: 0,
			},
			&debugscreens.BgMap{
				X:     350,
				Y:     260,
				MapId: 1,
			},
			&debugscreens.SpriteList{
				X: 0,
				Y: 300,
			},
		},
	}

	err = ag.loadOSBBackground()
	if err != nil {
		log.Panic(err)
	}
	ag.initOSB()

	return ag
}

func (a *AXWGameboyEbitenGameView) Update() error {
	a.handleKeyboardInputs()
	events := a.handleOSBInputs()
	events = append(events, a.handleKeyboardInputsForMiscEvents()...)
	a.handleMiscEvents(events)

	if a.isTerminated {
		if a.Gameboy.Options.SaveWriter != nil {
			a.Gameboy.Memory.Cartridge.SaveRam(a.Gameboy.Options.SaveWriter)
		}
		return Terminated
	}

	if !a.isPaused {
		a.Gameboy.Update(int(float32(internal.ClockSpeed) / 60.0 * a.Gameboy.Cpu.SpeedBoost * a.Gameboy.GetSpeedMultiplier()))
	} else {
		a.Gameboy.Inputs.ClearButtonList()
	}

	a.isOSBRequested = a.Gameboy.Options.OSBEnabled

	ebiten.SetWindowTitle(fmt.Sprintf("%s (%.02f TPS)", a.Gameboy.Memory.Cartridge.CartridgeHeader().Title, ebiten.CurrentTPS()))

	return nil
}

func (a *AXWGameboyEbitenGameView) Draw(screen *ebiten.Image) {
	gamescreen := screen.SubImage(image.Rect(0, 0, internal.ScreenWidth, internal.ScreenHeight)).(*ebiten.Image)
	gamescreen.ReplacePixels(a.Gameboy.GetReadyFramebufferAsBytearray())

	if a.isOSBEnabled {
		bounds := a.osbImg.Bounds()
		osbscreen := screen.SubImage(image.Rect(0, internal.ScreenHeight, bounds.Size().X, internal.ScreenHeight+bounds.Size().Y)).(*ebiten.Image)
		osbscreen.ReplacePixels(a.osbData)
	}

	if a.isDebugEnabled {
		for i := 0; i < len(a.debugScreens); i++ {
			debugscreen := a.debugScreens[i]
			debugscreen.Draw(a.Gameboy, screen)
		}
	}
}

func (a *AXWGameboyEbitenGameView) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if a.isOSBEnabled != a.isOSBRequested {
		a.isOSBEnabled = a.isOSBRequested
	}
	if a.isDebugEnabled != a.isDebugRequested {
		a.isDebugEnabled = a.isDebugRequested
	}

	if a.isDebugEnabled {
		return 700, 600
	}

	if a.isOSBEnabled {
		bounds := a.osbImg.Bounds()
		return internal.ScreenWidth, internal.ScreenHeight + bounds.Size().Y
	}
	return internal.ScreenWidth, internal.ScreenHeight
}

func (a *AXWGameboyEbitenGameView) GetExitResult() []byte {
	return []byte{}
}
