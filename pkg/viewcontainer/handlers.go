package viewcontainer

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"syscall/js"

	"go.janniklasrichter.de/axwgameboy/internal"
	"go.janniklasrichter.de/axwgameboy/pkg/gameview"
)

type loadRomData struct {
	romData      []byte
	SoundEnabled bool
	OSBEnabled   bool
}

func (ag *AXWGameboyViewContainer) installRomLoader() {
	romDataChannel := make(chan loadRomData)

	js.Global().Set("loadROM", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var rom []byte
		for _, el := range args[0].String() {
			rom = append(rom, byte(el))
		}
		romDataChannel <- loadRomData{
			romData:      rom,
			SoundEnabled: args[1].Get("soundEnabled").Bool(),
			OSBEnabled:   args[1].Get("osbEnabled").Bool(),
		}
		return nil
	}))

	go func() {
		for {
			select {
			case loadData := <-romDataChannel:
				romHash := sha256.Sum256(loadData.romData)
				savegame := &localStorageSavegame{
					romHash: fmt.Sprintf("%x", romHash),
				}
				options := &internal.GameboyOptions{
					RomReader:    bytes.NewReader(loadData.romData),
					SaveWriter:   savegame,
					SaveReader:   savegame,
					CGBEnabled:   true,
					SoundEnabled: loadData.SoundEnabled,
					OSBEnabled:   loadData.OSBEnabled,
					SoundVolume:  0.1,
				}
				ag.GameView = gameview.NewAXWGameboyEbitenGameView(options)
				ag.currentStage = GameRunning
			}
		}
	}()

	return
}

func (ag *AXWGameboyViewContainer) installSettingsHandler() {
	//settingsDataChannel := make(chan []byte)

	js.Global().Set("changeSettings", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if ag.GameView == nil || ag.GameView.Gameboy == nil {
			return nil
		}
		ag.GameView.Gameboy.Options.SoundEnabled = args[0].Get("soundEnabled").Bool()
		ag.GameView.Gameboy.Options.OSBEnabled = args[0].Get("osbEnabled").Bool()
		ag.GameView.Gameboy.CheatCodeManager.ReplaceCodeList(strings.NewReader(args[0].Get("cheats").String()))
		return nil
	}))
}
