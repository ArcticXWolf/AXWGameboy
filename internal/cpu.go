package internal

import (
	"fmt"
)

var ClockSpeed int = 4194304

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

func (r *Registers) FlagZ() bool {
	return r.Flags&0x80 != 0
}
func (r *Registers) FlagN() bool {
	return r.Flags&0x40 != 0
}
func (r *Registers) FlagH() bool {
	return r.Flags&0x20 != 0
}
func (r *Registers) FlagC() bool {
	return r.Flags&0x10 != 0
}

func (r *Registers) SetFlagZ(isSet bool) {
	if isSet {
		r.Flags |= uint8(0x80)
	} else {
		r.Flags &= ^uint8(0x80)
	}
}
func (r *Registers) SetFlagN(isSet bool) {
	if isSet {
		r.Flags |= uint8(0x40)
	} else {
		r.Flags &= ^uint8(0x40)
	}
}
func (r *Registers) SetFlagH(isSet bool) {
	if isSet {
		r.Flags |= uint8(0x20)
	} else {
		r.Flags &= ^uint8(0x20)
	}
}
func (r *Registers) SetFlagC(isSet bool) {
	if isSet {
		r.Flags |= uint8(0x10)
	} else {
		r.Flags &= ^uint8(0x10)
	}
}

func (r *Registers) String() string {
	return fmt.Sprintf("A(%02x) B(%02x) C(%02x) D(%02x) E(%02x) H(%02x) L(%02x) F(%08b) PC(%04x) SP(%04x) ", r.A, r.B, r.C, r.D, r.E, r.H, r.L, r.Flags, r.Pc, r.Sp)
}

type Cpu struct {
	Registers   *Registers
	ClockCycles int
	Memory      *Mmu
	Gpu         *Gpu
}

func NewCpu() *Cpu {
	fillUninplementedOpcodes()
	fillUninplementedOpcodesCb()
	g := &Gpu{
		CurrentScanline: 0x90,
	}
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
		Memory:      NewMemory(g),
		Gpu:         g,
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

func (c *Cpu) Tick(d *Debugger) {
	d.checkBreakpoint(c)
	d.breakIfNecessary(c)
	_, opcode := c.getNextOpcode()
	c.ClockCycles += opcode.Cycles
	opcode.Function(c)
}

func (c *Cpu) getNextOpcode() (uint8, *opcode) {
	code := c.popPc()
	return code, opcodes[code]
}

func (c *Cpu) peekPc(offset int) uint8 {
	return c.Memory.ReadByte(c.Registers.Pc + uint16(offset))
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
	return uint16(result2)<<8 | uint16(result1)
}

func (c *Cpu) String() string {
	step := fmt.Sprintf("(0x%04x) %02x, %15s", c.Registers.Pc, c.peekPc(0), opcodes[c.peekPc(0)].Label)
	peek := fmt.Sprintf("%02x %02x %02x", c.peekPc(1), c.peekPc(2), c.peekPc(3))
	return fmt.Sprintf("STEP: %s | PEEK: %s | REG: %s | CLOCK: %v", step, peek, c.Registers.String(), c.ClockCycles)
}
