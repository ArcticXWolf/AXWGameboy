package main

import (
	"log"
	"path"
	"time"

	"go.janniklasrichter.de/axwgameboy/internal/cpu"
	"go.janniklasrichter.de/axwgameboy/internal/memory"
)

var (
	version = "dev"
	date    = "dev"
	commit  = "dev"
)

func main() {
	log.Printf("AXWGameboy | Version %v | Builddate %v | Commit %v", version, date, commit)
	debugger := &cpu.Debugger{
		Enabled: true,
		Address: 0x0000,
	}
	cpu := cpu.New()
	cpu.Memory, _ = memory.NewFromRom(string(path.Join(path.Base("."), "cpu_instrs.gb")), cpu.Gpu)

	framesPerSecond := 2
	clockSpeed := 4194304
	cyclesPerFrame := clockSpeed / framesPerSecond
	frameDuration := time.Second / time.Duration(framesPerSecond)
	frameCount := 0

	ticker := time.NewTicker(frameDuration)

	for ; true; <-ticker.C {
		frameCount++

		for i := 0; i < cyclesPerFrame; i++ {
			cpu.Tick(debugger)
		}

		log.Printf("%s", cpu.String())
	}
}
