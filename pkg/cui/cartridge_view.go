package cui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type CartridgeViewWidget struct {
	cui  *GameboyCui
	name string
	x, y int
	w, h int
}

func NewCartridgeViewWidget(cui *GameboyCui, name string, x, y int) *CartridgeViewWidget {
	return &CartridgeViewWidget{cui: cui, name: name, x: x, y: y, w: 80, h: 2}
}

func (w *CartridgeViewWidget) Layout(g *gocui.Gui) error {
	_, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}
	return nil
}

func (w *CartridgeViewWidget) UpdateWidget(gb *internal.Gameboy) error {
	v, err := w.cui.gui.View(w.name)
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(
		v,
		"Cart: %v",
		gb.Memory.Cartridge.String(),
	)
	return nil
}
