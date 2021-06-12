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
	gb := internal.NewGameboy()
	gb.Debugger.Enabled = true
	gb.Debugger.Address = 0x0100
	gb.Run()
}
