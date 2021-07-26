package cui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type GpuViewWidget struct {
	cui  *GameboyCui
	name string
	x, y int
	w, h int
}

func NewGpuViewWidget(cui *GameboyCui, name string, x, y int) *GpuViewWidget {
	return &GpuViewWidget{cui: cui, name: name, x: x, y: y, w: 80, h: 50}
}

func (w *GpuViewWidget) Layout(g *gocui.Gui) error {
	_, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}
	return nil
}

func (w *GpuViewWidget) UpdateWidget(gb *internal.Gameboy) error {
	v, err := w.cui.gui.View(w.name)
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(
		v,
		"BGP: %v\nS1P: %v\nS2P: %v\n",
		gb.Gpu.BgPaletteMap,
		gb.Gpu.SpritePaletteMap[0],
		gb.Gpu.SpritePaletteMap[1],
	)
	for i, x := range gb.Gpu.SpriteObjectData {
		fmt.Fprintf(
			v,
			"S%d: %v\n",
			i,
			x,
		)
	}
	return nil
}
