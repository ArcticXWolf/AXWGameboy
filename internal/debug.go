package internal

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Debugger struct {
	AddressEnabled  bool
	LogOnly         bool
	LogEvery        int
	logEveryCurrent int
	Address         uint16
	Step            bool
}

func (d *Debugger) checkBreakpoint(gb *Gameboy) {
	if d.AddressEnabled && gb.Cpu.Registers.Pc == d.Address {
		d.Step = true
	} else if d.LogOnly {
		if d.LogEvery > 0 {
			d.logEveryCurrent++
			if d.logEveryCurrent > d.LogEvery {
				d.logEveryCurrent = 0
			} else {
				return
			}
		}
		log.Printf("%s", gb.String())
	}
}

func (d *Debugger) triggerBreakpoint(gb *Gameboy) {
	// log.Printf("%s", gb.String())
	str := d.BreakExecution()
	if str == "c" {
		d.Step = false
	} else if str == "step" {
		d.Step = true
	} else if ok, _ := regexp.MatchString("b[0-9a-fA-F]{4}", str); ok {
		d.Step = false
		addr, _ := strconv.ParseUint(str[1:5], 16, 16)
		d.Address = uint16(addr)
		log.Printf("Next breakpoint at 0x%04x", d.Address)
	} else if str == "gpu" {
		log.Printf("GPU: %v", gb.Gpu)
		d.triggerBreakpoint(gb)
	} else if str == "ipl" {
		d.identifyPalettes(gb)
		d.triggerBreakpoint(gb)
	} else if ok, _ := regexp.MatchString("t[0-9]{3}", str); ok {
		tile, _ := strconv.ParseInt(str[1:4], 10, 0)
		d.dumpTile(gb, int(tile))
		d.triggerBreakpoint(gb)
	} else if str == "du" {
		gb.Display.Render(gb)
		d.triggerBreakpoint(gb)
	} else if str == "rtr" {
		d.dumpScreendata(gb.ReadyToRender)
		d.triggerBreakpoint(gb)
	} else if str == "ws" {
		d.dumpScreendata(gb.WorkingScreen)
		d.triggerBreakpoint(gb)
	} else if str == "sp" {
		log.Printf("STACK: Starts at 0x%04x, dumping for 0x%04x",
			gb.Cpu.Registers.Sp,
			0xFFFF-gb.Cpu.Registers.Sp,
		)
		d.dumpMemory(gb, gb.Cpu.Registers.Sp, 0xFFFF-gb.Cpu.Registers.Sp)
		d.triggerBreakpoint(gb)
	} else if str == "mem" {
		d.dumpMemory(gb, 0, 0xFFFF)
		d.triggerBreakpoint(gb)
	} else if str == "vmem" {
		d.dumpMemory(gb, 0x8000, 0x1FFF)
		d.triggerBreakpoint(gb)
	} else if ok, _ := regexp.MatchString("m[0-9a-fA-F]{8}", str); ok {
		addr, _ := strconv.ParseUint(str[1:5], 16, 16)
		length, _ := strconv.ParseUint(str[5:9], 16, 16)
		d.dumpMemory(gb, uint16(addr), uint16(length))
		d.triggerBreakpoint(gb)
	} else if str == "q" {
		os.Exit(0)
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
			if data[x][y][0]+data[x][y][1]+data[x][y][2] > 200 {
				lineStr = fmt.Sprintf("%s#", lineStr)
			} else if data[x][y][0]+data[x][y][1]+data[x][y][2] > 100 {
				lineStr = fmt.Sprintf("%s0", lineStr)
			} else if data[x][y][0]+data[x][y][1]+data[x][y][2] > 0 {
				lineStr = fmt.Sprintf("%s:", lineStr)
			} else {
				lineStr = fmt.Sprintf("%s.", lineStr)
			}
		}
		log.Printf("%s", lineStr)
	}
}

func (d *Debugger) dumpTileset(gb *Gameboy) {
	for i := 0; i < len(gb.Gpu.tileSet); i++ {
		log.Printf("Tile %d:", i)
		for y := 0; y < 8; y++ {
			lineStr := ""
			for x := 0; x < 8; x++ {
				if gb.Gpu.tileSet[i][x][y] == 3 {
					lineStr = fmt.Sprintf("%s#", lineStr)
				} else if gb.Gpu.tileSet[i][x][y] == 2 {
					lineStr = fmt.Sprintf("%s0", lineStr)
				} else if gb.Gpu.tileSet[i][x][y] == 1 {
					lineStr = fmt.Sprintf("%s:", lineStr)
				} else {
					lineStr = fmt.Sprintf("%s.", lineStr)
				}
			}
			log.Printf("TILE |%s|", lineStr)
		}
	}
}

func (d *Debugger) dumpTile(gb *Gameboy, index int) {
	log.Printf("Tile %d:", index)
	for y := 0; y < 8; y++ {
		lineStr := ""
		for x := 0; x < 8; x++ {
			if gb.Gpu.tileSet[index][x][y] == 3 {
				lineStr = fmt.Sprintf("%s#", lineStr)
			} else if gb.Gpu.tileSet[index][x][y] == 2 {
				lineStr = fmt.Sprintf("%s0", lineStr)
			} else if gb.Gpu.tileSet[index][x][y] == 1 {
				lineStr = fmt.Sprintf("%s:", lineStr)
			} else {
				lineStr = fmt.Sprintf("%s.", lineStr)
			}
		}
		log.Printf("TILE |%s|", lineStr)
	}
}

func (d *Debugger) identifyPalettes(gb *Gameboy) {
	gb.Gpu.bgPaletteColors = [4]color.Color{
		color.RGBA{255, 0, 0, 255},
		color.RGBA{192, 0, 0, 255},
		color.RGBA{96, 0, 0, 255},
		color.RGBA{50, 0, 0, 255},
	}
	gb.Gpu.spritePaletteColors = [2][4]color.Color{
		{
			color.RGBA{0, 255, 0, 255},
			color.RGBA{0, 192, 0, 255},
			color.RGBA{0, 96, 0, 255},
			color.RGBA{0, 50, 0, 255},
		},
		{
			color.RGBA{0, 0, 255, 255},
			color.RGBA{0, 0, 192, 255},
			color.RGBA{0, 0, 96, 255},
			color.RGBA{0, 0, 50, 255},
		},
	}
}
