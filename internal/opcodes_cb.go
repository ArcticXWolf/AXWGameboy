package internal

import (
	"fmt"
	"log"
)

var opcodesCb = [0x100]*opcode{
	0x10: {"RL B", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.B&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.B = (gb.Cpu.Registers.B << 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.B |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.B == 0x0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},
	0x11: {"RL C", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.C&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.C = (gb.Cpu.Registers.C << 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.C |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.C == 0x0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},
	0x12: {"RL D", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.D&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.D = (gb.Cpu.Registers.D << 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.D |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.D == 0x0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},
	0x13: {"RL E", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.E&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.E = (gb.Cpu.Registers.E << 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.E |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.E == 0x0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},
	0x14: {"RL H", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.H&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.H = (gb.Cpu.Registers.H << 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.H |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.H == 0x0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},
	0x15: {"RL L", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.L&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.L = (gb.Cpu.Registers.L << 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.L |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.L == 0x0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},
	0x16: {"RL (HL)", 16, func(gb *Gameboy) {
		hlValue := gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8 + uint16(gb.Cpu.Registers.L))
		carryNew := hlValue&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		rotation := (hlValue << 1) & 0xFF
		if carryOld {
			rotation |= 0x1
		}
		gb.Memory.WriteByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L), rotation)

		gb.Cpu.Registers.SetFlagZ(rotation == 0x0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},
	0x17: {"RL A", 8, func(gb *Gameboy) {
		carryNew := gb.Cpu.Registers.A&0x80 != 0
		carryOld := gb.Cpu.Registers.FlagC()
		gb.Cpu.Registers.A = (gb.Cpu.Registers.A << 1) & 0xFF
		if carryOld {
			gb.Cpu.Registers.A |= 0x1
		}

		gb.Cpu.Registers.SetFlagZ(gb.Cpu.Registers.A == 0x0)
		gb.Cpu.Registers.SetFlagN(false)
		gb.Cpu.Registers.SetFlagH(false)
		gb.Cpu.Registers.SetFlagC(carryNew)
	}},

	0x78: {"BIT 7, B", 8, func(gb *Gameboy) {
		if gb.Cpu.Registers.B&0x80 == 0x0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0x79: {"BIT 7, C", 8, func(gb *Gameboy) {
		if gb.Cpu.Registers.C&0x80 == 0x0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0x7a: {"BIT 7, D", 8, func(gb *Gameboy) {
		if gb.Cpu.Registers.D&0x80 == 0x0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0x7b: {"BIT 7, E", 8, func(gb *Gameboy) {
		if gb.Cpu.Registers.E&0x80 == 0x0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0x7c: {"BIT 7, H", 8, func(gb *Gameboy) {
		if gb.Cpu.Registers.H&0x80 == 0x0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0x7d: {"BIT 7, L", 8, func(gb *Gameboy) {
		if gb.Cpu.Registers.L&0x80 == 0x0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0x7e: {"BIT 7, (HL)", 12, func(gb *Gameboy) {
		if gb.Memory.ReadByte(uint16(gb.Cpu.Registers.H)<<8+uint16(gb.Cpu.Registers.L))&0x80 == 0x0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
	0x7f: {"BIT 7, A", 8, func(gb *Gameboy) {
		if gb.Cpu.Registers.A&0x80 == 0x0 {
			gb.Cpu.Registers.Flags = 0x80
		} else {
			gb.Cpu.Registers.Flags = 0x0
		}
	}},
}

func fillUninplementedOpcodesCb() {
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
