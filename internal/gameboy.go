package internal

import (
	"log"
	"time"

	"go.janniklasrichter.de/axwgameboy/pkg/apu"
)

const (
	FramesPerSecond int = 60
	ScreenHeight    int = 144
	ScreenWidth     int = 160
)

type GameboyOptions struct {
	SavePath             string
	RomPath              string
	Palette              string
	SoundVolume          float64
	SerialOutputFunction func(byte)
	DisplayProvider      DisplayProvider
	InputProvider        InputProvider
	OnCycleFunction      func(*Gameboy)
	OnFrameFunction      func(*Gameboy)
}

type Gameboy struct {
	InputProvider InputProvider
	Cpu           *Cpu
	Memory        *Mmu
	Gpu           *Gpu
	Apu           Apu
	Timer         *Timer
	Inputs        *Inputs
	Debugger      *Debugger
	WorkingScreen [ScreenWidth][ScreenHeight][3]uint8
	ReadyToRender [ScreenWidth][ScreenHeight][3]uint8
	Halted        bool
	Quit          bool
	Options       *GameboyOptions
	LastSave      time.Time
}

func NewGameboy(options *GameboyOptions) (*Gameboy, error) {
	c := NewCpu()
	i := NewInputs()
	a := &apu.APU{}
	a.Init(true, options.SoundVolume)

	gb := &Gameboy{
		Cpu:           c,
		InputProvider: options.InputProvider,
		Memory:        nil,
		Gpu:           nil,
		Apu:           a,
		Timer:         nil,
		Inputs:        i,
		Debugger:      &Debugger{AddressEnabled: false},
		Halted:        false,
		Options:       options,
		LastSave:      time.Now(),
	}

	gb.Gpu = NewGpu(gb)
	gb.Timer = NewTimer(gb)
	var err error
	gb.Memory, err = NewMemory(gb)
	if err != nil {
		return nil, err
	}

	return gb, err
}

// Deprecated
func (gb *Gameboy) Run() {
	frameDuration := time.Second / time.Duration(FramesPerSecond)
	frameCount := 0
	lastSave := time.Now()

	ticker := time.NewTicker(frameDuration)

	for ; !gb.Quit; <-ticker.C {
		cyclesPerFrame := int(float32(ClockSpeed) / float32(FramesPerSecond) * gb.Cpu.SpeedBoost)
		frameCount++

		if gb.InputProvider != nil {
			gb.InputProvider.HandleInput(gb)
		}
		gb.Inputs.HandleInput(gb)
		gb.Inputs.ClearButtonList()

		cycles := 0
		for cycles <= cyclesPerFrame {
			cyclesCPU := gb.Cpu.Tick(gb)
			cycles += cyclesCPU
			gb.Gpu.Update(gb, cyclesCPU)
			gb.Memory.Cartridge.UpdateComponentsPerCycle()

			if gb.Options.OnCycleFunction != nil {
				gb.Options.OnCycleFunction(gb)
			}
			gb.Apu.Buffer(cyclesCPU, 1)
		}

		if gb.Options.OnFrameFunction != nil {
			gb.Options.OnFrameFunction(gb)
		}
		// if gb.DisplayProvider != nil {
		// 	gb.DisplayProvider.Render(gb)
		// }

		if time.Since(lastSave) > time.Minute {
			if gb.Options.SavePath != "" {
				err := gb.Memory.Cartridge.SaveRam(gb.Options.SavePath)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	if gb.Options.SavePath != "" {
		err := gb.Memory.Cartridge.SaveRam(gb.Options.SavePath)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (gb *Gameboy) UpdateFrame(cyclesPerFrame int) {
	if gb.InputProvider != nil {
		gb.InputProvider.HandleInput(gb)
	}
	gb.Inputs.HandleInput(gb)
	gb.Inputs.ClearButtonList()

	cycles := 0
	for cycles <= cyclesPerFrame {
		cyclesCPU := gb.Cpu.Tick(gb)
		cycles += cyclesCPU
		gb.Gpu.Update(gb, cyclesCPU)
		gb.Memory.Cartridge.UpdateComponentsPerCycle()

		if gb.Options.OnCycleFunction != nil {
			gb.Options.OnCycleFunction(gb)
		}

		gb.Apu.Buffer(cyclesCPU, 1)
	}

	if gb.Options.OnFrameFunction != nil {
		gb.Options.OnFrameFunction(gb)
	}

	if time.Since(gb.LastSave) > time.Minute {
		if gb.Options.SavePath != "" {
			err := gb.Memory.Cartridge.SaveRam(gb.Options.SavePath)
			if err != nil {
				log.Println(err)
			}
			gb.LastSave = time.Now()
		}
	}
}
