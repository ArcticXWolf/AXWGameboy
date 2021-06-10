package cpu

import (
	"fmt"
	"log"

	"go.janniklasrichter.de/axwgameboy/internal/utils"
)

type opcode struct {
	Label    string
	Cycles   int
	Function func(*Cpu)
}

var opcodesCb = [0x100]*opcode{
	0x10: {"RL B", 8, func(c *Cpu) {
		carryNew := c.Registers.B&0x80 != 0
		carryOld := c.Registers.FlagC()
		c.Registers.B = (c.Registers.B << 1) & 0xFF
		if carryOld {
			c.Registers.B |= 0x1
		}

		c.Registers.SetFlagZ(c.Registers.B == 0x0)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},
	0x11: {"RL C", 8, func(c *Cpu) {
		carryNew := c.Registers.C&0x80 != 0
		carryOld := c.Registers.FlagC()
		c.Registers.C = (c.Registers.C << 1) & 0xFF
		if carryOld {
			c.Registers.C |= 0x1
		}

		c.Registers.SetFlagZ(c.Registers.C == 0x0)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},
	0x12: {"RL D", 8, func(c *Cpu) {
		carryNew := c.Registers.D&0x80 != 0
		carryOld := c.Registers.FlagC()
		c.Registers.D = (c.Registers.D << 1) & 0xFF
		if carryOld {
			c.Registers.D |= 0x1
		}

		c.Registers.SetFlagZ(c.Registers.D == 0x0)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},
	0x13: {"RL E", 8, func(c *Cpu) {
		carryNew := c.Registers.E&0x80 != 0
		carryOld := c.Registers.FlagC()
		c.Registers.E = (c.Registers.E << 1) & 0xFF
		if carryOld {
			c.Registers.E |= 0x1
		}

		c.Registers.SetFlagZ(c.Registers.E == 0x0)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},
	0x14: {"RL H", 8, func(c *Cpu) {
		carryNew := c.Registers.H&0x80 != 0
		carryOld := c.Registers.FlagC()
		c.Registers.H = (c.Registers.H << 1) & 0xFF
		if carryOld {
			c.Registers.H |= 0x1
		}

		c.Registers.SetFlagZ(c.Registers.H == 0x0)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},
	0x15: {"RL L", 8, func(c *Cpu) {
		carryNew := c.Registers.L&0x80 != 0
		carryOld := c.Registers.FlagC()
		c.Registers.L = (c.Registers.L << 1) & 0xFF
		if carryOld {
			c.Registers.L |= 0x1
		}

		c.Registers.SetFlagZ(c.Registers.L == 0x0)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},
	0x16: {"RL (HL)", 16, func(c *Cpu) {
		hlValue := c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L))
		carryNew := hlValue&0x80 != 0
		carryOld := c.Registers.FlagC()
		rotation := (hlValue << 1) & 0xFF
		if carryOld {
			rotation |= 0x1
		}
		c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), rotation)

		c.Registers.SetFlagZ(rotation == 0x0)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},
	0x17: {"RL A", 8, func(c *Cpu) {
		carryNew := c.Registers.A&0x80 != 0
		carryOld := c.Registers.FlagC()
		c.Registers.A = (c.Registers.A << 1) & 0xFF
		if carryOld {
			c.Registers.A |= 0x1
		}

		c.Registers.SetFlagZ(c.Registers.A == 0x0)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},

	0x78: {"BIT 7, B", 8, func(c *Cpu) {
		if c.Registers.B&0x80 == 0x0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x79: {"BIT 7, C", 8, func(c *Cpu) {
		if c.Registers.C&0x80 == 0x0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x7a: {"BIT 7, D", 8, func(c *Cpu) {
		if c.Registers.D&0x80 == 0x0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x7b: {"BIT 7, E", 8, func(c *Cpu) {
		if c.Registers.E&0x80 == 0x0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x7c: {"BIT 7, H", 8, func(c *Cpu) {
		if c.Registers.H&0x80 == 0x0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x7d: {"BIT 7, L", 8, func(c *Cpu) {
		if c.Registers.L&0x80 == 0x0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x7e: {"BIT 7, (HL)", 12, func(c *Cpu) {
		if c.Memory.ReadByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L))&0x80 == 0x0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x7f: {"BIT 7, A", 8, func(c *Cpu) {
		if c.Registers.A&0x80 == 0x0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
}

var opcodes = [0x100]*opcode{
	0x00: {"NOP", 4, func(c *Cpu) {}},
	// 8-Bit Loads
	// LD r1,n
	0x06: {"LD B, n", 8, func(c *Cpu) { c.Registers.B = c.popPc() }},
	0x0e: {"LD C, n", 8, func(c *Cpu) { c.Registers.C = c.popPc() }},
	0x16: {"LD D, n", 8, func(c *Cpu) { c.Registers.D = c.popPc() }},
	0x1e: {"LD E, n", 8, func(c *Cpu) { c.Registers.E = c.popPc() }},
	0x26: {"LD H, n", 8, func(c *Cpu) { c.Registers.H = c.popPc() }},
	0x2e: {"LD L, n", 8, func(c *Cpu) { c.Registers.L = c.popPc() }},
	0x36: {"LD (HL), n", 12, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.popPc()) }},
	0x3e: {"LD A, n", 8, func(c *Cpu) { c.Registers.A = c.popPc() }},

	// LD r1,r2
	0x40: {"LD B, B", 4, func(c *Cpu) { c.Registers.B = c.Registers.B }},
	0x41: {"LD B, C", 4, func(c *Cpu) { c.Registers.B = c.Registers.C }},
	0x42: {"LD B, D", 4, func(c *Cpu) { c.Registers.B = c.Registers.D }},
	0x43: {"LD B, E", 4, func(c *Cpu) { c.Registers.B = c.Registers.E }},
	0x44: {"LD B, H", 4, func(c *Cpu) { c.Registers.B = c.Registers.H }},
	0x45: {"LD B, L", 4, func(c *Cpu) { c.Registers.B = c.Registers.L }},
	0x46: {"LD B, (HL)", 8, func(c *Cpu) { c.Registers.B = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L)) }},

	0x48: {"LD C, B", 4, func(c *Cpu) { c.Registers.C = c.Registers.B }},
	0x49: {"LD C, C", 4, func(c *Cpu) { c.Registers.C = c.Registers.C }},
	0x4a: {"LD C, D", 4, func(c *Cpu) { c.Registers.C = c.Registers.D }},
	0x4b: {"LD C, E", 4, func(c *Cpu) { c.Registers.C = c.Registers.E }},
	0x4c: {"LD C, H", 4, func(c *Cpu) { c.Registers.C = c.Registers.H }},
	0x4d: {"LD C, L", 4, func(c *Cpu) { c.Registers.C = c.Registers.L }},
	0x4e: {"LD C, (HL)", 8, func(c *Cpu) { c.Registers.C = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L)) }},

	0x50: {"LD D, B", 4, func(c *Cpu) { c.Registers.D = c.Registers.B }},
	0x51: {"LD D, C", 4, func(c *Cpu) { c.Registers.D = c.Registers.C }},
	0x52: {"LD D, D", 4, func(c *Cpu) { c.Registers.D = c.Registers.D }},
	0x53: {"LD D, E", 4, func(c *Cpu) { c.Registers.D = c.Registers.E }},
	0x54: {"LD D, H", 4, func(c *Cpu) { c.Registers.D = c.Registers.H }},
	0x55: {"LD D, L", 4, func(c *Cpu) { c.Registers.D = c.Registers.L }},
	0x56: {"LD D, (HL)", 8, func(c *Cpu) { c.Registers.D = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L)) }},

	0x58: {"LD E, B", 4, func(c *Cpu) { c.Registers.E = c.Registers.B }},
	0x59: {"LD E, C", 4, func(c *Cpu) { c.Registers.E = c.Registers.C }},
	0x5a: {"LD E, D", 4, func(c *Cpu) { c.Registers.E = c.Registers.D }},
	0x5b: {"LD E, E", 4, func(c *Cpu) { c.Registers.E = c.Registers.E }},
	0x5c: {"LD E, H", 4, func(c *Cpu) { c.Registers.E = c.Registers.H }},
	0x5d: {"LD E, L", 4, func(c *Cpu) { c.Registers.E = c.Registers.L }},
	0x5e: {"LD E, (HL)", 8, func(c *Cpu) { c.Registers.E = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L)) }},

	0x60: {"LD H, B", 4, func(c *Cpu) { c.Registers.H = c.Registers.B }},
	0x61: {"LD H, C", 4, func(c *Cpu) { c.Registers.H = c.Registers.C }},
	0x62: {"LD H, D", 4, func(c *Cpu) { c.Registers.H = c.Registers.D }},
	0x63: {"LD H, E", 4, func(c *Cpu) { c.Registers.H = c.Registers.E }},
	0x64: {"LD H, H", 4, func(c *Cpu) { c.Registers.H = c.Registers.H }},
	0x65: {"LD H, L", 4, func(c *Cpu) { c.Registers.H = c.Registers.L }},
	0x66: {"LD H, (HL)", 8, func(c *Cpu) { c.Registers.H = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L)) }},

	0x68: {"LD L, B", 4, func(c *Cpu) { c.Registers.L = c.Registers.B }},
	0x69: {"LD L, C", 4, func(c *Cpu) { c.Registers.L = c.Registers.C }},
	0x6a: {"LD L, D", 4, func(c *Cpu) { c.Registers.L = c.Registers.D }},
	0x6b: {"LD L, E", 4, func(c *Cpu) { c.Registers.L = c.Registers.E }},
	0x6c: {"LD L, H", 4, func(c *Cpu) { c.Registers.L = c.Registers.H }},
	0x6d: {"LD L, L", 4, func(c *Cpu) { c.Registers.L = c.Registers.L }},
	0x6e: {"LD L, (HL)", 8, func(c *Cpu) { c.Registers.L = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L)) }},

	0x70: {"LD (HL), B", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.B) }},
	0x71: {"LD (HL), C", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.C) }},
	0x72: {"LD (HL), D", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.D) }},
	0x73: {"LD (HL), E", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.E) }},
	0x74: {"LD (HL), H", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.H) }},
	0x75: {"LD (HL), L", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.L) }},

	0x78: {"LD A, B", 4, func(c *Cpu) { c.Registers.A = c.Registers.B }},
	0x79: {"LD A, C", 4, func(c *Cpu) { c.Registers.A = c.Registers.C }},
	0x7a: {"LD A, D", 4, func(c *Cpu) { c.Registers.A = c.Registers.D }},
	0x7b: {"LD A, E", 4, func(c *Cpu) { c.Registers.A = c.Registers.E }},
	0x7c: {"LD A, H", 4, func(c *Cpu) { c.Registers.A = c.Registers.H }},
	0x7d: {"LD A, L", 4, func(c *Cpu) { c.Registers.A = c.Registers.L }},
	0x7e: {"LD A, (HL)", 8, func(c *Cpu) { c.Registers.A = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L)) }},
	0x7f: {"LD A, A", 4, func(c *Cpu) { c.Registers.A = c.Registers.A }},
	0x0a: {"LD A, (BC)", 8, func(c *Cpu) { c.Registers.A = c.Memory.ReadByte(uint16(c.Registers.B)<<8 + uint16(c.Registers.C)) }},
	0x1a: {"LD A, (DE)", 8, func(c *Cpu) { c.Registers.A = c.Memory.ReadByte(uint16(c.Registers.D)<<8 + uint16(c.Registers.E)) }},
	0xfa: {"LD A, (nn)", 16, func(c *Cpu) { c.Registers.A = c.Memory.ReadByte(c.popPc16()) }},

	0x47: {"LD B, A", 4, func(c *Cpu) { c.Registers.B = c.Registers.A }},
	0x4f: {"LD C, A", 4, func(c *Cpu) { c.Registers.C = c.Registers.A }},
	0x57: {"LD D, A", 4, func(c *Cpu) { c.Registers.D = c.Registers.A }},
	0x5f: {"LD E, A", 4, func(c *Cpu) { c.Registers.E = c.Registers.A }},
	0x67: {"LD H, A", 4, func(c *Cpu) { c.Registers.H = c.Registers.A }},
	0x6f: {"LD L, A", 4, func(c *Cpu) { c.Registers.L = c.Registers.A }},
	0x02: {"LD (BC), A", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.B)<<8+uint16(c.Registers.C), c.Registers.A) }},
	0x12: {"LD (DE), A", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.D)<<8+uint16(c.Registers.E), c.Registers.A) }},
	0x77: {"LD (HL), A", 8, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.A) }},
	0xEA: {"LD (HL), A", 16, func(c *Cpu) { c.Memory.WriteByte(c.popPc16(), c.Registers.A) }},

	// LD A and ($FF00+C)
	0xF2: {"LD A, ($FF00+C)", 8, func(c *Cpu) {
		c.Registers.A = c.Memory.ReadByte(uint16(0xFF00) + uint16(c.Registers.C))
	}},
	0xE2: {"LD ($FF00+C), A", 8, func(c *Cpu) {
		c.Memory.WriteByte(uint16(0xFF00)+uint16(c.Registers.C), c.Registers.A)
	}},

	// LD A and HL-
	0x3a: {"LD A, (HL-)", 8, func(c *Cpu) {
		c.Registers.A = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L))
		c.Registers.L = (c.Registers.L - 1) & 0xFF
		if c.Registers.L == 0xFF {
			c.Registers.H = (c.Registers.H - 1) & 0xFF
		}
	}},
	0x32: {"LD (HL-), A", 8, func(c *Cpu) {
		c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.A)
		c.Registers.L = (c.Registers.L - 1) & 0xFF
		if c.Registers.L == 0xFF {
			c.Registers.H = (c.Registers.H - 1) & 0xFF
		}
	}},

	// LD A and HL+
	0x2a: {"LD A, (HL+)", 8, func(c *Cpu) {
		c.Registers.A = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L))
		c.Registers.L = (c.Registers.L + 1) & 0xFF
		if c.Registers.L == 0x0 {
			c.Registers.H = (c.Registers.H + 1) & 0xFF
		}
	}},
	0x22: {"LD (HL+), A", 8, func(c *Cpu) {
		c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.Registers.A)
		c.Registers.L = (c.Registers.L + 1) & 0xFF
		if c.Registers.L == 0x0 {
			c.Registers.H = (c.Registers.H + 1) & 0xFF
		}
	}},

	// LD A and ($FF00+n)
	0xF0: {"LD A, ($FF00+n)", 8, func(c *Cpu) {
		c.Registers.A = c.Memory.ReadByte(uint16(0xFF00) + uint16(c.popPc()))
	}},
	0xE0: {"LD ($FF00+n), A", 8, func(c *Cpu) {
		c.Memory.WriteByte(uint16(0xFF00)+uint16(c.popPc()), c.Registers.A)
	}},

	// 16-Bit Loads
	0x01: {"LD BC, nn", 12, func(c *Cpu) { c.Registers.C = c.popPc(); c.Registers.B = c.popPc() }},
	0x11: {"LD DE, nn", 12, func(c *Cpu) { c.Registers.E = c.popPc(); c.Registers.D = c.popPc() }},
	0x21: {"LD HL, nn", 12, func(c *Cpu) { c.Registers.L = c.popPc(); c.Registers.H = c.popPc() }},
	0x31: {"LD SP, nn", 12, func(c *Cpu) { c.Registers.Sp = c.popPc16() }},

	0xc5: {"PUSH BC", 16, func(c *Cpu) {
		c.Registers.Sp--
		c.Memory.WriteByte(c.Registers.Sp, c.Registers.B)
		c.Registers.Sp--
		c.Memory.WriteByte(c.Registers.Sp, c.Registers.C)
	}},
	0xd5: {"PUSH DE", 16, func(c *Cpu) {
		c.Registers.Sp--
		c.Memory.WriteByte(c.Registers.Sp, c.Registers.D)
		c.Registers.Sp--
		c.Memory.WriteByte(c.Registers.Sp, c.Registers.E)
	}},
	0xe5: {"PUSH HL", 16, func(c *Cpu) {
		c.Registers.Sp--
		c.Memory.WriteByte(c.Registers.Sp, c.Registers.H)
		c.Registers.Sp--
		c.Memory.WriteByte(c.Registers.Sp, c.Registers.L)
	}},
	0xf5: {"PUSH AF", 16, func(c *Cpu) {
		c.Registers.Sp--
		c.Memory.WriteByte(c.Registers.Sp, c.Registers.A)
		c.Registers.Sp--
		c.Memory.WriteByte(c.Registers.Sp, c.Registers.Flags)
	}},

	0xc1: {"POP BC", 12, func(c *Cpu) {
		c.Registers.C = c.Memory.ReadByte(c.Registers.Sp)
		c.Registers.Sp++
		c.Registers.B = c.Memory.ReadByte(c.Registers.Sp)
		c.Registers.Sp++
	}},
	0xd1: {"POP DE", 12, func(c *Cpu) {
		c.Registers.E = c.Memory.ReadByte(c.Registers.Sp)
		c.Registers.Sp++
		c.Registers.D = c.Memory.ReadByte(c.Registers.Sp)
		c.Registers.Sp++
	}},
	0xe1: {"POP HL", 12, func(c *Cpu) {
		c.Registers.L = c.Memory.ReadByte(c.Registers.Sp)
		c.Registers.Sp++
		c.Registers.H = c.Memory.ReadByte(c.Registers.Sp)
		c.Registers.Sp++
	}},
	0xf1: {"POP AF", 12, func(c *Cpu) {
		c.Registers.Flags = c.Memory.ReadByte(c.Registers.Sp)
		c.Registers.Sp++
		c.Registers.A = c.Memory.ReadByte(c.Registers.Sp)
		c.Registers.Sp++
	}},

	// XOR
	0xee: {"XOR n", 8, func(c *Cpu) {
		c.Registers.A ^= c.popPc()
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0xaf: {"XOR A", 4, func(c *Cpu) {
		c.Registers.A ^= c.Registers.A
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0xa8: {"XOR B", 4, func(c *Cpu) {
		c.Registers.A ^= c.Registers.B
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0xa9: {"XOR C", 4, func(c *Cpu) {
		c.Registers.A ^= c.Registers.C
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0xaa: {"XOR D", 4, func(c *Cpu) {
		c.Registers.A ^= c.Registers.D
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0xab: {"XOR E", 4, func(c *Cpu) {
		c.Registers.A ^= c.Registers.E
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0xac: {"XOR H", 4, func(c *Cpu) {
		c.Registers.A ^= c.Registers.H
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0xad: {"XOR L", 4, func(c *Cpu) {
		c.Registers.A ^= c.Registers.L
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0xae: {"XOR (HL)", 8, func(c *Cpu) {
		c.Registers.A ^= c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L))
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},

	// INC r
	0x04: {"INC B", 4, func(c *Cpu) {
		c.Registers.B++
		c.Registers.B &= 0xFF
		if c.Registers.B == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x0C: {"INC C", 4, func(c *Cpu) {
		c.Registers.C++
		c.Registers.C &= 0xFF
		if c.Registers.C == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x14: {"INC D", 4, func(c *Cpu) {
		c.Registers.D++
		c.Registers.D &= 0xFF
		if c.Registers.D == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x1C: {"INC E", 4, func(c *Cpu) {
		c.Registers.E++
		c.Registers.E &= 0xFF
		if c.Registers.E == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x24: {"INC H", 4, func(c *Cpu) {
		c.Registers.H++
		c.Registers.H &= 0xFF
		if c.Registers.H == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x2C: {"INC L", 4, func(c *Cpu) {
		c.Registers.L++
		c.Registers.L &= 0xFF
		if c.Registers.L == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x34: {"INC (HL)", 4, func(c *Cpu) {
		incValue := c.Memory.ReadByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L)) + 1
		incValue &= 0xFF
		c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), incValue)
		if incValue == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},
	0x3C: {"INC A", 4, func(c *Cpu) {
		c.Registers.A++
		c.Registers.A &= 0xFF
		if c.Registers.A == 0 {
			c.Registers.Flags = 0x80
		} else {
			c.Registers.Flags = 0x0
		}
	}},

	0x03: {"INC BC", 8, func(c *Cpu) {
		c.Registers.C++
		c.Registers.C &= 0xFF
		if c.Registers.C == 0 {
			c.Registers.B++
			c.Registers.B &= 0xFF
		}
	}},
	0x13: {"INC DE", 8, func(c *Cpu) {
		c.Registers.E++
		c.Registers.E &= 0xFF
		if c.Registers.E == 0 {
			c.Registers.D++
			c.Registers.D &= 0xFF
		}
	}},
	0x23: {"INC HL", 8, func(c *Cpu) {
		c.Registers.L++
		c.Registers.L &= 0xFF
		if c.Registers.L == 0 {
			c.Registers.H++
			c.Registers.H &= 0xFF
		}
	}},
	0x33: {"INC SP", 8, func(c *Cpu) {
		c.Registers.Sp++
		c.Registers.Sp &= 0xFFFF
	}},

	// DEC
	0x05: {"DEC B", 4, func(c *Cpu) {
		org := c.Registers.B
		c.Registers.B = (org - 1) & 0xFF
		c.Registers.SetFlagZ(c.Registers.B == 0)
		c.Registers.SetFlagN(true)
		c.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x0D: {"DEC C", 4, func(c *Cpu) {
		org := c.Registers.C
		c.Registers.C = (org - 1) & 0xFF
		c.Registers.SetFlagZ(c.Registers.C == 0)
		c.Registers.SetFlagN(true)
		c.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x15: {"DEC D", 4, func(c *Cpu) {
		org := c.Registers.D
		c.Registers.D = (org - 1) & 0xFF
		c.Registers.SetFlagZ(c.Registers.D == 0)
		c.Registers.SetFlagN(true)
		c.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x1D: {"DEC E", 4, func(c *Cpu) {
		org := c.Registers.E
		c.Registers.E = (org - 1) & 0xFF
		c.Registers.SetFlagZ(c.Registers.E == 0)
		c.Registers.SetFlagN(true)
		c.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x25: {"DEC H", 4, func(c *Cpu) {
		org := c.Registers.H
		c.Registers.H = (org - 1) & 0xFF
		c.Registers.SetFlagZ(c.Registers.H == 0)
		c.Registers.SetFlagN(true)
		c.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x2d: {"DEC L", 4, func(c *Cpu) {
		org := c.Registers.L
		c.Registers.L = (org - 1) & 0xFF
		c.Registers.SetFlagZ(c.Registers.L == 0)
		c.Registers.SetFlagN(true)
		c.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x35: {"DEC (HL)", 12, func(c *Cpu) {
		org := c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L))
		new := (org - 1) & 0xFF
		c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), new)
		c.Registers.SetFlagZ(new == 0)
		c.Registers.SetFlagN(true)
		c.Registers.SetFlagH(org&0x0F == 0)
	}},
	0x3D: {"DEC A", 4, func(c *Cpu) {
		org := c.Registers.A
		c.Registers.A = (org - 1) & 0xFF
		c.Registers.SetFlagZ(c.Registers.A == 0)
		c.Registers.SetFlagN(true)
		c.Registers.SetFlagH(org&0x0F == 0)
	}},

	// ADD
	0x80: {"ADD A, B", 4, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.Registers.B, false)
	}},
	0x81: {"ADD A, C", 4, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.Registers.C, false)
	}},
	0x82: {"ADD A, D", 4, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.Registers.D, false)
	}},
	0x83: {"ADD A, E", 4, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.Registers.E, false)
	}},
	0x84: {"ADD A, H", 4, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.Registers.H, false)
	}},
	0x85: {"ADD A, L", 4, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.Registers.L, false)
	}},
	0x86: {"ADD A, (HL)", 8, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.Memory.ReadByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L)), false)
	}},
	0x87: {"ADD A, A", 4, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.Registers.A, false)
	}},
	0xC6: {"ADD A, n", 8, func(c *Cpu) {
		c.Registers.A = instructionAddition(c, c.Registers.A, c.popPc(), false)
	}},

	0x90: {"SUB A, B", 4, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.Registers.B, false)
	}},
	0x91: {"SUB A, C", 4, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.Registers.C, false)
	}},
	0x92: {"SUB A, D", 4, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.Registers.D, false)
	}},
	0x93: {"SUB A, E", 4, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.Registers.E, false)
	}},
	0x94: {"SUB A, H", 4, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.Registers.H, false)
	}},
	0x95: {"SUB A, L", 4, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.Registers.L, false)
	}},
	0x96: {"SUB A, (HL)", 8, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.Memory.ReadByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L)), false)
	}},
	0x97: {"SUB A, A", 4, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.Registers.A, false)
	}},
	0xD6: {"SUB A, n", 8, func(c *Cpu) {
		c.Registers.A = instructionSubstraction(c, c.Registers.A, c.popPc(), false)
	}},

	// Compares
	0xb8: {"CP A,B", 4, func(c *Cpu) {
		instructionCompare(c, c.Registers.B, c.Registers.A)
	}},
	0xb9: {"CP A,C", 4, func(c *Cpu) {
		instructionCompare(c, c.Registers.C, c.Registers.A)
	}},
	0xba: {"CP A,D", 4, func(c *Cpu) {
		instructionCompare(c, c.Registers.D, c.Registers.A)
	}},
	0xbb: {"CP A,E", 4, func(c *Cpu) {
		instructionCompare(c, c.Registers.E, c.Registers.A)
	}},
	0xbc: {"CP A,H", 4, func(c *Cpu) {
		instructionCompare(c, c.Registers.H, c.Registers.A)
	}},
	0xbd: {"CP A,L", 4, func(c *Cpu) {
		instructionCompare(c, c.Registers.L, c.Registers.A)
	}},
	0xbe: {"CP A,(HL)", 8, func(c *Cpu) {
		instructionCompare(c, c.Memory.ReadByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L)), c.Registers.A)
	}},
	0xbf: {"CP A,A", 4, func(c *Cpu) {
		instructionCompare(c, c.Registers.A, c.Registers.A)
	}},
	0xfe: {"CP A,n", 8, func(c *Cpu) {
		instructionCompare(c, c.popPc(), c.Registers.A)
	}},

	// Jumps
	0x18: {"JR n", 8, func(c *Cpu) {
		jumpDistance := int16(c.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		c.Registers.Pc += uint16(jumpDistance)
		c.ClockCycles += 4
	}},
	0x20: {"JR NZ, n", 8, func(c *Cpu) {
		jumpDistance := int16(c.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		if (c.Registers.Flags & 0x80) == 0x0 {
			c.Registers.Pc += uint16(jumpDistance)
			c.ClockCycles += 4
		}
	}},
	0x28: {"JR Z, n", 8, func(c *Cpu) {
		jumpDistance := int16(c.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		if (c.Registers.Flags & 0x80) == 0x80 {
			c.Registers.Pc += uint16(jumpDistance)
			c.ClockCycles += 4
		}
	}},
	0x30: {"JR NC, n", 8, func(c *Cpu) {
		jumpDistance := int16(c.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		if (c.Registers.Flags & 0x10) == 0x0 {
			c.Registers.Pc += uint16(jumpDistance)
			c.ClockCycles += 4
		}
	}},
	0x38: {"JR C, n", 8, func(c *Cpu) {
		jumpDistance := int16(c.popPc())
		if jumpDistance > 0x7F {
			jumpDistance = -((^jumpDistance + 1) & 0xFF)
		}
		if (c.Registers.Flags & 0x10) == 0x10 {
			c.Registers.Pc += uint16(jumpDistance)
			c.ClockCycles += 4
		}
	}},

	// CALL
	0xcd: {"CALL nn", 12, func(c *Cpu) {
		c.Registers.Sp -= 2
		jumpAddr := c.popPc16()
		c.Memory.WriteWord(c.Registers.Sp, c.Registers.Pc)
		c.Registers.Pc = jumpAddr
	}},
	0xc4: {"CALL NZ,nn", 12, func(c *Cpu) {
		jumpAddr := c.popPc16()
		if (c.Registers.Flags & 0x80) == 0x0 {
			c.Registers.Sp -= 2
			c.Memory.WriteWord(c.Registers.Sp, c.Registers.Pc)
			c.ClockCycles += 8
			c.Registers.Pc = jumpAddr
		}
	}},
	0xcc: {"CALL Z,nn", 12, func(c *Cpu) {
		jumpAddr := c.popPc16()
		if (c.Registers.Flags & 0x80) == 0x80 {
			c.Registers.Sp -= 2
			c.Memory.WriteWord(c.Registers.Sp, c.Registers.Pc)
			c.ClockCycles += 8
			c.Registers.Pc = jumpAddr
		}
	}},
	0xd4: {"CALL NC,nn", 12, func(c *Cpu) {
		jumpAddr := c.popPc16()
		if (c.Registers.Flags & 0x10) == 0x0 {
			c.Registers.Sp -= 2
			c.Memory.WriteWord(c.Registers.Sp, c.Registers.Pc)
			c.ClockCycles += 8
			c.Registers.Pc = jumpAddr
		}
	}},
	0xdc: {"CALL C,nn", 12, func(c *Cpu) {
		jumpAddr := c.popPc16()
		if (c.Registers.Flags & 0x10) == 0x10 {
			c.Registers.Sp -= 2
			c.Memory.WriteWord(c.Registers.Sp, c.Registers.Pc)
			c.ClockCycles += 8
			c.Registers.Pc = jumpAddr
		}
	}},

	// RET
	0xc9: {"RET", 12, func(c *Cpu) {
		c.Registers.Pc = c.Memory.ReadWord(c.Registers.Sp)
		c.Registers.Sp += 2
	}},
	0xc0: {"RET NZ", 4, func(c *Cpu) {
		if !c.Registers.FlagZ() {
			c.Registers.Pc = c.Memory.ReadWord(c.Registers.Sp)
			c.Registers.Sp += 2
			c.ClockCycles += 8
		}
	}},
	0xc8: {"RET Z", 4, func(c *Cpu) {
		if c.Registers.FlagZ() {
			c.Registers.Pc = c.Memory.ReadWord(c.Registers.Sp)
			c.Registers.Sp += 2
			c.ClockCycles += 8
		}
	}},
	0xd0: {"RET NC", 4, func(c *Cpu) {
		if !c.Registers.FlagC() {
			c.Registers.Pc = c.Memory.ReadWord(c.Registers.Sp)
			c.Registers.Sp += 2
			c.ClockCycles += 8
		}
	}},
	0xd8: {"RET C", 4, func(c *Cpu) {
		if c.Registers.FlagC() {
			c.Registers.Pc = c.Memory.ReadWord(c.Registers.Sp)
			c.Registers.Sp += 2
			c.ClockCycles += 8
		}
	}},

	// ROT
	0x17: {"RLA", 8, func(c *Cpu) {
		carryNew := c.Registers.A&0x80 != 0
		carryOld := c.Registers.FlagC()
		c.Registers.A = (c.Registers.A << 1) & 0xFF
		if carryOld {
			c.Registers.A |= 0x1
		}

		c.Registers.SetFlagZ(false)
		c.Registers.SetFlagN(false)
		c.Registers.SetFlagH(false)
		c.Registers.SetFlagC(carryNew)
	}},

	// CB Mapper
	0xcb: {"PREFIX CB", 0, func(c *Cpu) {
		opcodeByte := c.popPc()
		opcodeCb := opcodesCb[opcodeByte]
		c.ClockCycles += opcodeCb.Cycles
		opcodeCb.Function(c)
	}},
}

func fillUninplementedOpcodes() {
	for k, v := range opcodes {
		if v == nil {
			opcodeByte := k
			opcodes[k] = &opcode{
				fmt.Sprintf("UNIMP: %02x", k),
				1,
				func(c *Cpu) {
					log.Printf("Opcode not implemented: %02x", opcodeByte)
					log.Print(c.String())
					utils.BreakExecution()
				},
			}
		}
	}
}

func fillUninplementedOpcodesCb() {
	for k, v := range opcodesCb {
		if v == nil {
			opcodeByte := k
			opcodesCb[k] = &opcode{
				fmt.Sprintf("UNIMP CB: %02x", k),
				1,
				func(c *Cpu) {
					log.Printf("OpcodeCb not implemented: %02x", opcodeByte)
					log.Print(c.String())
					utils.BreakExecution()
				},
			}
		}
	}
}
