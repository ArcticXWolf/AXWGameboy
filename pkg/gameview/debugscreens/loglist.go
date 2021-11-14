package debugscreens

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type LogList struct {
	X int
	Y int
}

func (t *LogList) GetXPos() int   { return t.X }
func (t *LogList) GetYPos() int   { return t.Y }
func (t *LogList) GetWidth() int  { return 300 }
func (t *LogList) GetHeight() int { return 10 * 60 }

func (t *LogList) Draw(gb *internal.Gameboy, screen *ebiten.Image) {
	subscreen := screen.SubImage(image.Rect(t.GetXPos(), t.GetYPos(), t.GetXPos()+t.GetWidth(), t.GetYPos()+t.GetHeight())).(*ebiten.Image)
	messages := gb.RingLogger.GetAllMessages()
	var i int
	for k, v := range messages {
		if k == "gpu" {
			for _, line := range v {
				ebitenutil.DebugPrintAt(subscreen, fmt.Sprintf("%-6s %s", k, line), t.X, t.Y+10*i)
				i++
			}
		}
	}
}
