package main

import (
	"log"

	"github.com/faiface/pixel/pixelgl"
	"go.janniklasrichter.de/axwgameboy/internal"
)

var (
	version = "dev"
	date    = "dev"
	commit  = "dev"
)

func main() {
	pixelgl.Run(start)
}

func start() {
	log.Printf("AXWGameboy | Version %v | Builddate %v | Commit %v", version, date, commit)
	options := &internal.GameboyOptions{
		RomPath: "./roms/01.gb",
		SerialOutputFunction: func(b byte) {
			log.Printf("Got serial output: %s", string(b))
		},
	}
	gb, err := internal.NewGameboy(options)
	if err != nil {
		log.Panicf("Error loading rom: %s", err)
	}

	gb.Debugger.AddressEnabled = false
	gb.Debugger.Address = 0x0100
	gb.Debugger.LogOnly = false
	gb.Debugger.LogEvery = 10

	gb.Run()
}
