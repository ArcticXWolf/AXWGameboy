package internal

import (
	"io"
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
	SaveWriter           io.Writer
	SaveReader           io.Reader
	RomReader            io.Reader
	Palette              string
	SoundEnabled         bool
	SoundVolume          float64
	OSBEnabled           bool
	CGBEnabled           bool
	SerialOutputFunction func(byte)
	DisplayProvider      DisplayProvider
	InputProvider        InputProvider
	OnCycleFunction      func(*Gameboy)
	OnFrameFunction      func(*Gameboy)
}

type Gameboy struct {
	InputProvider        InputProvider
	Cpu                  *Cpu
	Memory               *Mmu
	Gpu                  *Gpu
	Apu                  Apu
	Timer                *Timer
	Inputs               *Inputs
	Debugger             *Debugger
	CheatCodeManager     *CheatCodeManager
	WorkingScreen        [ScreenWidth][ScreenHeight][3]uint8
	ReadyToRender        [ScreenWidth][ScreenHeight][3]uint8
	CgbModeEnabled       bool
	doubleSpeed          bool
	doubleSpeedRequested bool
	Halted               bool
	Quit                 bool
	RingLogger           *RingLogger
	Options              *GameboyOptions
	LastSave             time.Time
}

func NewGameboy(options *GameboyOptions) (*Gameboy, error) {
	c := NewCpu()
	i := NewInputs()

	gb := &Gameboy{
		Cpu:              c,
		InputProvider:    options.InputProvider,
		Memory:           nil,
		Gpu:              nil,
		Apu:              nil,
		Timer:            nil,
		Inputs:           i,
		Debugger:         &Debugger{AddressEnabled: false},
		CheatCodeManager: &CheatCodeManager{genieEnabled: true, sharkEnabled: true},
		RingLogger:       NewRingLogger(60),
		Halted:           false,
		Options:          options,
		LastSave:         time.Now(),
	}

	a := apu.NewApu()
	a.Init(options.SoundEnabled, options.SoundVolume)
	gb.Apu = a

	gb.Gpu = NewGpu(gb)
	gb.Timer = NewTimer(gb)
	var err error
	gb.Memory, _, err = NewMemory(gb)
	if err != nil {
		return nil, err
	}

	gb.CgbModeEnabled = options.CGBEnabled

	return gb, err
}

func (gb *Gameboy) Update(cyclesPerFrame int) {
	if gb.InputProvider != nil {
		gb.InputProvider.HandleInput(gb)
	}
	gb.Inputs.HandleInput(gb)
	gb.Inputs.ClearButtonList()

	gb.Apu.ToggleSound(gb.Options.SoundEnabled)

	cycles := 0
	for cycles <= cyclesPerFrame {
		cyclesCPU := gb.Cpu.Tick(gb)
		cycles += cyclesCPU
		gb.Gpu.Update(gb, cyclesCPU)
		gb.Memory.Cartridge.UpdateComponentsPerCycle(uint16(cyclesCPU))

		if gb.Options.OnCycleFunction != nil {
			gb.Options.OnCycleFunction(gb)
		}

		if gb.Apu.IsSoundEnabled() {
			gb.Apu.Buffer(cyclesCPU, int(gb.GetSpeedMultiplier()), float64(gb.Cpu.SpeedBoost))
		}
	}

	if gb.Options.OnFrameFunction != nil {
		gb.Options.OnFrameFunction(gb)
	}

	if time.Since(gb.LastSave) > time.Second {
		if gb.Options.SaveWriter != nil {
			err := gb.Memory.Cartridge.SaveRam(gb.Options.SaveWriter)
			if err != nil {
				log.Println(err)
			}
		}
		gb.LastSave = time.Now()
	}
}
