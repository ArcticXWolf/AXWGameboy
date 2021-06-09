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

var opcodes = [0x100]*opcode{
	0x00: {"NOP", 8, func(c *Cpu) {}},
	0x06: {"LD B, n", 8, func(c *Cpu) { c.Registers.B = c.Memory.ReadByte(c.Registers.Pc) }},
}

func fillUninplementedOpcodes() {
	for k, v := range opcodes {
		if v == nil {
			opcodes[k] = &opcode{
				fmt.Sprintf("UNIMPLEMENTED: %02x", k),
				1,
				func(c *Cpu) {
					log.Printf("Opcode not implemented: %02x", k)
					log.Print(c.String())
					utils.BreakExecution()
				},
			}
		}
	}
}
