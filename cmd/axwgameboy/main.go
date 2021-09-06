package main

import (
	"flag"
	"log"
	"runtime"

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
	soundEnabled bool
	soundVolume  float64
	serialOutput bool
	colorEnabled bool
	cuiEnabled   bool
	osbEnabled   bool
)

func init() {
	flag.StringVar(&savePath, "save", "", "Enables RAM persistence and has to contain the path to the desired savefile.")
	flag.StringVar(&romPath, "rom", "", "Set to the path of the rom, which will be used. If not specified, then a file selector will be shown with all roms in the current folder.")
	flag.StringVar(&paletteName, "palette", "DMG", "For non-color-mode: Specify which color palette shall be used. Currently available: dmg, red, white")
	flag.BoolVar(&serialOutput, "serial", false, "Print bytes of serial output to console as ASCII characters.")
	flag.BoolVar(&cuiEnabled, "cui", false, "Disable normal console output and show console debug gui instead")
	flag.BoolVar(&colorEnabled, "color", true, "Defaults to true. If set to false, it forces all games to be in non-color-mode.")
	flag.BoolVar(&osbEnabled, "osb", false, "Enable on-screen-button display.")
	flag.BoolVar(&soundEnabled, "sound", true, "Defaults to true. Enable/Disable sound.")
	flag.Float64Var(&soundVolume, "volume", 0.5, "Sets the starting master volume. Specify between 0 and 1.")
}

func main() {
	flag.Parse()
	start()
}

func start() {
	log.Printf("AXWGameboy | Version %v | Builddate %v | Commit %v", version, date, commit)

	options := &internal.GameboyOptions{
		RomPath:      romPath,
		SavePath:     savePath,
		Palette:      paletteName,
		SoundEnabled: soundEnabled,
		SoundVolume:  soundVolume,
		OSBEnabled:   osbEnabled,
		CGBEnabled:   colorEnabled,
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

	if runtime.GOOS == "android" {
		options.OSBEnabled = true
	}

	ebitenGame := ebitenprovider.NewAXWGameboyEbitenGame(options)
	ebiten.SetWindowResizable(true)
	ebiten.RunGame(ebitenGame)
}
