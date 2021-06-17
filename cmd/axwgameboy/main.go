package main

import (
	"flag"
	"log"

	"github.com/faiface/pixel/pixelgl"
	"go.janniklasrichter.de/axwgameboy/internal"
	debugCui "go.janniklasrichter.de/axwgameboy/pkg/cui"
)

var (
	version      = "dev"
	date         = "dev"
	commit       = "dev"
	romPath      string
	headless     bool
	serialOutput bool
	cuiEnabled   bool
)

func init() {
	flag.StringVar(&romPath, "rom", "./cpu_instrs.gb", "Rom to use")
	flag.BoolVar(&headless, "headless", false, "Run in headless (aka no display) mode")
	flag.BoolVar(&serialOutput, "serial", false, "Show serial output in console")
	flag.BoolVar(&cuiEnabled, "cui", false, "Enable debug console interface")
}

func main() {
	flag.Parse()
	pixelgl.Run(start)
}

func start() {
	log.Printf("AXWGameboy | Version %v | Builddate %v | Commit %v", version, date, commit)

	options := &internal.GameboyOptions{
		RomPath:  romPath,
		Headless: headless,
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

	gb.Debugger.AddressEnabled = true
	gb.Debugger.Address = 0x0150
	gb.Debugger.LogOnly = false
	gb.Debugger.LogEvery = 10

	gb.Run()
}
