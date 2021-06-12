package internal

import (
	"image/color"
)

type Gpu struct {
	backgroundActivated bool
	backgroundMap       bool
	backgroundTile      bool
	lcdActivated        bool
	currentMode         uint8
	scrollX             uint8
	scrollY             uint8
	modeClock           int
	CurrentScanline     uint8
	tileSet             [512][8][8]uint8
	vram                [0x2000]byte
	bgPalette           [4]color.Color
}

func NewGpu() *Gpu {

	return &Gpu{
		bgPalette: [4]color.Color{
			color.RGBA{255, 255, 255, 255},
			color.RGBA{192, 192, 192, 255},
			color.RGBA{96, 96, 96, 255},
			color.RGBA{0, 0, 0, 255},
		},
	}
}

func (g *Gpu) ReadByte(address uint16) (result uint8) {
	switch address {
	case 0xFF40:
		var bA, bM, bT, lA uint8
		if g.backgroundActivated {
			bA = 0x01
		}
		if g.backgroundMap {
			bM = 0x08
		}
		if !g.backgroundTile {
			bT = 0x10
		}
		if g.lcdActivated {
			lA = 0x80
		}
		return bA | bM | bT | lA
	case 0xFF41:
		return g.currentMode & 0x03
	case 0xFF42:
		return g.scrollY
	case 0xFF43:
		return g.scrollX
	case 0xFF44:
		return g.CurrentScanline
	default:
		return 0x00
	}
}

func (g *Gpu) WriteByte(address uint16, value uint8) {
	switch address {
	case 0xFF40:
		g.backgroundActivated = value&0x01 != 0
		g.backgroundMap = value&0x08 != 0
		g.backgroundTile = value&0x10 == 0
		g.lcdActivated = value&0x80 != 0
	case 0xFF41:
		g.currentMode = value & 0x03
	case 0xFF42:
		g.scrollY = value
	case 0xFF43:
		g.scrollX = value
	case 0xFF47:
	default:
	}
}

func (g *Gpu) updateTile(address uint16) {
	address &= 0x1FFE
	tileIndex := (address >> 4) & 0x1FF
	y := (address >> 1) & 0x7

	for x := 0; x < 8; x++ {
		bitIndex := 1 << (7 - x)
		lowerBit := 0
		higherBit := 0

		if g.vram[address]&byte(bitIndex) != 0 {
			lowerBit = 1
		}
		if g.vram[address+1]&byte(bitIndex) != 0 {
			higherBit = 2
		}

		g.tileSet[tileIndex][x][y] = uint8(lowerBit) + uint8(higherBit)
	}
}

func (g *Gpu) Update(gb *Gameboy, cyclesUsed int) {
	if !g.lcdActivated {
		g.CurrentScanline = 0
		g.currentMode = 0
		gb.clearScreen()
		return
	}

	g.modeClock += cyclesUsed

	switch g.currentMode {
	case 0: // HBlank
		if g.modeClock >= 204 {
			g.modeClock = 0
			g.CurrentScanline++
			if g.CurrentScanline > uint8(ScreenHeight)-1 {
				g.currentMode = 1
				gb.ReadyToRender = gb.WorkingScreen
				gb.WorkingScreen = [ScreenWidth][ScreenHeight][3]uint8{}
			} else {
				g.currentMode = 2
			}
		}
	case 1: // VBlank
		if g.modeClock >= 456 {
			g.modeClock = 0
			g.CurrentScanline++
			if g.CurrentScanline > uint8(ScreenHeight)+10-1 {
				g.currentMode = 2
				g.CurrentScanline = 0
			}
		}
	case 2:
		if g.modeClock >= 80 {
			g.modeClock = 0
			g.currentMode = 3
		}
	case 3:
		if g.modeClock >= 172 {
			g.modeClock = 0
			g.currentMode = 0

			g.RenderScanline(gb)
		}
	}
}

func (g *Gpu) RenderScanline(gb *Gameboy) {
	g.renderTiles(gb)
}

func (g *Gpu) renderTiles(gb *Gameboy) {
	var mapOffset, lineOffset, tile uint16
	var x, y uint8

	mapOffset = 0x1800
	if g.backgroundMap {
		mapOffset = 0x1C00
	}
	mapOffset += ((uint16(g.CurrentScanline+g.scrollY) & 0xFF) >> 3) << 5

	lineOffset = uint16(g.scrollX) >> 3
	x = g.scrollX & 0x7
	y = (g.CurrentScanline + g.scrollY) & 0x7

	tile = uint16(g.vram[mapOffset+lineOffset])
	if g.backgroundTile && tile < 128 {
		tile += 256
	}

	for i := 0; i < ScreenWidth; i++ {
		red, green, blue, _ := g.bgPalette[g.tileSet[tile][x][y]].RGBA()
		gb.WorkingScreen[i][g.CurrentScanline][0] = uint8(red)
		gb.WorkingScreen[i][g.CurrentScanline][1] = uint8(green)
		gb.WorkingScreen[i][g.CurrentScanline][2] = uint8(blue)

		x++
		if x >= 8 {
			x = 0
			lineOffset = (lineOffset + 1) & 0x1F
			tile = uint16(g.vram[mapOffset+lineOffset])
			if g.backgroundTile && tile < 128 {
				tile += 256
			}
		}
	}
}

func (g *Gpu) Reset(gb *Gameboy) {
	g.CurrentScanline = 0
	g.currentMode = 0
	gb.clearScreen()

	g.tileSet = [512][8][8]uint8{}
}

func (gb *Gameboy) clearScreen() {
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			gb.WorkingScreen[x][y][0] = 255
			gb.WorkingScreen[x][y][1] = 255
			gb.WorkingScreen[x][y][2] = 255
		}
	}
	gb.ReadyToRender = gb.WorkingScreen
	gb.WorkingScreen = [ScreenWidth][ScreenHeight][3]uint8{}
}
