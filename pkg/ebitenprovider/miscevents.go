package ebitenprovider

type MiscEvent int

const (
	ShutdownGame MiscEvent = iota
	SpeedboostToggle
	PauseToggle
	SoundChannel1Toggle
	SoundChannel2Toggle
	SoundChannel3Toggle
	SoundChannel4Toggle
)

func (a *AXWGameboyEbitenGame) handleMiscEvents(events []MiscEvent) {
	for _, event := range events {
		if event == SpeedboostToggle {
			a.toggleSpeedboost()
		} else if event == PauseToggle {
			a.togglePause()
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
		}
	}
}

func (a *AXWGameboyEbitenGame) toggleSpeedboost() {
	a.isSpeedboostActive = !a.isSpeedboostActive

	if a.isSpeedboostActive {
		a.Gameboy.Cpu.SpeedBoost = 3.0
	} else {
		a.Gameboy.Cpu.SpeedBoost = 1.0
	}
}

func (a *AXWGameboyEbitenGame) togglePause() {
	a.isPaused = !a.isPaused
}

func (a *AXWGameboyEbitenGame) markGameForShutdown() {
	a.isTerminated = true
}
