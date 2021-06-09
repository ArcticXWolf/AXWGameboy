package cpu

import (
	"fmt"
	"log"

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
	fillUninplementedOpcodesCb()
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
	code, opcode := c.getNextOpcode()
	log.Printf("Got next opcode: Code: %02x, %s", code, opcode.Label)
	c.ClockCycles += opcode.Cycles
	opcode.Function(c)
}

func (c *Cpu) getNextOpcode() (uint8, *opcode) {
	code := c.popPc()
	return code, opcodes[code]
}

func (c *Cpu) popPc() uint8 {
	result := c.Memory.ReadByte(c.Registers.Pc)
	c.Registers.Pc++
	c.Registers.Pc &= 0xFFFF
	return result
}

func (c *Cpu) popPc16() uint16 {
	result1 := c.popPc()
	result2 := c.popPc()
	return uint16(result2<<8) | uint16(result1)
}

func (c *Cpu) String() string {
	return fmt.Sprintf("Current State:\n%s\nClock: %v\n", c.Registers.String(), c.ClockCycles)
}
