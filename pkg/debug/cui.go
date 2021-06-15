package debug

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type DisassemblyLogWidget struct {
	name string
	x, y int
	w, h int
}

func NewDisassemblyLogWidget(name string, x, y int) *DisassemblyLogWidget {
	return &DisassemblyLogWidget{name: name, x: x, y: y, w: 50, h: 30}
}

func (w *DisassemblyLogWidget) Layout(g *gocui.Gui) error {
	_, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}
	return nil
}

type GameboyCui struct {
	gui *gocui.Gui
}

func NewGui() *GameboyCui {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	d := NewDisassemblyLogWidget("Disassembly", 0, 0)
	g.SetManager(d)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	return &GameboyCui{
		gui: g,
	}
}

func (gc *GameboyCui) UpdateView(gb *internal.Gameboy) {
	gc.gui.UpdateAsync(func(g *gocui.Gui) error {
		v, err := g.View("Disassembly")
		if err != nil {
			log.Panicln(err)
		}
		v.Clear()
		fmt.Fprintf(
			v,
			"Clock: %010d\n\nA:  0x%02x\nB:  0x%02x\nC:  0x%02x\nD:  0x%02x\nE:  0x%02x\nH:  0x%02x\nL:  0x%02x\n\nPC: 0x%04x %s\nSP: 0x%04x\n\nIE: %v",
			gb.Cpu.ClockCycles,
			gb.Cpu.Registers.A,
			gb.Cpu.Registers.B,
			gb.Cpu.Registers.C,
			gb.Cpu.Registers.D,
			gb.Cpu.Registers.E,
			gb.Cpu.Registers.H,
			gb.Cpu.Registers.L,
			gb.Cpu.Registers.Pc,
			internal.Opcodes[gb.PeekPc(0)].Label,
			gb.Cpu.Registers.Sp,
			gb.Cpu.Registers.Ime,
		)
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
