package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
	debugCui "go.janniklasrichter.de/axwgameboy/pkg/cui"
	"go.janniklasrichter.de/axwgameboy/pkg/ebitenprovider"
)

var (
	version      = "dev"
	date         = "dev"
	commit       = "dev"
	savePath     string
	romPath      string
	paletteName  string
	headless     bool
	serialOutput bool
	cuiEnabled   bool
)

func init() {
	flag.StringVar(&savePath, "save", "", "Savefile to use")
	flag.StringVar(&romPath, "rom", "", "Rom to use")
	flag.StringVar(&paletteName, "palette", "white", "Name of a palette to use")
	flag.BoolVar(&serialOutput, "serial", false, "Show serial output in console")
	flag.BoolVar(&cuiEnabled, "cui", false, "Enable debug console interface")
}

func main() {
	flag.Parse()
	start()
}

func start() {
	log.Printf("AXWGameboy | Version %v | Builddate %v | Commit %v", version, date, commit)

	options := &internal.GameboyOptions{
		RomPath:  romPath,
		SavePath: savePath,
		Palette:  paletteName,
	}
	if serialOutput {
		options.SerialOutputFunction = func(b byte) {
			log.Print(string(b))
		}
	}
	var cui *debugCui.GameboyCui
	if cuiEnabled {
		cui = debugCui.NewGui()
		options.OnFrameFunction = func(g *internal.Gameboy) {
			cui.UpdateView(g)
		}
		go cui.RunLoop()
	}

	gb, err := internal.NewGameboy(options)
	if err != nil {
		log.Panicf("Error loading rom: %s", err)
	}
	log.Printf("Loaded %s", gb.Memory.Cartridge.CartridgeInfo())

	ebitenGame := ebitenprovider.NewAXWGameboyEbitenGame(gb, true)
	ebiten.SetWindowResizable(true)
	ebiten.RunGame(ebitenGame)
}
