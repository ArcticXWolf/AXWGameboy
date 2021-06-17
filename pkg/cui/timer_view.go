package cui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type TimerViewWidget struct {
	cui  *GameboyCui
	name string
	x, y int
	w, h int
}

func NewTimerViewWidget(cui *GameboyCui, name string, x, y int) *TimerViewWidget {
	return &TimerViewWidget{cui: cui, name: name, x: x, y: y, w: 15, h: 8}
}

func (w *TimerViewWidget) Layout(g *gocui.Gui) error {
	_, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}
	return nil
}

func (w *TimerViewWidget) UpdateWidget(gb *internal.Gameboy) error {
	v, err := w.cui.gui.View(w.name)
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(
		v,
		"SYS:  0x%04x\nDIV: 0x%02x\n\nTACE: %v\nTMA:  %04d\nTIMA: %04d\nTACS: %d",
		gb.Timer.DividerValue,
		(gb.Timer.DividerValue >> 8),
		gb.Timer.ControlFlag&0x4 > 0,
		gb.Timer.ModuloValue,
		gb.Timer.CounterValue,
		gb.Timer.ControlFlag&0x03,
	)
	return nil
}
