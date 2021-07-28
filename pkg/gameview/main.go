package gameview

import (
	"errors"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type AXWGameboyEbitenGameView struct {
	Gameboy            *internal.Gameboy
	isPaused           bool
	isSpeedboostActive bool
	isTerminated       bool
	isOSBRequested     bool
	isOSBEnabled       bool
	isTilemapRequested bool
	isTilemapEnabled   bool
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

func (a *AXWGameboyEbitenGameView) Draw(screen *ebiten.Image) {
	gamescreen := screen.SubImage(image.Rect(0, 0, internal.ScreenWidth, internal.ScreenHeight)).(*ebiten.Image)
	gamescreen.ReplacePixels(a.Gameboy.GetReadyFramebufferAsBytearray())

	if a.isOSBEnabled || a.isTilemapEnabled {
		bounds := a.osbImg.Bounds()
		osbscreen := screen.SubImage(image.Rect(0, internal.ScreenHeight, bounds.Size().X, internal.ScreenHeight+bounds.Size().Y)).(*ebiten.Image)
		if a.isTilemapEnabled {
			osbscreen.ReplacePixels(a.Gameboy.Gpu.GetTilemapAsBytearray(1))
		} else {
			osbscreen.ReplacePixels(a.osbData)
		}
	}
}

func (a *AXWGameboyEbitenGameView) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if a.isOSBEnabled != a.isOSBRequested {
		a.isOSBEnabled = a.isOSBRequested
	}
	if a.isTilemapEnabled != a.isTilemapRequested {
		a.isTilemapEnabled = a.isTilemapRequested
	}

	if a.isOSBEnabled || a.isTilemapEnabled {
		bounds := a.osbImg.Bounds()
		return internal.ScreenWidth, internal.ScreenHeight + bounds.Size().Y
	}
	return internal.ScreenWidth, internal.ScreenHeight
}

func (a *AXWGameboyEbitenGameView) GetExitResult() []byte {
	return []byte{}
}
