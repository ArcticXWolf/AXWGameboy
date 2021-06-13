package internal

import (
	"fmt"
)

var (
	ClockSpeed int     = 4194304
	SpeedBoost float32 = 1.0
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

	Pc  uint16
	Sp  uint16
	Ime bool
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
}

func NewCpu() *Cpu {
	fillUninplementedOpcodes()
	fillOpcodesCb(true)
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
			Ime:   false,
		},
		ClockCycles: 0,
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
	c.Registers.Ime = true
	c.ClockCycles = 0
}

func (c *Cpu) Tick(gb *Gameboy) int {
	gb.Debugger.checkBreakpoint(gb)
	if gb.Debugger.Step {
		gb.Debugger.triggerBreakpoint(gb)
	}
	_, opcode := c.getNextOpcode(gb)
	c.ClockCycles += opcode.Cycles
	opcode.Function(gb)

	cyclesInterrupt := c.handleInterrupts(gb)

	return opcode.Cycles + cyclesInterrupt
}

func (c *Cpu) getNextOpcode(gb *Gameboy) (uint8, *opcode) {
	code := gb.popPc()
	return code, opcodes[code]
}

func (c *Cpu) handleInterrupts(gb *Gameboy) int {
	cycles := 0
	if c.Registers.Ime {
		triggeredInterrupts := gb.Memory.GetInterruptFlags().EnableFlags & gb.Memory.GetInterruptFlags().TriggeredFlags

		if triggeredInterrupts&0x01 != 0 {
			gb.Memory.GetInterruptFlags().TriggeredFlags &= ^uint8(0x01)
			instructionInterrupt(gb, 0x0040)
			cycles += 20
		} else if triggeredInterrupts&0x02 != 0 {
			gb.Memory.GetInterruptFlags().TriggeredFlags &= ^uint8(0x02)
			instructionInterrupt(gb, 0x0048)
			cycles += 20
		} else if triggeredInterrupts&0x04 != 0 {
			gb.Memory.GetInterruptFlags().TriggeredFlags &= ^uint8(0x04)
			instructionInterrupt(gb, 0x0050)
			cycles += 20
		} else if triggeredInterrupts&0x08 != 0 {
			gb.Memory.GetInterruptFlags().TriggeredFlags &= ^uint8(0x08)
			instructionInterrupt(gb, 0x0058)
			cycles += 20
		} else if triggeredInterrupts&0x10 != 0 {
			gb.Memory.GetInterruptFlags().TriggeredFlags &= ^uint8(0x10)
			instructionInterrupt(gb, 0x0060)
			cycles += 20
		}
	}
	return cycles
}

func (gb *Gameboy) peekPc(offset int) uint8 {
	return gb.Memory.ReadByte(gb.Cpu.Registers.Pc + uint16(offset))
}

func (gb *Gameboy) popPc() uint8 {
	result := gb.Memory.ReadByte(gb.Cpu.Registers.Pc)
	gb.Cpu.Registers.Pc++
	gb.Cpu.Registers.Pc &= 0xFFFF
	return result
}

func (gb *Gameboy) popPc16() uint16 {
	result1 := gb.popPc()
	result2 := gb.popPc()
	return uint16(result2)<<8 | uint16(result1)
}

func (gb *Gameboy) String() string {
	step := fmt.Sprintf("(0x%04x) %02x, %15s", gb.Cpu.Registers.Pc, gb.peekPc(0), opcodes[gb.peekPc(0)].Label)
	peek := fmt.Sprintf("%02x %02x %02x", gb.peekPc(1), gb.peekPc(2), gb.peekPc(3))
	isr := fmt.Sprintf("%02x %02x", gb.Memory.GetInterruptFlags().EnableFlags, gb.Memory.GetInterruptFlags().TriggeredFlags)
	return fmt.Sprintf("STEP: %s | PEEK: %s | REG: %s | ISR: %s | CLOCK: %v", step, peek, gb.Cpu.Registers.String(), isr, gb.Cpu.ClockCycles)
}
