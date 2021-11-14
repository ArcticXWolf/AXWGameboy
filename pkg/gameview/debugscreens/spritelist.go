package debugscreens

import (
	"fmt"
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
func (t *SpriteList) GetWidth() int  { return 300 }
func (t *SpriteList) GetHeight() int { return 10 * 40 }

func (t *SpriteList) Draw(gb *internal.Gameboy, screen *ebiten.Image) {
	subscreen := screen.SubImage(image.Rect(t.GetXPos(), t.GetYPos(), t.GetXPos()+t.GetWidth(), t.GetYPos()+t.GetHeight())).(*ebiten.Image)
	for i := 0; i < len(gb.Gpu.SpriteObjectData); i++ {
		str := gb.Gpu.SpriteObjectData[i].String()
		oam0 := gb.Gpu.Oam[4*i]
		oam1 := gb.Gpu.Oam[4*i+1]
		oam2 := gb.Gpu.Oam[4*i+2]
		oam3 := gb.Gpu.Oam[4*i+3]
		str = fmt.Sprintf("%s B%02x %02x %02x %02x", str, oam0, oam1, oam2, oam3)
		ebitenutil.DebugPrintAt(subscreen, str, t.X, t.Y+10*i)
	}
}
