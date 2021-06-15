package internal

import (
	"fmt"
	"time"
)

const (
	FramesPerSecond int = 60
	ScreenHeight    int = 144
	ScreenWidth     int = 160
)

type GameboyOptions struct {
	RomPath              string
	SerialOutputFunction func(byte)
	Headless             bool
}

type Gameboy struct {
	Display       *Display
	Cpu           *Cpu
	Memory        MemoryDevice
	Gpu           *Gpu
	Timer         *Timer
	Inputs        *Inputs
	Debugger      *Debugger
	WorkingScreen [ScreenWidth][ScreenHeight][3]uint8
	ReadyToRender [ScreenWidth][ScreenHeight][3]uint8
	Halted        bool
	Options       *GameboyOptions
}

func NewGameboy(options *GameboyOptions) (*Gameboy, error) {
	c := NewCpu()
	t := NewTimer()
	i := NewInputs()

	var d *Display
	if !options.Headless {
		d = NewDisplay()
	}

	gb := &Gameboy{
		Cpu:      c,
		Display:  d,
		Memory:   nil,
		Gpu:      nil,
		Timer:    t,
		Inputs:   i,
		Debugger: &Debugger{AddressEnabled: false},
		Halted:   false,
		Options:  options,
	}

	gb.Gpu = NewGpu(gb)
	var err error
	gb.Memory, err = NewMemory(gb)
	if err != nil {
		return nil, err
	}

	return gb, err
}

func (gb *Gameboy) Run() {
	cyclesPerFrame := int(float32(ClockSpeed) / float32(FramesPerSecond) * SpeedBoost)
	frameDuration := time.Second / time.Duration(FramesPerSecond)
	frameCount := 0
	lastFpsUpdate := time.Now()

	ticker := time.NewTicker(frameDuration)

	for ; true; <-ticker.C {
		frameCount++

		gb.Display.HandleInput(gb)
		gb.Inputs.HandleInput(gb)
		gb.Inputs.ClearButtonList()

		cycles := 0
		for cycles <= cyclesPerFrame {
			cyclesCPU := gb.Cpu.Tick(gb)
			cycles += cyclesCPU
			gb.Gpu.Update(gb, cyclesCPU)
		}

		if !gb.Options.Headless {
			gb.Display.Render(gb)
		}

		since := time.Since(lastFpsUpdate)
		if since > time.Second {
			lastFpsUpdate = time.Now()
			gb.Display.window.SetTitle(fmt.Sprintf("AXWGameboy (%d FPS)", frameCount))
			frameCount = 0
		}
	}
}
