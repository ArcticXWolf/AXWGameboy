package cui

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type CpuViewWidget struct {
	cui  *GameboyCui
	name string
	x, y int
	w, h int
}

func NewCpuViewWidget(cui *GameboyCui, name string, x, y int) *CpuViewWidget {
	return &CpuViewWidget{cui: cui, name: name, x: x, y: y, w: 25, h: 15}
}

func (w *CpuViewWidget) Layout(g *gocui.Gui) error {
	_, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h, 0)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
	}
	return nil
}

func (w *CpuViewWidget) UpdateWidget(gb *internal.Gameboy) error {
	v, err := w.cui.gui.View(w.name)
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintf(
		v,
		"Clock: %011d\n\nA:  0x%02x\nB:  0x%02x\nC:  0x%02x\nD:  0x%02x\nE:  0x%02x\nH:  0x%02x\nL:  0x%02x\n\nPC: 0x%04x %s\nSP: 0x%04x\n\nIE: %v",
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
}
