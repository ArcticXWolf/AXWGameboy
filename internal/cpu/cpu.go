package cpu

import (
	"fmt"

	"go.janniklasrichter.de/axwgameboy/internal/memory"
)

type Registers struct {
	A     uint8
	B     uint8
	C     uint8
	D     uint8
	E     uint8
	H     uint8
	L     uint8
	Flags uint8

	Pc uint16
	Sp uint16
}

func (r *Registers) String() string {
	return fmt.Sprintf("Registers: A(%02x) B(%02x) C(%02x) D(%02x) E(%02x) H(%02x) L(%02x) F(%b) PC(%04x) SP(%04x) ", r.A, r.B, r.C, r.D, r.E, r.H, r.L, r.Flags, r.Pc, r.Sp)
}

type Cpu struct {
	Registers   *Registers
	ClockCycles int
	Memory      *memory.Mmu
}

func New() *Cpu {
	fillUninplementedOpcodes()
	return &Cpu{
		Registers: &Registers{
			A:     0,
			B:     0,
			C:     0,
			D:     0,
			E:     0,
			H:     0,
			L:     0,
			Flags: 0,
			Pc:    0,
			Sp:    0,
		},
		ClockCycles: 0,
		Memory:      memory.New(),
	}
}

func (c *Cpu) Reset() {
	c.Registers.A = 0
	c.Registers.B = 0
	c.Registers.C = 0
	c.Registers.D = 0
	c.Registers.E = 0
	c.Registers.H = 0
	c.Registers.L = 0
	c.Registers.Flags = 0
	c.Registers.Pc = 0
	c.Registers.Sp = 0
	c.ClockCycles = 0
}

func (c *Cpu) Tick() {
	opcode := c.getNextOpcode()
	opcode.Function(c)
	c.Registers.Pc++
	c.Registers.Pc &= 0xFFFF
	c.ClockCycles += opcode.Cycles
}

func (c *Cpu) getNextOpcode() *opcode {
	return opcodes[c.Memory.ReadByte(c.Registers.Pc)]
}

func (c *Cpu) String() string {
	return fmt.Sprintf("Current State:\n%s\nClock: %v\n%s\n", c.Registers.String(), c.ClockCycles, c.Memory.String())
}
