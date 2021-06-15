package internal

import (
	"fmt"
	"log"
)

type Opcode struct {
	Label    string
	Cycles   int
	Function func(*Gameboy)
}

var Opcodes = [0x100]*Opcode{
	0x00: {"NOP", 4, func(gb *Gameboy) {}},
	// 8-Bit Loads
	// LD r1,n
	0x06: {"LD B, n", 8, func(gb *Gameboy) { gb.Cpu.Registers.B = gb.popPc() }},
	0x0e: {"LD C, n", 8, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.popPc() }},
	0x16: {"LD D, n", 8, func(gb *Gameboy) { gb.Cpu.Registers.D = gb.popPc() }},
	0x1e: {"LD E, n", 8, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.popPc() }},
	0x26: {"LD H, n", 8, func(gb *Gameboy) { gb.Cpu.Registers.H = gb.popPc() }},
	0x2e: {"LD L, n", 8, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.popPc() }},
	0x36: {"LD (HL), n", 12, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.popPc())
	}},
	0x3e: {"LD A, n", 8, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.popPc() }},

	// LD r1,r2
	0x40: {"LD B, B", 4, func(gb *Gameboy) { gb.Cpu.Registers.B = gb.Cpu.Registers.B }},
	0x41: {"LD B, C", 4, func(gb *Gameboy) { gb.Cpu.Registers.B = gb.Cpu.Registers.C }},
	0x42: {"LD B, D", 4, func(gb *Gameboy) { gb.Cpu.Registers.B = gb.Cpu.Registers.D }},
	0x43: {"LD B, E", 4, func(gb *Gameboy) { gb.Cpu.Registers.B = gb.Cpu.Registers.E }},
	0x44: {"LD B, H", 4, func(gb *Gameboy) { gb.Cpu.Registers.B = gb.Cpu.Registers.H }},
	0x45: {"LD B, L", 4, func(gb *Gameboy) { gb.Cpu.Registers.B = gb.Cpu.Registers.L }},
	0x46: {"LD B, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
	}},

	0x48: {"LD C, B", 4, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.Cpu.Registers.B }},
	0x49: {"LD C, C", 4, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.Cpu.Registers.C }},
	0x4a: {"LD C, D", 4, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.Cpu.Registers.D }},
	0x4b: {"LD C, E", 4, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.Cpu.Registers.E }},
	0x4c: {"LD C, H", 4, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.Cpu.Registers.H }},
	0x4d: {"LD C, L", 4, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.Cpu.Registers.L }},
	0x4e: {"LD C, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
	}},

	0x50: {"LD D, B", 4, func(gb *Gameboy) { gb.Cpu.Registers.D = gb.Cpu.Registers.B }},
	0x51: {"LD D, C", 4, func(gb *Gameboy) { gb.Cpu.Registers.D = gb.Cpu.Registers.C }},
	0x52: {"LD D, D", 4, func(gb *Gameboy) { gb.Cpu.Registers.D = gb.Cpu.Registers.D }},
	0x53: {"LD D, E", 4, func(gb *Gameboy) { gb.Cpu.Registers.D = gb.Cpu.Registers.E }},
	0x54: {"LD D, H", 4, func(gb *Gameboy) { gb.Cpu.Registers.D = gb.Cpu.Registers.H }},
	0x55: {"LD D, L", 4, func(gb *Gameboy) { gb.Cpu.Registers.D = gb.Cpu.Registers.L }},
	0x56: {"LD D, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
	}},

	0x58: {"LD E, B", 4, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.Cpu.Registers.B }},
	0x59: {"LD E, C", 4, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.Cpu.Registers.C }},
	0x5a: {"LD E, D", 4, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.Cpu.Registers.D }},
	0x5b: {"LD E, E", 4, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.Cpu.Registers.E }},
	0x5c: {"LD E, H", 4, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.Cpu.Registers.H }},
	0x5d: {"LD E, L", 4, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.Cpu.Registers.L }},
	0x5e: {"LD E, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
	}},

	0x60: {"LD H, B", 4, func(gb *Gameboy) { gb.Cpu.Registers.H = gb.Cpu.Registers.B }},
	0x61: {"LD H, C", 4, func(gb *Gameboy) { gb.Cpu.Registers.H = gb.Cpu.Registers.C }},
	0x62: {"LD H, D", 4, func(gb *Gameboy) { gb.Cpu.Registers.H = gb.Cpu.Registers.D }},
	0x63: {"LD H, E", 4, func(gb *Gameboy) { gb.Cpu.Registers.H = gb.Cpu.Registers.E }},
	0x64: {"LD H, H", 4, func(gb *Gameboy) { gb.Cpu.Registers.H = gb.Cpu.Registers.H }},
	0x65: {"LD H, L", 4, func(gb *Gameboy) { gb.Cpu.Registers.H = gb.Cpu.Registers.L }},
	0x66: {"LD H, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
	}},

	0x68: {"LD L, B", 4, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.Cpu.Registers.B }},
	0x69: {"LD L, C", 4, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.Cpu.Registers.C }},
	0x6a: {"LD L, D", 4, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.Cpu.Registers.D }},
	0x6b: {"LD L, E", 4, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.Cpu.Registers.E }},
	0x6c: {"LD L, H", 4, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.Cpu.Registers.H }},
	0x6d: {"LD L, L", 4, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.Cpu.Registers.L }},
	0x6e: {"LD L, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
	}},

	0x70: {"LD (HL), B", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.B)
	}},
	0x71: {"LD (HL), C", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.C)
	}},
	0x72: {"LD (HL), D", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.D)
	}},
	0x73: {"LD (HL), E", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.E)
	}},
	0x74: {"LD (HL), H", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.H)
	}},
	0x75: {"LD (HL), L", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.L)
	}},

	0x78: {"LD A, B", 4, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.Cpu.Registers.B }},
	0x79: {"LD A, C", 4, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.Cpu.Registers.C }},
	0x7a: {"LD A, D", 4, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.Cpu.Registers.D }},
	0x7b: {"LD A, E", 4, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.Cpu.Registers.E }},
	0x7c: {"LD A, H", 4, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.Cpu.Registers.H }},
	0x7d: {"LD A, L", 4, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.Cpu.Registers.L }},
	0x7e: {"LD A, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
	}},
	0x7f: {"LD A, A", 4, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.Cpu.Registers.A }},
	0x0a: {"LD A, (BC)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.B)<<8 + uint16(gb.Cpu.Registers.C))
	}},
	0x1a: {"LD A, (DE)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.D)<<8 + uint16(gb.Cpu.Registers.E))
	}},
	0xfa: {"LD A, (nn)", 16, func(gb *Gameboy) { gb.Cpu.Registers.A = gb.Memory.ReadByte(gb.popPc16()) }},

	0x47: {"LD B, A", 4, func(gb *Gameboy) { gb.Cpu.Registers.B = gb.Cpu.Registers.A }},
	0x4f: {"LD C, A", 4, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.Cpu.Registers.A }},
	0x57: {"LD D, A", 4, func(gb *Gameboy) { gb.Cpu.Registers.D = gb.Cpu.Registers.A }},
	0x5f: {"LD E, A", 4, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.Cpu.Registers.A }},
	0x67: {"LD H, A", 4, func(gb *Gameboy) { gb.Cpu.Registers.H = gb.Cpu.Registers.A }},
	0x6f: {"LD L, A", 4, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.Cpu.Registers.A }},
	0x02: {"LD (BC), A", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.B)<<8+uint16(gb.Cpu.Registers.C), gb.Cpu.Registers.A)
	}},
	0x12: {"LD (DE), A", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.D)<<8+uint16(gb.Cpu.Registers.E), gb.Cpu.Registers.A)
	}},
	0x77: {"LD (HL), A", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.A)
	}},
	0xEA: {"LD (HL), A", 16, func(gb *Gameboy) { gb.Memory.WriteByte(gb.popPc16(), gb.Cpu.Registers.A) }},

	// LD A and ($FF00+C)
	0xF2: {"LD A, ($FF00+C)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = gb.Memory.ReadByte(uint16(0xFF00) + uint16(gb.Cpu.Registers.C))
	}},
	0xE2: {"LD ($FF00+C), A", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(0xFF00)+uint16(gb.Cpu.Registers.C), gb.Cpu.Registers.A)
	}},

	// LD A and HL-
	0x3a: {"LD A, (HL-)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
		gb.Cpu.Registers.L = (gb.Cpu.Registers.L - 1) & 0xFF
		if gb.Cpu.Registers.L == 0xFF {
			gb.Cpu.Registers.H = (gb.Cpu.Registers.H - 1) & 0xFF
		}
	}},
	0x32: {"LD (HL-), A", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.A)
		gb.Cpu.Registers.L = (gb.Cpu.Registers.L - 1) & 0xFF
		if gb.Cpu.Registers.L == 0xFF {
			gb.Cpu.Registers.H = (gb.Cpu.Registers.H - 1) & 0xFF
		}
	}},

	// LD A and HL+
	0x2a: {"LD A, (HL+)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
		gb.Cpu.Registers.L = (gb.Cpu.Registers.L + 1) & 0xFF
		if gb.Cpu.Registers.L == 0x0 {
			gb.Cpu.Registers.H = (gb.Cpu.Registers.H + 1) & 0xFF
		}
	}},
	0x22: {"LD (HL+), A", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.A)
		gb.Cpu.Registers.L = (gb.Cpu.Registers.L + 1) & 0xFF
		if gb.Cpu.Registers.L == 0x0 {
			gb.Cpu.Registers.H = (gb.Cpu.Registers.H + 1) & 0xFF
		}
	}},

	// LD A and ($FF00+n)
	0xF0: {"LD A, ($FF00+n)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = gb.Memory.ReadByte(uint16(0xFF00) + uint16(gb.popPc()))
	}},
	0xE0: {"LD ($FF00+n), A", 8, func(gb *Gameboy) {
		gb.Memory.WriteByte(uint16(0xFF00)+uint16(gb.popPc()), gb.Cpu.Registers.A)
	}},

	// 16-Bit Loads
	0x01: {"LD BC, nn", 12, func(gb *Gameboy) { gb.Cpu.Registers.C = gb.popPc(); gb.Cpu.Registers.B = gb.popPc() }},
	0x11: {"LD DE, nn", 12, func(gb *Gameboy) { gb.Cpu.Registers.E = gb.popPc(); gb.Cpu.Registers.D = gb.popPc() }},
	0x21: {"LD HL, nn", 12, func(gb *Gameboy) { gb.Cpu.Registers.L = gb.popPc(); gb.Cpu.Registers.H = gb.popPc() }},
	0x31: {"LD SP, nn", 12, func(gb *Gameboy) { gb.Cpu.Registers.Sp = gb.popPc16() }},
	0xF8: {"LD HL, SP+n", 12, func(gb *Gameboy) {
		n := int8(gb.popPc())
		sum := uint16(int32(gb.Cpu.Registers.Sp) + int32(n))
		gb.Cpu.Registers.L = uint8(sum) & 0xFF
		gb.Cpu.Registers.H = uint8((sum & 0xFF00) >> 8)
		sumTmp := gb.Cpu.Registers.Sp ^ uint16(n) ^ sum
		gb.Cpu.Registers.SetFlagZ(false)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH((sumTmp & 0x10) == 0x10)
		gb.Cpu.Registers.SetFlagC((sumTmp & 0x100) == 0x100)
	}},
	0xF9: {"LD SP, HL", 8, func(gb *Gameboy) { gb.Cpu.Registers.Sp = uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L) }},
	0x08: {"LD (nn), SP", 20, func(gb *Gameboy) { address := gb.popPc16(); gb.Memory.WriteWord(address, gb.Cpu.Registers.Sp) }},

	0xc5: {"PUSH BC", 16, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp--
		gb.Memory.WriteByte(gb.Cpu.Registers.Sp, gb.Cpu.Registers.B)
		gb.Cpu.Registers.Sp--
		gb.Memory.WriteByte(gb.Cpu.Registers.Sp, gb.Cpu.Registers.C)
	}},
	0xd5: {"PUSH DE", 16, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp--
		gb.Memory.WriteByte(gb.Cpu.Registers.Sp, gb.Cpu.Registers.D)
		gb.Cpu.Registers.Sp--
		gb.Memory.WriteByte(gb.Cpu.Registers.Sp, gb.Cpu.Registers.E)
	}},
	0xe5: {"PUSH HL", 16, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp--
		gb.Memory.WriteByte(gb.Cpu.Registers.Sp, gb.Cpu.Registers.H)
		gb.Cpu.Registers.Sp--
		gb.Memory.WriteByte(gb.Cpu.Registers.Sp, gb.Cpu.Registers.L)
	}},
	0xf5: {"PUSH AF", 16, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp--
		gb.Memory.WriteByte(gb.Cpu.Registers.Sp, gb.Cpu.Registers.A)
		gb.Cpu.Registers.Sp--
		gb.Memory.WriteByte(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Flags)
	}},

	0xc1: {"POP BC", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.C = gb.Memory.ReadByte(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp++
		gb.Cpu.Registers.B = gb.Memory.ReadByte(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp++
	}},
	0xd1: {"POP DE", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.E = gb.Memory.ReadByte(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp++
		gb.Cpu.Registers.D = gb.Memory.ReadByte(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp++
	}},
	0xe1: {"POP HL", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.L = gb.Memory.ReadByte(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp++
		gb.Cpu.Registers.H = gb.Memory.ReadByte(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp++
	}},
	0xf1: {"POP AF", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Flags = gb.Memory.ReadByte(gb.Cpu.Registers.Sp) & 0xF0
		gb.Cpu.Registers.Sp++
		gb.Cpu.Registers.A = gb.Memory.ReadByte(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp++
	}},

	// AND
	0xe6: {"AND n", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.popPc()
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	0xa0: {"AND B", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.Cpu.Registers.B
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	0xa1: {"AND C", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.Cpu.Registers.C
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	0xa2: {"AND D", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.Cpu.Registers.D
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	0xa3: {"AND E", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.Cpu.Registers.E
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	0xa4: {"AND H", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.Cpu.Registers.H
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	0xa5: {"AND L", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.Cpu.Registers.L
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	0xa6: {"AND (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	0xa7: {"AND A", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A &= gb.Cpu.Registers.A
		gb.Cpu.Registers.A &= 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(true)
		gb.Cpu.Registers.SetFlagC(false)
	}},
	// OR
	0xf6: {"OR n", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.popPc()
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xb0: {"OR B", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.Cpu.Registers.B
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xb1: {"OR C", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.Cpu.Registers.C
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xb2: {"OR D", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.Cpu.Registers.D
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xb3: {"OR E", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.Cpu.Registers.E
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xb4: {"OR H", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.Cpu.Registers.H
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xb5: {"OR L", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.Cpu.Registers.L
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xb6: {"OR (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xb7: {"OR A", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A |= gb.Cpu.Registers.A
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	// XOR
	0xee: {"XOR n", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.popPc()
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xaf: {"XOR A", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.Cpu.Registers.A
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xa8: {"XOR B", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.Cpu.Registers.B
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xa9: {"XOR C", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.Cpu.Registers.C
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xaa: {"XOR D", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.Cpu.Registers.D
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xab: {"XOR E", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.Cpu.Registers.E
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xac: {"XOR H", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.Cpu.Registers.H
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xad: {"XOR L", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.Cpu.Registers.L
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0xae: {"XOR (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A ^= gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
		gb.Cpu.Registers.A &= 0xFF
		if gb.Cpu.Registers.A == 0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},

	// INC r
	0x04: {"INC B", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionIncrement(gb, gb.Cpu.Registers.B)
	}},
	0x0C: {"INC C", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionIncrement(gb, gb.Cpu.Registers.C)
	}},
	0x14: {"INC D", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionIncrement(gb, gb.Cpu.Registers.D)
	}},
	0x1C: {"INC E", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionIncrement(gb, gb.Cpu.Registers.E)
	}},
	0x24: {"INC H", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionIncrement(gb, gb.Cpu.Registers.H)
	}},
	0x2C: {"INC L", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionIncrement(gb, gb.Cpu.Registers.L)
	}},
	0x34: {"INC (HL)", 4, func(gb *Gameboy) {
		incValue := gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), instructionIncrement(gb, incValue))
	}},
	0x3C: {"INC A", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionIncrement(gb, gb.Cpu.Registers.A)
	}},

	0x03: {"INC BC", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C++
		gb.Cpu.Registers.C &= 0xFF
		if gb.Cpu.Registers.C == 0 {
			gb.Cpu.Registers.B++
			gb.Cpu.Registers.B &= 0xFF
		}
	}},
	0x13: {"INC DE", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E++
		gb.Cpu.Registers.E &= 0xFF
		if gb.Cpu.Registers.E == 0 {
			gb.Cpu.Registers.D++
			gb.Cpu.Registers.D &= 0xFF
		}
	}},
	0x23: {"INC HL", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L++
		gb.Cpu.Registers.L &= 0xFF
		if gb.Cpu.Registers.L == 0 {
			gb.Cpu.Registers.H++
			gb.Cpu.Registers.H &= 0xFF
		}
	}},
	0x33: {"INC SP", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp++
		gb.Cpu.Registers.Sp &= 0xFFFF
	}},

	// DEC
	0x05: {"DEC B", 4, func(gb *Gameboy) {
		org := gb.Cpu.Registers.B
		gb.Cpu.Registers.B = (org - 1) & 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.B == 0)
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x0D: {"DEC C", 4, func(gb *Gameboy) {
		org := gb.Cpu.Registers.C
		gb.Cpu.Registers.C = (org - 1) & 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.C == 0)
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x15: {"DEC D", 4, func(gb *Gameboy) {
		org := gb.Cpu.Registers.D
		gb.Cpu.Registers.D = (org - 1) & 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.D == 0)
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x1D: {"DEC E", 4, func(gb *Gameboy) {
		org := gb.Cpu.Registers.E
		gb.Cpu.Registers.E = (org - 1) & 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.E == 0)
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x25: {"DEC H", 4, func(gb *Gameboy) {
		org := gb.Cpu.Registers.H
		gb.Cpu.Registers.H = (org - 1) & 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.H == 0)
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x2d: {"DEC L", 4, func(gb *Gameboy) {
		org := gb.Cpu.Registers.L
		gb.Cpu.Registers.L = (org - 1) & 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.L == 0)
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x35: {"DEC (HL)", 12, func(gb *Gameboy) {
		org := gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
		new := (org - 1) & 0xFF
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), new)
		gb.Cpu.Registers.SetFlagZ(new == 0)
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x3D: {"DEC A", 4, func(gb *Gameboy) {
		org := gb.Cpu.Registers.A
		gb.Cpu.Registers.A = (org - 1) & 0xFF
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(org&0x0F == 0)
	}},

	0x0B: {"DEC BC", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C--
		gb.Cpu.Registers.C &= 0xFF
		if gb.Cpu.Registers.C == 0xFF {
			gb.Cpu.Registers.B--
			gb.Cpu.Registers.B &= 0xFF
		}
	}},
	0x1B: {"DEC DE", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E--
		gb.Cpu.Registers.E &= 0xFF
		if gb.Cpu.Registers.E == 0xFF {
			gb.Cpu.Registers.D--
			gb.Cpu.Registers.D &= 0xFF
		}
	}},
	0x2B: {"DEC HL", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L--
		gb.Cpu.Registers.L &= 0xFF
		if gb.Cpu.Registers.L == 0xFF {
			gb.Cpu.Registers.H--
			gb.Cpu.Registers.H &= 0xFF
		}
	}},
	0x3B: {"DEC SP", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp--
		gb.Cpu.Registers.Sp &= 0xFFFF
	}},

	// ADD
	0x80: {"ADD A, B", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.B, false)
	}},
	0x81: {"ADD A, C", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.C, false)
	}},
	0x82: {"ADD A, D", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.D, false)
	}},
	0x83: {"ADD A, E", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.E, false)
	}},
	0x84: {"ADD A, H", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.H, false)
	}},
	0x85: {"ADD A, L", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.L, false)
	}},
	0x86: {"ADD A, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L)), false)
	}},
	0x87: {"ADD A, A", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.A, false)
	}},
	0xC6: {"ADD A, n", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.popPc(), false)
	}},
	0x09: {"ADD HL, BC", 8, func(gb *Gameboy) {
		result := instructionAddition16(gb, uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), uint16(gb.Cpu.Registers.B)<<8+uint16(gb.Cpu.Registers.C))
		gb.Cpu.Registers.H = uint8((result & 0xFF00) >> 8)
		gb.Cpu.Registers.L = uint8(result & 0x00FF)
	}},
	0x19: {"ADD HL, DE", 8, func(gb *Gameboy) {
		result := instructionAddition16(gb, uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), uint16(gb.Cpu.Registers.D)<<8+uint16(gb.Cpu.Registers.E))
		gb.Cpu.Registers.H = uint8((result & 0xFF00) >> 8)
		gb.Cpu.Registers.L = uint8(result & 0x00FF)
	}},
	0x29: {"ADD HL, HL", 8, func(gb *Gameboy) {
		result := instructionAddition16(gb, uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L))
		gb.Cpu.Registers.H = uint8((result & 0xFF00) >> 8)
		gb.Cpu.Registers.L = uint8(result & 0x00FF)
	}},
	0x39: {"ADD HL, SP", 8, func(gb *Gameboy) {
		result := instructionAddition16(gb, uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.H = uint8((result & 0xFF00) >> 8)
		gb.Cpu.Registers.L = uint8(result & 0x00FF)
	}},
	0xe8: {"ADD SP, n", 8, func(gb *Gameboy) {
		n := int8(gb.popPc())
		sum := uint16(int32(gb.Cpu.Registers.Sp) + int32(n))
		sumTmp := gb.Cpu.Registers.Sp ^ uint16(n) ^ sum
		gb.Cpu.Registers.Sp = sum
		gb.Cpu.Registers.SetFlagZ(false)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH((sumTmp & 0x10) == 0x10)
		gb.Cpu.Registers.SetFlagC((sumTmp & 0x100) == 0x100)
	}},
	0x88: {"ADC A, B", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.B, true)
	}},
	0x89: {"ADC A, C", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.C, true)
	}},
	0x8a: {"ADC A, D", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.D, true)
	}},
	0x8b: {"ADC A, E", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.E, true)
	}},
	0x8c: {"ADC A, H", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.H, true)
	}},
	0x8d: {"ADC A, L", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.L, true)
	}},
	0x8e: {"ADC A, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L)), true)
	}},
	0x8f: {"ADC A, A", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.A, true)
	}},
	0xce: {"ADC A, n", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionAddition(gb, gb.Cpu.Registers.A, gb.popPc(), true)
	}},

	0x90: {"SUB A, B", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.B, false)
	}},
	0x91: {"SUB A, C", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.C, false)
	}},
	0x92: {"SUB A, D", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.D, false)
	}},
	0x93: {"SUB A, E", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.E, false)
	}},
	0x94: {"SUB A, H", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.H, false)
	}},
	0x95: {"SUB A, L", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.L, false)
	}},
	0x96: {"SUB A, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L)), false)
	}},
	0x97: {"SUB A, A", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.A, false)
	}},
	0xD6: {"SUB A, n", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.popPc(), false)
	}},
	0x98: {"SBC A, B", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.B, true)
	}},
	0x99: {"SBC A, C", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.C, true)
	}},
	0x9a: {"SBC A, D", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.D, true)
	}},
	0x9b: {"SBC A, E", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.E, true)
	}},
	0x9c: {"SBC A, H", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.H, true)
	}},
	0x9d: {"SBC A, L", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.L, true)
	}},
	0x9e: {"SBC A, (HL)", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L)), true)
	}},
	0x9f: {"SBC A, A", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.A, true)
	}},
	0xde: {"SBC A, n", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSubstraction(gb, gb.Cpu.Registers.A, gb.popPc(), true)
	}},

	// Compares
	0xb8: {"CP A,B", 4, func(gb *Gameboy) {
		instructionCompare(gb, gb.Cpu.Registers.B, gb.Cpu.Registers.A)
	}},
	0xb9: {"CP A,C", 4, func(gb *Gameboy) {
		instructionCompare(gb, gb.Cpu.Registers.C, gb.Cpu.Registers.A)
	}},
	0xba: {"CP A,D", 4, func(gb *Gameboy) {
		instructionCompare(gb, gb.Cpu.Registers.D, gb.Cpu.Registers.A)
	}},
	0xbb: {"CP A,E", 4, func(gb *Gameboy) {
		instructionCompare(gb, gb.Cpu.Registers.E, gb.Cpu.Registers.A)
	}},
	0xbc: {"CP A,H", 4, func(gb *Gameboy) {
		instructionCompare(gb, gb.Cpu.Registers.H, gb.Cpu.Registers.A)
	}},
	0xbd: {"CP A,L", 4, func(gb *Gameboy) {
		instructionCompare(gb, gb.Cpu.Registers.L, gb.Cpu.Registers.A)
	}},
	0xbe: {"CP A,(HL)", 8, func(gb *Gameboy) {
		instructionCompare(gb, gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L)), gb.Cpu.Registers.A)
	}},
	0xbf: {"CP A,A", 4, func(gb *Gameboy) {
		instructionCompare(gb, gb.Cpu.Registers.A, gb.Cpu.Registers.A)
	}},
	0xfe: {"CP A,n", 8, func(gb *Gameboy) {
		instructionCompare(gb, gb.popPc(), gb.Cpu.Registers.A)
	}},

	// Jumps
	0xc3: {"JP nn", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Pc = gb.popPc16()
	}},
	0xc2: {"JPNZ nn", 12, func(gb *Gameboy) {
		jumpAddr := gb.popPc16()
		if !gb.Cpu.Registers.FlagZ() {
			gb.Cpu.Registers.Pc = jumpAddr
			gb.Cpu.ClockCycles += 4
		}
	}},
	0xca: {"JPZ nn", 12, func(gb *Gameboy) {
		jumpAddr := gb.popPc16()
		if gb.Cpu.Registers.FlagZ() {
			gb.Cpu.Registers.Pc = jumpAddr
			gb.Cpu.ClockCycles += 4
		}
	}},
	0xd2: {"JPNC nn", 12, func(gb *Gameboy) {
		jumpAddr := gb.popPc16()
		if !gb.Cpu.Registers.FlagC() {
			gb.Cpu.Registers.Pc = jumpAddr
			gb.Cpu.ClockCycles += 4
		}
	}},
	0xda: {"JPC nn", 12, func(gb *Gameboy) {
		jumpAddr := gb.popPc16()
		if gb.Cpu.Registers.FlagC() {
			gb.Cpu.Registers.Pc = jumpAddr
			gb.Cpu.ClockCycles += 4
		}
	}},
	0xe9: {"JP (HL)", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.Pc = uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
	}},
	0x18: {"JR n", 8, func(gb *Gameboy) {
		jumpDistance := int16(gb.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		gb.Cpu.Registers.Pc += uint16(jumpDistance)
		gb.Cpu.ClockCycles += 4
	}},
	0x20: {"JR NZ, n", 8, func(gb *Gameboy) {
		jumpDistance := int16(gb.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		if (gb.Cpu.Registers.Flags & 0x80) == 0x0 {
			gb.Cpu.Registers.Pc += uint16(jumpDistance)
			gb.Cpu.ClockCycles += 4
		}
	}},
	0x28: {"JR Z, n", 8, func(gb *Gameboy) {
		jumpDistance := int16(gb.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		if (gb.Cpu.Registers.Flags & 0x80) == 0x80 {
			gb.Cpu.Registers.Pc += uint16(jumpDistance)
			gb.Cpu.ClockCycles += 4
		}
	}},
	0x30: {"JR NC, n", 8, func(gb *Gameboy) {
		jumpDistance := int16(gb.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		if (gb.Cpu.Registers.Flags & 0x10) == 0x0 {
			gb.Cpu.Registers.Pc += uint16(jumpDistance)
			gb.Cpu.ClockCycles += 4
		}
	}},
	0x38: {"JR C, n", 8, func(gb *Gameboy) {
		jumpDistance := int16(gb.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		if (gb.Cpu.Registers.Flags & 0x10) == 0x10 {
			gb.Cpu.Registers.Pc += uint16(jumpDistance)
			gb.Cpu.ClockCycles += 4
		}
	}},

	// CALL
	0xcd: {"CALL nn", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := gb.popPc16()
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},
	0xc4: {"CALL NZ,nn", 12, func(gb *Gameboy) {
		jumpAddr := gb.popPc16()
		if (gb.Cpu.Registers.Flags & 0x80) == 0x0 {
			gb.Cpu.Registers.Sp -= 2
			gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
			gb.Cpu.ClockCycles += 8
			gb.Cpu.Registers.Pc = jumpAddr
		}
	}},
	0xcc: {"CALL Z,nn", 12, func(gb *Gameboy) {
		jumpAddr := gb.popPc16()
		if (gb.Cpu.Registers.Flags & 0x80) == 0x80 {
			gb.Cpu.Registers.Sp -= 2
			gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
			gb.Cpu.ClockCycles += 8
			gb.Cpu.Registers.Pc = jumpAddr
		}
	}},
	0xd4: {"CALL NC,nn", 12, func(gb *Gameboy) {
		jumpAddr := gb.popPc16()
		if (gb.Cpu.Registers.Flags & 0x10) == 0x0 {
			gb.Cpu.Registers.Sp -= 2
			gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
			gb.Cpu.ClockCycles += 8
			gb.Cpu.Registers.Pc = jumpAddr
		}
	}},
	0xdc: {"CALL C,nn", 12, func(gb *Gameboy) {
		jumpAddr := gb.popPc16()
		if (gb.Cpu.Registers.Flags & 0x10) == 0x10 {
			gb.Cpu.Registers.Sp -= 2
			gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
			gb.Cpu.ClockCycles += 8
			gb.Cpu.Registers.Pc = jumpAddr
		}
	}},

	// RST
	0xc7: {"RST 00", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := uint16(0x0000)
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},
	0xcf: {"RST 08", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := uint16(0x0008)
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},
	0xd7: {"RST 10", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := uint16(0x0010)
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},
	0xdf: {"RST 18", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := uint16(0x0018)
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},
	0xe7: {"RST 20", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := uint16(0x0020)
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},
	0xef: {"RST 28", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := uint16(0x0028)
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},
	0xf7: {"RST 30", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := uint16(0x0030)
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},
	0xff: {"RST 38", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Sp -= 2
		jumpAddr := uint16(0x0038)
		gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)
		gb.Cpu.Registers.Pc = jumpAddr
	}},

	// RET
	0xc9: {"RET", 12, func(gb *Gameboy) {
		gb.Cpu.Registers.Pc = gb.Memory.ReadWord(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp += 2
	}},
	0xc0: {"RET NZ", 4, func(gb *Gameboy) {
		if !gb.Cpu.Registers.FlagZ() {
			gb.Cpu.Registers.Pc = gb.Memory.ReadWord(gb.Cpu.Registers.Sp)
			gb.Cpu.Registers.Sp += 2
			gb.Cpu.ClockCycles += 8
		}
	}},
	0xc8: {"RET Z", 4, func(gb *Gameboy) {
		if gb.Cpu.Registers.FlagZ() {
			gb.Cpu.Registers.Pc = gb.Memory.ReadWord(gb.Cpu.Registers.Sp)
			gb.Cpu.Registers.Sp += 2
			gb.Cpu.ClockCycles += 8
		}
	}},
	0xd0: {"RET NC", 4, func(gb *Gameboy) {
		if !gb.Cpu.Registers.FlagC() {
			gb.Cpu.Registers.Pc = gb.Memory.ReadWord(gb.Cpu.Registers.Sp)
			gb.Cpu.Registers.Sp += 2
			gb.Cpu.ClockCycles += 8
		}
	}},
	0xd8: {"RET C", 4, func(gb *Gameboy) {
		if gb.Cpu.Registers.FlagC() {
			gb.Cpu.Registers.Pc = gb.Memory.ReadWord(gb.Cpu.Registers.Sp)
			gb.Cpu.Registers.Sp += 2
			gb.Cpu.ClockCycles += 8
		}
	}},
	0xd9: {"RETI", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.Pc = gb.Memory.ReadWord(gb.Cpu.Registers.Sp)
		gb.Cpu.Registers.Sp += 2
		gb.Cpu.Registers.Ime = true
	}},

	// ROT
	0x07: {"RLCA", 8, func(gb *Gameboy) {
		carry := gb.Cpu.Registers.A&0x80 != 0
		gb.Cpu.Registers.A = (gb.Cpu.Registers.A << 1) & 0xFF
		if carry {
			gb.Cpu.Registers.A |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(false)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carry)
	}},
	0x17: {"RLA", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.A&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.A = (gb.Cpu.Registers.A << 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.A |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(false)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},
	0x0F: {"RRCA", 8, func(gb *Gameboy) {
		carry := gb.Cpu.Registers.A&0x1 != 0
		gb.Cpu.Registers.A = (gb.Cpu.Registers.A >> 1) & 0xFF
		if carry {
			gb.Cpu.Registers.A |= 0x80
		}

		gb.Cpu.Registers.SetFlagZ(false)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carry)
	}},
	0x1F: {"RRA", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.A&0x1 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.A = (gb.Cpu.Registers.A >> 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.A |= 0x80
		}

		gb.Cpu.Registers.SetFlagZ(false)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},

	// Misc
	0x27: {"DAA", 4, func(gb *Gameboy) {
		if !gb.Cpu.Registers.FlagN() {
			if gb.Cpu.Registers.FlagC() || gb.Cpu.Registers.A > 0x99 {
				gb.Cpu.Registers.A += 0x60
				gb.Cpu.Registers.SetFlagC(true)
			}
			if gb.Cpu.Registers.FlagH() || gb.Cpu.Registers.A&0xF > 0x9 {
				gb.Cpu.Registers.A += 0x06
				gb.Cpu.Registers.SetFlagH(false)
			}
		} else if gb.Cpu.Registers.FlagC() && gb.Cpu.Registers.FlagH() {
			gb.Cpu.Registers.A += 0x9a
			gb.Cpu.Registers.SetFlagH(false)
		} else if gb.Cpu.Registers.FlagC() {
			gb.Cpu.Registers.A += 0xA0
		} else if gb.Cpu.Registers.FlagH() {
			gb.Cpu.Registers.A += 0xFA
			gb.Cpu.Registers.SetFlagH(false)
		}
		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0)
	}},
	0x2F: {"CPL", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.A = 0xFF ^ gb.Cpu.Registers.A
		gb.Cpu.Registers.SetFlagN(true)
		gb.Cpu.Registers.SetFlagH(true)
	}},
	0x37: {"SCF", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(true)
	}},
	0x3F: {"CCF", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(!gb.Cpu.Registers.FlagC())
	}},

	0xF3: {"DI", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.Ime = false
	}},
	0xFB: {"EI", 4, func(gb *Gameboy) {
		gb.Cpu.Registers.Ime = true
	}},

	// CB Mapper
	0xcb: {"PREFIX CB", 0, func(gb *Gameboy) {
		opcodeByte := gb.popPc()
		opcodeCb := OpcodesCb[opcodeByte]
		gb.Cpu.ClockCycles += opcodeCb.Cycles
		opcodeCb.Function(gb)
	}},
	0x76: {"HALT", 4, func(gb *Gameboy) {
		gb.Halted = true
	}},
	0x10: {"STOP", 4, func(gb *Gameboy) {
		gb.Halted = true
		gb.popPc()
	}},
	// Disallowed
	0xD3: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xD3, gb.Cpu.Registers.Pc)
	}},
	0xDB: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xDB, gb.Cpu.Registers.Pc)
	}},
	0xDD: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xDD, gb.Cpu.Registers.Pc)
	}},
	0xE3: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xE3, gb.Cpu.Registers.Pc)
	}},
	0xE4: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xE4, gb.Cpu.Registers.Pc)
	}},
	0xEB: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xEB, gb.Cpu.Registers.Pc)
	}},
	0xEC: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xEC, gb.Cpu.Registers.Pc)
	}},
	0xED: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xED, gb.Cpu.Registers.Pc)
	}},
	0xF4: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xF4, gb.Cpu.Registers.Pc)
	}},
	0xFC: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xFC, gb.Cpu.Registers.Pc)
	}},
	0xFD: {"ILLEGAL", 4, func(gb *Gameboy) {
		log.Panicf("Received illegal Opcode 0x%02x at 0x%04x", 0xFD, gb.Cpu.Registers.Pc)
	}},
}

func fillUninplementedOpcodes() {
	for k, v := range Opcodes {
		if v == nil {
			opcodeByte := k
			Opcodes[k] = &Opcode{
				fmt.Sprintf("UNIMP: %02x", k),
				1,
				func(gb *Gameboy) {
					log.Printf("Opcode not implemented: %02x", opcodeByte)
					log.Print(gb.String())
				},
			}
		}
	}
}
