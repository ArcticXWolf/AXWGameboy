package gameview

import (
	"log"
	"os"
	"runtime/pprof"
)

type MiscEvent int

const (
	ShutdownGame MiscEvent = iota
	SpeedboostToggle
	PauseToggle
	DebugToggle
	SoundChannel1Toggle
	SoundChannel2Toggle
	SoundChannel3Toggle
	SoundChannel4Toggle
	VolumeUp
	VolumeDown
	StartProfiling
	StopProfiling
)

func (a *AXWGameboyEbitenGameView) handleMiscEvents(events []MiscEvent) {
	for _, event := range events {
		if event == SpeedboostToggle {
			a.toggleSpeedboost()
		} else if event == PauseToggle {
			a.TogglePause()
		} else if event == DebugToggle {
			a.toggleDebugview()
		} else if event == ShutdownGame {
			a.markGameForShutdown()
		} else if event == SoundChannel1Toggle {
			a.Gameboy.Apu.ToggleSoundChannel(1)
		} else if event == SoundChannel2Toggle {
			a.Gameboy.Apu.ToggleSoundChannel(2)
		} else if event == SoundChannel3Toggle {
			a.Gameboy.Apu.ToggleSoundChannel(3)
		} else if event == SoundChannel4Toggle {
			a.Gameboy.Apu.ToggleSoundChannel(4)
		} else if event == VolumeUp {
			a.Gameboy.Apu.ChangeVolume(0.1)
		} else if event == VolumeDown {
			a.Gameboy.Apu.ChangeVolume(-0.1)
		} else if event == StartProfiling {
			a.startProfiling()
		} else if event == StopProfiling {
			a.stopProfiling()
		}
	}
}

func (a *AXWGameboyEbitenGameView) toggleSpeedboost() {
	a.isSpeedboostActive = !a.isSpeedboostActive

	if a.isSpeedboostActive {
		a.Gameboy.Cpu.SpeedBoost = 3.0
	} else {
		a.Gameboy.Cpu.SpeedBoost = 1.0
	}
}

func (a *AXWGameboyEbitenGameView) TogglePause() {
	a.isPaused = !a.isPaused
}

func (a *AXWGameboyEbitenGameView) toggleDebugview() {
	a.isDebugRequested = !a.isDebugRequested
}

func (a *AXWGameboyEbitenGameView) markGameForShutdown() {
	a.isTerminated = true
}

func (a *AXWGameboyEbitenGameView) startProfiling() {
	f, err := os.Create("cpu.profile")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
}

func (a *AXWGameboyEbitenGameView) stopProfiling() {
	pprof.StopCPUProfile()
}
