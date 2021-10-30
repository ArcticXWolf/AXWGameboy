package viewcontainer

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"syscall/js"

	"go.janniklasrichter.de/axwgameboy/internal"
	"go.janniklasrichter.de/axwgameboy/pkg/gameview"
)

type loadRomData struct {
	romData      []byte
	SoundEnabled bool
	OSBEnabled   bool
	CheatCodes   string
}

func (ag *AXWGameboyViewContainer) installRomLoader() {
	js.Global().Set("loadROM", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var rom []byte
		for _, el := range args[0].String() {
			rom = append(rom, byte(el))
		}
		loadData := loadRomData{
			romData:      rom,
			SoundEnabled: args[1].Get("soundEnabled").Bool(),
			OSBEnabled:   args[1].Get("osbEnabled").Bool(),
			CheatCodes:   args[1].Get("cheats").String(),
		}

		ag.startGameboyFromLoadData(loadData)
		ag.changeSettingsFromLoadData(loadData)

		return nil
	}))

	return
}

func (ag *AXWGameboyViewContainer) installSettingsHandler() {
	js.Global().Set("changeSettings", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		loadData := loadRomData{
			SoundEnabled: args[0].Get("soundEnabled").Bool(),
			OSBEnabled:   args[0].Get("osbEnabled").Bool(),
			CheatCodes:   args[0].Get("cheats").String(),
		}

		ag.changeSettingsFromLoadData(loadData)
		return nil
	}))
}

func (ag *AXWGameboyViewContainer) installSavegameHandler() {
	js.Global().Set("loadSave", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if ag.GameView == nil || ag.GameView.Gameboy == nil {
			return nil
		}

		ag.currentStage = GameSelection

		var save []byte
		for _, el := range args[0].String() {
			save = append(save, byte(el))
		}
		loadData := loadRomData{
			romData:      ag.GameView.Gameboy.Memory.Cartridge.GetBinaryData(),
			SoundEnabled: ag.GameView.Gameboy.Options.SoundEnabled,
			OSBEnabled:   ag.GameView.Gameboy.Options.OSBEnabled,
			CheatCodes:   ag.GameView.Gameboy.CheatCodeManager.GetCodeList(),
		}

		romHash := sha256.Sum256(loadData.romData)
		savegame := &localStorageSavegame{
			romHash: fmt.Sprintf("%x", romHash),
		}
		_, err := savegame.Write(save)
		if err != nil {
			log.Printf("replacement of savegame failed")
			return nil
		}

		ag.startGameboyFromLoadData(loadData)
		ag.changeSettingsFromLoadData(loadData)

		return nil
	}))

	js.Global().Set("downloadSave", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if ag.GameView == nil || ag.GameView.Gameboy == nil {
			return nil
		}

		romHash := sha256.Sum256(ag.GameView.Gameboy.Memory.Cartridge.GetBinaryData())
		savegame := &localStorageSavegame{
			romHash: fmt.Sprintf("%x", romHash),
		}
		save, err := ioutil.ReadAll(savegame)
		if err != nil {
			log.Printf("reading of savegame failed")
			return nil
		}

		dstArray := js.Global().Get("Uint8Array").New(len(save))
		js.CopyBytesToJS(dstArray, save)

		return dstArray
	}))
}

func (ag *AXWGameboyViewContainer) startGameboyFromLoadData(loadData loadRomData) {
	romHash := sha256.Sum256(loadData.romData)
	savegame := &localStorageSavegame{
		romHash: fmt.Sprintf("%x", romHash),
	}
	options := &internal.GameboyOptions{
		RomReader:   bytes.NewReader(loadData.romData),
		SaveWriter:  savegame,
		SaveReader:  savegame,
		CGBEnabled:  true,
		SoundVolume: 0.1,
	}
	ag.GameView = gameview.NewAXWGameboyEbitenGameView(options)
	ag.currentStage = GameRunning
}

func (ag *AXWGameboyViewContainer) changeSettingsFromLoadData(loadData loadRomData) {
	if ag.GameView == nil || ag.GameView.Gameboy == nil {
		return
	}
	ag.GameView.Gameboy.Options.SoundEnabled = loadData.SoundEnabled
	ag.GameView.Gameboy.Options.OSBEnabled = loadData.OSBEnabled
	ag.GameView.Gameboy.CheatCodeManager.ReplaceCodeList(strings.NewReader(loadData.CheatCodes))
}
