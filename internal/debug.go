package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Debugger struct {
	Enabled bool
	Address uint16
	Step    bool
}

func (d *Debugger) checkBreakpoint(gb *Gameboy) {
	if d.Enabled && gb.Cpu.Registers.Pc == d.Address {
		d.Step = true
	}
}

func (d *Debugger) breakIfNecessary(gb *Gameboy) {
	if d.Step {
		log.Printf("%s", gb.String())
		str := d.BreakExecution()
		if str == "c" {
			d.Step = false
		} else if ok, _ := regexp.MatchString("b[0-9a-fA-F]{4}", str); ok {
			d.Step = false
			addr, _ := strconv.ParseUint(str[1:5], 16, 16)
			d.Address = uint16(addr)
			log.Printf("Next breakpoint at 0x%04x", d.Address)
		} else if str == "gpu" {
			log.Printf("GPU: %v", gb.Gpu)
			d.dumpTileset(gb)
			d.breakIfNecessary(gb)
		} else if str == "du" {
			gb.Display.Render(gb)
			d.breakIfNecessary(gb)
		} else if str == "rtr" {
			d.dumpScreendata(gb.ReadyToRender)
			d.breakIfNecessary(gb)
		} else if str == "ws" {
			d.dumpScreendata(gb.WorkingScreen)
			d.breakIfNecessary(gb)
		} else if str == "sp" {
			log.Printf("STACK: Starts at 0x%04x, dumping for 0x%04x",
				gb.Cpu.Registers.Sp,
				0xFFFF-gb.Cpu.Registers.Sp,
			)
			d.dumpMemory(gb, gb.Cpu.Registers.Sp, 0xFFFF-gb.Cpu.Registers.Sp)
			d.breakIfNecessary(gb)
		} else if str == "mem" {
			d.dumpMemory(gb, 0, 0xFFFF)
			d.breakIfNecessary(gb)
		} else if str == "vmem" {
			d.dumpMemory(gb, 0x8000, 0x1FFF)
			d.breakIfNecessary(gb)
		} else if ok, _ := regexp.MatchString("m[0-9a-fA-F]{8}", str); ok {
			addr, _ := strconv.ParseUint(str[1:5], 16, 16)
			length, _ := strconv.ParseUint(str[5:9], 16, 16)
			d.dumpMemory(gb, uint16(addr), uint16(length))
			d.breakIfNecessary(gb)
		} else if str == "q" {
			os.Exit(0)
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
	return strings.TrimSpace(str)
}

func (d *Debugger) dumpMemory(gb *Gameboy, start uint16, length uint16) {
	start &= 0xFFF0
	length &= 0xFFF0
	var i uint32
	for i = uint32(start); i <= uint32(start+length); i += 0x10 {
		d.dump16Byte(gb, uint16(i))
	}
}

func (d *Debugger) dump16Byte(gb *Gameboy, address uint16) {
	b0 := gb.Memory.ReadByte(address)
	b1 := gb.Memory.ReadByte(address + 0x1)
	b2 := gb.Memory.ReadByte(address + 0x2)
	b3 := gb.Memory.ReadByte(address + 0x3)
	b4 := gb.Memory.ReadByte(address + 0x4)
	b5 := gb.Memory.ReadByte(address + 0x5)
	b6 := gb.Memory.ReadByte(address + 0x6)
	b7 := gb.Memory.ReadByte(address + 0x7)
	b8 := gb.Memory.ReadByte(address + 0x8)
	b9 := gb.Memory.ReadByte(address + 0x9)
	ba := gb.Memory.ReadByte(address + 0xa)
	bb := gb.Memory.ReadByte(address + 0xb)
	bc := gb.Memory.ReadByte(address + 0xc)
	bd := gb.Memory.ReadByte(address + 0xd)
	be := gb.Memory.ReadByte(address + 0xe)
	bf := gb.Memory.ReadByte(address + 0xf)
	log.Printf("(0x%04x)    %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x %02x%02x", address, b0, b1, b2, b3, b4, b5, b6, b7, b8, b9, ba, bb, bc, bd, be, bf)
}

func (d *Debugger) dumpScreendata(data [ScreenWidth][ScreenHeight][3]uint8) {
	for y := 0; y < ScreenHeight; y++ {
		lineStr := ""
		for x := 0; x < ScreenWidth; x++ {
			if data[x][y][0]+data[x][y][1]+data[x][y][2] > 0 {
				lineStr = fmt.Sprintf("%s#", lineStr)
			} else {
				lineStr = fmt.Sprintf("%s.", lineStr)
			}
		}
		log.Printf("%s", lineStr)
	}
}

func (d *Debugger) dumpTileset(gb *Gameboy) {
	for i := 0; i < 100; i++ {
		log.Printf("Tile %d:", i)
		for y := 0; y < 8; y++ {
			lineStr := ""
			for x := 0; x < 8; x++ {
				if gb.Gpu.tileSet[i][x][y] > 0 {
					lineStr = fmt.Sprintf("%s#", lineStr)
				} else {
					lineStr = fmt.Sprintf("%s.", lineStr)
				}
			}
			log.Printf("TILE |%s|", lineStr)
		}
	}
}
