package debugscreens

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type SpriteList struct {
	X int
	Y int
}

func (t *SpriteList) GetXPos() int   { return t.X }
func (t *SpriteList) GetYPos() int   { return t.Y }
func (t *SpriteList) GetWidth() int  { return 200 }
func (t *SpriteList) GetHeight() int { return 10 * 40 }

func (t *SpriteList) Draw(gb *internal.Gameboy, screen *ebiten.Image) {
	subscreen := screen.SubImage(image.Rect(t.GetXPos(), t.GetYPos(), t.GetXPos()+t.GetWidth(), t.GetYPos()+t.GetHeight())).(*ebiten.Image)
	for i := 0; i < len(gb.Gpu.SpriteObjectData); i++ {
		ebitenutil.DebugPrintAt(subscreen, gb.Gpu.SpriteObjectData[i].String(), t.X, t.Y+10*i)
	}
}
