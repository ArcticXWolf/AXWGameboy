package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Debugger struct {
	Enabled bool
	Address uint16
	Step    bool
}

func (d *Debugger) checkBreakpoint(c *Cpu) {
	if d.Enabled && c.Registers.Pc == d.Address {
		d.Step = true
	}
}

func (d *Debugger) breakIfNecessary(c *Cpu) {
	if d.Step {
		log.Printf("%s", c.String())
		str := d.BreakExecution()
		if str == "s\n" {
			d.Step = false
		} else if str == "sp\n" {
			log.Printf("STACK: %02x %02x %02x %02x %02x %02x %02x %02x",
				c.Memory.ReadByte(c.Registers.Sp),
				c.Memory.ReadByte(c.Registers.Sp+1),
				c.Memory.ReadByte(c.Registers.Sp+2),
				c.Memory.ReadByte(c.Registers.Sp+3),
				c.Memory.ReadByte(c.Registers.Sp+4),
				c.Memory.ReadByte(c.Registers.Sp+5),
				c.Memory.ReadByte(c.Registers.Sp+6),
				c.Memory.ReadByte(c.Registers.Sp+7),
			)
			d.breakIfNecessary(c)
		} else if str == "hl\n" {
			log.Printf("HL: 0x%02x%02x, Value: 0x%02x",
				c.Registers.H,
				c.Registers.L,
				c.Memory.ReadByte(uint16(c.Registers.H)<<8+uint16(c.Registers.L)),
			)
			d.breakIfNecessary(c)
		} else if str == "q\n" {
			os.Exit(0)
		} else if ok, _ := regexp.MatchString("b[0-9a-fA-F]{4}\n", str); ok {
			d.Step = false
			addr, _ := strconv.ParseUint(str[1:5], 16, 16)
			d.Address = uint16(addr)
			log.Printf("Next breakpoint at 0x%04x", d.Address)
		}
	}
}

func (d *Debugger) BreakExecution() string {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Inputerror: %v", err)
		return d.BreakExecution()
	}
	return str
}
