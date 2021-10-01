package viewcontainer

import (
	"log"
	"syscall/js"

	"go.janniklasrichter.de/axwgameboy/internal"
	"go.janniklasrichter.de/axwgameboy/pkg/gameview"
)

var (
	window       = js.Global().Get("window")
	document     = js.Global().Get("document")
	gameSelector js.Value
	fileInput    js.Value
)

func (a *AXWGameboyViewContainer) createGameSelector() {
	gameSelector = document.Call("createElement", "div")
	gameSelector.Set("id", "gameSelector")
	document.Get("body").Call("appendChild", gameSelector)

	gameSelectorStyle := gameSelector.Get("style")
	gameSelectorStyle.Set("width", "100%")
	gameSelectorStyle.Set("height", "100%")
	gameSelectorStyle.Set("margin", "0")
	gameSelectorStyle.Set("padding", "0")
	gameSelectorStyle.Set("background", "#ffffff")
	gameSelectorStyle.Set("position", "absolute")
	gameSelectorStyle.Set("top", "0")

	fileInput = document.Call("createElement", "input")
	fileInput.Set("type", "file")
	fileInput.Set("id", "fileInput")

	fileInputLabel := document.Call("createElement", "label")
	fileInputLabel.Set("innerHTML", "Open ROM")
	fileInputLabel.Call("setAttribute", "for", "fileInput")
	fileInputLabel.Get("style").Set("border", "1px solid #ccc")
	fileInputLabel.Get("style").Set("display", "inline-block")
	fileInputLabel.Get("style").Set("padding", "6px 12px")
	fileInputLabel.Get("style").Set("cursor", "pointer")
	fileInputLabel.Get("style").Set("width", "300px")
	fileInputLabel.Get("style").Set("height", "50px")
	fileInputLabel.Get("style").Set("margin", "-25px 0px 0px -150px")
	fileInputLabel.Get("style").Set("position", "absolute")
	fileInputLabel.Get("style").Set("top", "50%")
	fileInputLabel.Get("style").Set("left", "50%")
	fileInputLabel.Get("style").Set("text-align", "center")
	fileInputLabel.Get("style").Set("font-size", "40px")

	gameSelector.Call("appendChild", fileInputLabel)
	gameSelector.Call("appendChild", fileInput)

	fileInput.Get("style").Set("display", "none")

	fileInput.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.Global().Call("readFileInputToFunction", "fileInput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			array := args[0]
			buf := make([]byte, array.Get("length").Int())
			js.CopyBytesToGo(buf, array)
			a.handleGameSelector(buf)
			return nil
		}))
		return nil
	}))
}

func (a *AXWGameboyViewContainer) handleGameSelector(romData []byte) {
	log.Printf("got ROM.")
	a.destroyGameSelector()

	options := &internal.GameboyOptions{
		RomData:      romData,
		OSBEnabled:   true,
		SoundEnabled: true,
		SoundVolume:  1.0,
		CGBEnabled:   true,
	}

	a.GameView = gameview.NewAXWGameboyEbitenGameView(options)
	a.currentStage = GameRunning
}

func (a *AXWGameboyViewContainer) destroyGameSelector() {
	document.Get("body").Call("removeChild", gameSelector)
}
