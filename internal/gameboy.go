package internal

import (
	"path"
	"time"
)

const (
	FramesPerSecond int = 60
	ScreenHeight    int = 144
	ScreenWidth     int = 160
)

type Gameboy struct {
	Display       *Display
	Cpu           *Cpu
	Memory        MemoryDevice
	Gpu           *Gpu
	Debugger      *Debugger
	WorkingScreen [ScreenWidth][ScreenHeight][3]uint8
	ReadyToRender [ScreenWidth][ScreenHeight][3]uint8
}

func NewGameboy() *Gameboy {
	g := NewGpu()
	c := NewCpu()
	m, _ := NewFromRom(string(path.Join(path.Base("."), "cpu_instrs.gb")), g)
	d := NewDisplay()

	return &Gameboy{
		Cpu:      c,
		Display:  d,
		Memory:   m,
		Gpu:      g,
		Debugger: &Debugger{Enabled: false},
	}
}

func (gb *Gameboy) Run() {
	cyclesPerFrame := int(float32(ClockSpeed) / float32(FramesPerSecond) * SpeedBoost)
	frameDuration := time.Second / time.Duration(FramesPerSecond)
	frameCount := 0

	ticker := time.NewTicker(frameDuration)

	for ; true; <-ticker.C {
		frameCount++

		for i := 0; i < cyclesPerFrame; i++ {
			cycles := gb.Cpu.Tick(gb)
			gb.Gpu.Update(gb, cycles)
		}

		gb.Display.Render(gb)
	}
}
