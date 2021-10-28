package debugscreens

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type DebugScreen interface {
	GetXPos() int
	GetYPos() int
	GetWidth() int
	GetHeight() int
	Draw(gb *internal.Gameboy, image *ebiten.Image)
}
