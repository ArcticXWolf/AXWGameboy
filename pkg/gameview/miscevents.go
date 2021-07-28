package gameview

type MiscEvent int

const (
	ShutdownGame MiscEvent = iota
	SpeedboostToggle
	PauseToggle
	DebugToggle
	Tilemap0Toggle
	Tilemap1Toggle
	SoundChannel1Toggle
	SoundChannel2Toggle
	SoundChannel3Toggle
	SoundChannel4Toggle
)

func (a *AXWGameboyEbitenGameView) handleMiscEvents(events []MiscEvent) {
	for _, event := range events {
		if event == SpeedboostToggle {
			a.toggleSpeedboost()
		} else if event == PauseToggle {
			a.togglePause()
		} else if event == DebugToggle {
			a.Gameboy.Debugger.TriggerBreakpoint(a.Gameboy)
		} else if event == Tilemap0Toggle {
			a.toggleTilemap(0)
		} else if event == Tilemap1Toggle {
			a.toggleTilemap(1)
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

func (a *AXWGameboyEbitenGameView) toggleSpeedboost() {
	a.isSpeedboostActive = !a.isSpeedboostActive

	if a.isSpeedboostActive {
		a.Gameboy.Cpu.SpeedBoost = 3.0
	} else {
		a.Gameboy.Cpu.SpeedBoost = 1.0
	}
}

func (a *AXWGameboyEbitenGameView) togglePause() {
	a.isPaused = !a.isPaused
}

func (a *AXWGameboyEbitenGameView) toggleTilemap(number int) {
	if a.tilemapVram != number && a.isTilemapEnabled {
		a.tilemapVram = number
	} else {
		a.isTilemapRequested = !a.isTilemapRequested
	}
}

func (a *AXWGameboyEbitenGameView) markGameForShutdown() {
	a.isTerminated = true
}
