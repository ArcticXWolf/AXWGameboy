package internal

import (
	"fmt"
	"log"
)

var opcodesCb = [0x100]*opcode{
	0x00: {"RLC B", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionCBRLC(gb, gb.Cpu.Registers.B)
	}},
	0x01: {"RLC C", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionCBRLC(gb, gb.Cpu.Registers.C)
	}},
	0x02: {"RLC D", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionCBRLC(gb, gb.Cpu.Registers.D)
	}},
	0x03: {"RLC E", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionCBRLC(gb, gb.Cpu.Registers.E)
	}},
	0x04: {"RLC H", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionCBRLC(gb, gb.Cpu.Registers.H)
	}},
	0x05: {"RLC L", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionCBRLC(gb, gb.Cpu.Registers.L)
	}},
	0x06: {"RLC (HL)", 16, func(gb *Gameboy) {
		addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
		gb.Memory.WriteByte(addr, instructionCBRLC(gb, gb.Memory.ReadByte(addr)))
	}},
	0x07: {"RLC A", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionCBRLC(gb, gb.Cpu.Registers.A)
	}},
	0x08: {"RRC B", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionCBRRC(gb, gb.Cpu.Registers.B)
	}},
	0x09: {"RRC C", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionCBRRC(gb, gb.Cpu.Registers.C)
	}},
	0x0a: {"RRC D", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionCBRRC(gb, gb.Cpu.Registers.D)
	}},
	0x0b: {"RRC E", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionCBRRC(gb, gb.Cpu.Registers.E)
	}},
	0x0c: {"RRC H", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionCBRRC(gb, gb.Cpu.Registers.H)
	}},
	0x0d: {"RRC L", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionCBRRC(gb, gb.Cpu.Registers.L)
	}},
	0x0e: {"RRC (HL)", 16, func(gb *Gameboy) {
		addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
		gb.Memory.WriteByte(addr, instructionCBRRC(gb, gb.Memory.ReadByte(addr)))
	}},
	0x0f: {"RRC A", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionCBRRC(gb, gb.Cpu.Registers.A)
	}},
	0x10: {"RL B", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionCBRL(gb, gb.Cpu.Registers.B)
	}},
	0x11: {"RL C", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionCBRL(gb, gb.Cpu.Registers.C)
	}},
	0x12: {"RL D", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionCBRL(gb, gb.Cpu.Registers.D)
	}},
	0x13: {"RL E", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionCBRL(gb, gb.Cpu.Registers.E)
	}},
	0x14: {"RL H", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionCBRL(gb, gb.Cpu.Registers.H)
	}},
	0x15: {"RL L", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionCBRL(gb, gb.Cpu.Registers.L)
	}},
	0x16: {"RL (HL)", 16, func(gb *Gameboy) {
		addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
		gb.Memory.WriteByte(addr, instructionCBRL(gb, gb.Memory.ReadByte(addr)))
	}},
	0x17: {"RL A", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionCBRL(gb, gb.Cpu.Registers.A)
	}},
	0x18: {"RR B", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionCBRR(gb, gb.Cpu.Registers.B)
	}},
	0x19: {"RR C", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionCBRR(gb, gb.Cpu.Registers.C)
	}},
	0x1a: {"RR D", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionCBRR(gb, gb.Cpu.Registers.D)
	}},
	0x1b: {"RR E", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionCBRR(gb, gb.Cpu.Registers.E)
	}},
	0x1c: {"RR H", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionCBRR(gb, gb.Cpu.Registers.H)
	}},
	0x1d: {"RR L", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionCBRR(gb, gb.Cpu.Registers.L)
	}},
	0x1e: {"RR (HL)", 16, func(gb *Gameboy) {
		addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
		gb.Memory.WriteByte(addr, instructionCBRR(gb, gb.Memory.ReadByte(addr)))
	}},
	0x1f: {"RR A", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionCBRR(gb, gb.Cpu.Registers.A)
	}},
	0x20: {"SLA B", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionCBSLA(gb, gb.Cpu.Registers.B)
	}},
	0x21: {"SLA C", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionCBSLA(gb, gb.Cpu.Registers.C)
	}},
	0x22: {"SLA D", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionCBSLA(gb, gb.Cpu.Registers.D)
	}},
	0x23: {"SLA E", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionCBSLA(gb, gb.Cpu.Registers.E)
	}},
	0x24: {"SLA H", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionCBSLA(gb, gb.Cpu.Registers.H)
	}},
	0x25: {"SLA L", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionCBSLA(gb, gb.Cpu.Registers.L)
	}},
	0x26: {"SLA (HL)", 16, func(gb *Gameboy) {
		addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
		gb.Memory.WriteByte(addr, instructionCBSLA(gb, gb.Memory.ReadByte(addr)))
	}},
	0x27: {"SLA A", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionCBSLA(gb, gb.Cpu.Registers.A)
	}},
	0x28: {"SRA B", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionCBSRA(gb, gb.Cpu.Registers.B)
	}},
	0x29: {"SRA C", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionCBSRA(gb, gb.Cpu.Registers.C)
	}},
	0x2a: {"SRA D", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionCBSRA(gb, gb.Cpu.Registers.D)
	}},
	0x2b: {"SRA E", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionCBSRA(gb, gb.Cpu.Registers.E)
	}},
	0x2c: {"SRA H", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionCBSRA(gb, gb.Cpu.Registers.H)
	}},
	0x2d: {"SRA L", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionCBSRA(gb, gb.Cpu.Registers.L)
	}},
	0x2e: {"SRA (HL)", 16, func(gb *Gameboy) {
		addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
		gb.Memory.WriteByte(addr, instructionCBSRA(gb, gb.Memory.ReadByte(addr)))
	}},
	0x2f: {"SRA A", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionCBSRA(gb, gb.Cpu.Registers.A)
	}},
	0x30: {"SWAP B", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionSwap(gb, gb.Cpu.Registers.B)
	}},
	0x31: {"SWAP C", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionSwap(gb, gb.Cpu.Registers.C)
	}},
	0x32: {"SWAP D", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionSwap(gb, gb.Cpu.Registers.D)
	}},
	0x33: {"SWAP E", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionSwap(gb, gb.Cpu.Registers.E)
	}},
	0x34: {"SWAP H", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionSwap(gb, gb.Cpu.Registers.H)
	}},
	0x35: {"SWAP L", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionSwap(gb, gb.Cpu.Registers.L)
	}},
	0x36: {"SWAP (HL)", 16, func(gb *Gameboy) {
		addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
		gb.Memory.WriteByte(addr, instructionSwap(gb, gb.Memory.ReadByte(addr)))
	}},
	0x37: {"SWAP A", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionSwap(gb, gb.Cpu.Registers.A)
	}},
	0x38: {"SRL B", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.B = instructionCBSRL(gb, gb.Cpu.Registers.B)
	}},
	0x39: {"SRL C", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.C = instructionCBSRL(gb, gb.Cpu.Registers.C)
	}},
	0x3a: {"SRL D", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.D = instructionCBSRL(gb, gb.Cpu.Registers.D)
	}},
	0x3b: {"SRL E", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.E = instructionCBSRL(gb, gb.Cpu.Registers.E)
	}},
	0x3c: {"SRL H", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.H = instructionCBSRL(gb, gb.Cpu.Registers.H)
	}},
	0x3d: {"SRL L", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.L = instructionCBSRL(gb, gb.Cpu.Registers.L)
	}},
	0x3e: {"SRL (HL)", 16, func(gb *Gameboy) {
		addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
		gb.Memory.WriteByte(addr, instructionCBSRL(gb, gb.Memory.ReadByte(addr)))
	}},
	0x3f: {"SRL A", 8, func(gb *Gameboy) {
		gb.Cpu.Registers.A = instructionCBSRL(gb, gb.Cpu.Registers.A)
	}},
}

func fillOpcodesCb(handleUnimplemented bool) {
	// fill repeating patterns
	for x := 0; x < 8; x++ {
		i := x

		// BIT Tests
		opcodesCb[0x40+0x8*x] = &opcode{
			fmt.Sprintf("BIT %d, B", i),
			8,
			func(gb *Gameboy) {
				instructionTestBit(gb, gb.Cpu.Registers.B, uint8(i))
			},
		}
		opcodesCb[0x41+0x8*x] = &opcode{
			fmt.Sprintf("BIT %d, C", i),
			8,
			func(gb *Gameboy) {
				instructionTestBit(gb, gb.Cpu.Registers.C, uint8(i))
			},
		}
		opcodesCb[0x42+0x8*x] = &opcode{
			fmt.Sprintf("BIT %d, D", i),
			8,
			func(gb *Gameboy) {
				instructionTestBit(gb, gb.Cpu.Registers.D, uint8(i))
			},
		}
		opcodesCb[0x43+0x8*x] = &opcode{
			fmt.Sprintf("BIT %d, E", i),
			8,
			func(gb *Gameboy) {
				instructionTestBit(gb, gb.Cpu.Registers.E, uint8(i))
			},
		}
		opcodesCb[0x44+0x8*x] = &opcode{
			fmt.Sprintf("BIT %d, H", i),
			8,
			func(gb *Gameboy) {
				instructionTestBit(gb, gb.Cpu.Registers.H, uint8(i))
			},
		}
		opcodesCb[0x45+0x8*x] = &opcode{
			fmt.Sprintf("BIT %d, L", i),
			8,
			func(gb *Gameboy) {
				instructionTestBit(gb, gb.Cpu.Registers.L, uint8(i))
			},
		}
		opcodesCb[0x46+0x8*x] = &opcode{
			fmt.Sprintf("BIT %d, (HL)", i),
			16,
			func(gb *Gameboy) {
				instructionTestBit(gb, gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L)), uint8(i))
			},
		}
		opcodesCb[0x47+0x8*x] = &opcode{
			fmt.Sprintf("BIT %d, A", i),
			8,
			func(gb *Gameboy) {
				instructionTestBit(gb, gb.Cpu.Registers.A, uint8(i))
			},
		}

		// RES
		opcodesCb[0x80+0x8*x] = &opcode{
			fmt.Sprintf("RES %d, B", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.B = instructionResetBit(gb, gb.Cpu.Registers.B, uint8(i))
			},
		}
		opcodesCb[0x81+0x8*x] = &opcode{
			fmt.Sprintf("RES %d, C", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.C = instructionResetBit(gb, gb.Cpu.Registers.C, uint8(i))
			},
		}
		opcodesCb[0x82+0x8*x] = &opcode{
			fmt.Sprintf("RES %d, D", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.D = instructionResetBit(gb, gb.Cpu.Registers.D, uint8(i))
			},
		}
		opcodesCb[0x83+0x8*x] = &opcode{
			fmt.Sprintf("RES %d, E", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.E = instructionResetBit(gb, gb.Cpu.Registers.E, uint8(i))
			},
		}
		opcodesCb[0x84+0x8*x] = &opcode{
			fmt.Sprintf("RES %d, H", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.H = instructionResetBit(gb, gb.Cpu.Registers.H, uint8(i))
			},
		}
		opcodesCb[0x85+0x8*x] = &opcode{
			fmt.Sprintf("RES %d, L", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.L = instructionResetBit(gb, gb.Cpu.Registers.L, uint8(i))
			},
		}
		opcodesCb[0x86+0x8*x] = &opcode{
			fmt.Sprintf("RES %d, (HL)", i),
			16,
			func(gb *Gameboy) {
				addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
				gb.Memory.WriteByte(addr, instructionResetBit(gb, gb.Memory.ReadByte(addr), uint8(i)))
			},
		}
		opcodesCb[0x87+0x8*x] = &opcode{
			fmt.Sprintf("RES %d, A", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.A = instructionResetBit(gb, gb.Cpu.Registers.A, uint8(i))
			},
		}

		// SET
		opcodesCb[0xc0+0x8*x] = &opcode{
			fmt.Sprintf("SET %d, B", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.B = instructionSetBit(gb, gb.Cpu.Registers.B, uint8(i))
			},
		}
		opcodesCb[0xc1+0x8*x] = &opcode{
			fmt.Sprintf("SET %d, C", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.C = instructionSetBit(gb, gb.Cpu.Registers.C, uint8(i))
			},
		}
		opcodesCb[0xc2+0x8*x] = &opcode{
			fmt.Sprintf("SET %d, D", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.D = instructionSetBit(gb, gb.Cpu.Registers.D, uint8(i))
			},
		}
		opcodesCb[0xc3+0x8*x] = &opcode{
			fmt.Sprintf("SET %d, E", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.E = instructionSetBit(gb, gb.Cpu.Registers.E, uint8(i))
			},
		}
		opcodesCb[0xc4+0x8*x] = &opcode{
			fmt.Sprintf("SET %d, H", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.H = instructionSetBit(gb, gb.Cpu.Registers.H, uint8(i))
			},
		}
		opcodesCb[0xc5+0x8*x] = &opcode{
			fmt.Sprintf("SET %d, L", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.L = instructionSetBit(gb, gb.Cpu.Registers.L, uint8(i))
			},
		}
		opcodesCb[0xc6+0x8*x] = &opcode{
			fmt.Sprintf("SET %d, (HL)", i),
			16,
			func(gb *Gameboy) {
				addr := uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L)
				gb.Memory.WriteByte(addr, instructionSetBit(gb, gb.Memory.ReadByte(addr), uint8(i)))
			},
		}
		opcodesCb[0xc7+0x8*x] = &opcode{
			fmt.Sprintf("SET %d, A", i),
			8,
			func(gb *Gameboy) {
				gb.Cpu.Registers.A = instructionSetBit(gb, gb.Cpu.Registers.A, uint8(i))
			},
		}
	}

	if handleUnimplemented {
		// Add handler for missing ones
		for k, v := range opcodesCb {
			if v == nil {
				opcodeByte := k
				opcodesCb[k] = &opcode{
					fmt.Sprintf("UNIMP CB: %02x", k),
					1,
					func(gb *Gameboy) {
						log.Printf("OpcodeCb not implemented: %02x", opcodeByte)
						log.Print(gb.String())
					},
				}
			}
		}
	}
}
