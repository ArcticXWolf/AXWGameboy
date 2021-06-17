package cui

import (
	"errors"
	"log"

	"github.com/awesome-gocui/gocui"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type GameboyCui struct {
	gui     *gocui.Gui
	widgets []Widget
}

type Widget interface {
	gocui.Manager
	UpdateWidget(gb *internal.Gameboy) error
}

func NewGui() *GameboyCui {
	var err error
	var managers []gocui.Manager
	cui := &GameboyCui{
		gui:     nil,
		widgets: make([]Widget, 0),
	}

	cui.gui, err = gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer cui.gui.Close()

	w1 := NewCpuViewWidget(cui, "CPU View", 0, 0)
	cui.widgets = append(cui.widgets, w1)
	managers = append(managers, w1)

	w2 := NewTimerViewWidget(cui, "Timer View", 50, 0)
	cui.widgets = append(cui.widgets, w2)
	managers = append(managers, w2)

	cui.gui.SetManager(managers...)

	if err := cui.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	return cui
}

func (gc *GameboyCui) UpdateView(gb *internal.Gameboy) {
	gc.gui.UpdateAsync(func(g *gocui.Gui) error {
		for _, v := range gc.widgets {
			err := v.UpdateWidget(gb)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (gc *GameboyCui) RunLoop() {
	if err := gc.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
