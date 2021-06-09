package main

import (
	"log"
	"time"

	"go.janniklasrichter.de/axwgameboy/internal/cpu"
)

var (
	version = "dev"
	date    = "dev"
	commit  = "dev"
)

func main() {
	log.Printf("AXWGameboy | Version %v | Builddate %v | Commit %v", version, date, commit)
	cpu := cpu.New()

	framesPerSecond := 60
	frameDuration := time.Second / time.Duration(framesPerSecond)
	frameCount := 0

	ticker := time.NewTicker(frameDuration)

	for ; true; <-ticker.C {
		frameCount++

		cpu.Tick()
	}
}
