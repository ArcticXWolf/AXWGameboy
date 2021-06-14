package internal

import (
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
	g := NewGpu()
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
		Gpu:      g,
		Timer:    t,
		Inputs:   i,
		Debugger: &Debugger{AddressEnabled: false},
		Halted:   false,
		Options:  options,
	}

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

	ticker := time.NewTicker(frameDuration)

	for ; true; <-ticker.C {
		frameCount++

		gb.Display.HandleInput(gb)
		gb.Inputs.HandleInput(gb)
		gb.Inputs.ClearButtonList()

		for i := 0; i < cyclesPerFrame; i++ {
			cycles := gb.Cpu.Tick(gb)
			gb.Gpu.Update(gb, cycles)
		}

		if !gb.Options.Headless {
			gb.Display.Render(gb)
		}
	}
}
