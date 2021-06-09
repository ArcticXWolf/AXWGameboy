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

var opcodesCb = [0x100]*opcode{}

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

	// LD r1,r2
	0x78: {"LD A, B", 4, func(c *Cpu) { c.Registers.A = c.Registers.B }},
	0x79: {"LD A, C", 4, func(c *Cpu) { c.Registers.A = c.Registers.C }},
	0x7a: {"LD A, D", 4, func(c *Cpu) { c.Registers.A = c.Registers.D }},
	0x7b: {"LD A, E", 4, func(c *Cpu) { c.Registers.A = c.Registers.E }},
	0x7c: {"LD A, H", 4, func(c *Cpu) { c.Registers.A = c.Registers.H }},
	0x7d: {"LD A, L", 4, func(c *Cpu) { c.Registers.A = c.Registers.L }},
	0x7e: {"LD A, (HL)", 8, func(c *Cpu) { c.Registers.A = c.Memory.ReadByte(uint16(c.Registers.H)<<8 + uint16(c.Registers.L)) }},
	0x7f: {"LD A, A", 4, func(c *Cpu) { c.Registers.A = c.Registers.A }},

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
	0x36: {"LD (HL), n", 12, func(c *Cpu) { c.Memory.WriteByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L), c.popPc()) }},

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

	0x21: {"LD HL, nn", 12, func(c *Cpu) { c.Registers.L = c.popPc(); c.Registers.H = c.popPc() }},
	0x31: {"LD SP, nn", 12, func(c *Cpu) { c.Registers.Sp = c.popPc16() }},
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

	// CB Mapper
	0xcb: {"PREFIX CB", 0, func(c *Cpu) {
		opcodeByte := c.popPc()
		opcodeCb := opcodesCb[opcodeByte]
		log.Printf("Got next opcodeCB: Code: %02x, %s", opcodeByte, opcodeCb.Label)
		c.ClockCycles += opcodeCb.Cycles
		opcodeCb.Function(c)
	}},
}

func fillUninplementedOpcodes() {
	for k, v := range opcodes {
		if v == nil {
			opcodeByte := k
			opcodes[k] = &opcode{
				fmt.Sprintf("UNIMPLEMENTED: %02x", k),
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
				fmt.Sprintf("UNIMPLEMENTED CB: %02x", k),
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
