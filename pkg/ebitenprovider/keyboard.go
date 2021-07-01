package ebitenprovider

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"go.janniklasrichter.de/axwgameboy/internal"
)

var keyboardButtonMap = map[ebiten.Key]internal.Button{
	ebiten.KeyZ:     internal.ButtonA,
	ebiten.KeyX:     internal.ButtonB,
	ebiten.KeyAlt:   internal.ButtonSelect,
	ebiten.KeySpace: internal.ButtonStart,
	ebiten.KeyUp:    internal.ButtonUp,
	ebiten.KeyDown:  internal.ButtonDown,
	ebiten.KeyRight: internal.ButtonRight,
	ebiten.KeyLeft:  internal.ButtonLeft,
}

var keyboardMiscEventMap = map[ebiten.Key]MiscEvent{
	ebiten.KeyEscape:    ShutdownGame,
	ebiten.KeyShiftLeft: SpeedboostToggle,
	ebiten.KeyP:         PauseToggle,
}

func (a *AXWGameboyEbitenGame) handleKeyboardInputs() {
	for key, button := range keyboardButtonMap {
		if inpututil.IsKeyJustPressed(key) {
			a.Gameboy.Inputs.ButtonsPressed = append(a.Gameboy.Inputs.ButtonsPressed, button)
		}
		if inpututil.IsKeyJustReleased(key) {
			a.Gameboy.Inputs.ButtonsReleased = append(a.Gameboy.Inputs.ButtonsReleased, button)
		}
	}
}

func (a *AXWGameboyEbitenGame) handleKeyboardInputsForMiscEvents() (events []MiscEvent) {
	for key, event := range keyboardMiscEventMap {
		if inpututil.IsKeyJustPressed(key) {
			events = append(events, event)
		}
	}
	return
}
