package internal

import (
	"path"
	"time"
)

var (
	FramesPerSecond int     = 60
	ScreenHeight    float64 = 144
	ScreenWidth     float64 = 160
)

type Gameboy struct {
	Display  *Display
	Cpu      *Cpu
	Debugger *Debugger
}

func NewGameboy() *Gameboy {
	c := NewCpu()
	c.Memory, _ = NewFromRom(string(path.Join(path.Base("."), "cpu_instrs.gb")), c.Gpu)

	d := NewDisplay()

	return &Gameboy{
		Cpu:      c,
		Display:  d,
		Debugger: &Debugger{Enabled: false},
	}
}

func (gb *Gameboy) Run() {
	cyclesPerFrame := ClockSpeed / FramesPerSecond
	frameDuration := time.Second / time.Duration(FramesPerSecond)
	frameCount := 0

	ticker := time.NewTicker(frameDuration)

	for ; true; <-ticker.C {
		frameCount++

		for i := 0; i < cyclesPerFrame; i++ {
			gb.Cpu.Tick(gb.Debugger)
		}
	}
}
