package ebitenprovider

type MiscEvent int

const (
	ShutdownGame MiscEvent = iota
	SpeedboostToggle
	PauseToggle
)

func (a *AXWGameboyEbitenGame) handleMiscEvents(events []MiscEvent) {
	for _, event := range events {
		if event == SpeedboostToggle {
			a.toggleSpeedboost()
		} else if event == PauseToggle {
			a.togglePause()
		} else if event == ShutdownGame {
			a.markGameForShutdown()
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
