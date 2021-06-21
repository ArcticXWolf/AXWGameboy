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
	var cycles int

	gb.Debugger.checkBreakpoint(gb)
	if gb.Debugger.Step {
		gb.Debugger.triggerBreakpoint(gb)
	}

	gb.Timer.Update(gb)

	var cyclesOp int
	if !gb.Halted {
		_, opcode := c.getNextOpcode(gb)
		cyclesOp += opcode.Cycles
		opcode.Function(gb)
	} else {
		cyclesOp += 4
	}
	c.ClockCycles += cyclesOp
	cycles += cyclesOp

	cyclesInt := c.handleInterrupts(gb)
	c.ClockCycles += cyclesInt
	cycles += cyclesInt

	return cycles
}

func (c *Cpu) getNextOpcode(gb *Gameboy) (uint8, *Opcode) {
	code := gb.popPc()
	return code, Opcodes[code]
}

func (c *Cpu) handleInterrupts(gb *Gameboy) int {
	if !c.Registers.Ime && !gb.Halted {
		return 0
	}

	enabled := gb.Memory.GetInterruptFlags().EnableFlags
	triggered := gb.Memory.GetInterruptFlags().TriggeredFlags

	for i := 0; i < 5; i++ {
		if (enabled>>i)&0x1 > 0 && (triggered>>i)&0x1 > 0 {
			cycles := 0
			if gb.Halted {
				cycles += 4
			}

			serviced := instructionInterrupt(gb, i)

			if serviced {
				cycles += 20
			}
			return cycles
		}
	}
	return 0
}

func (gb *Gameboy) PeekPc(offset int) uint8 {
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
	step := fmt.Sprintf("(0x%04x) %02x, %15s", gb.Cpu.Registers.Pc, gb.PeekPc(0), Opcodes[gb.PeekPc(0)].Label)
	peek := fmt.Sprintf("%02x %02x %02x", gb.PeekPc(1), gb.PeekPc(2), gb.PeekPc(3))
	// isr := fmt.Sprintf("%v E%02x T%02x", gb.Cpu.Registers.Ime, gb.Memory.GetInterruptFlags().EnableFlags, gb.Memory.GetInterruptFlags().TriggeredFlags)
	return fmt.Sprintf("%010d STEP: %s | PEEK: %s | REG: %s", gb.Cpu.ClockCycles, step, peek, gb.Cpu.Registers.String())
}
