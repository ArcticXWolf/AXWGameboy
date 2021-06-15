package internal

import (
	"image/color"
)

type SpriteObject struct {
	x                int
	y                int
	tile             uint8
	useSecondPalette bool
	xflip            bool
	yflip            bool
	priority         bool
}

type Gpu struct {
	gb                  *Gameboy
	backgroundActivated bool
	spritesActivated    bool
	bigSpritesActivated bool
	backgroundMap       bool
	backgroundTile      bool
	windowActivated     bool
	windowMap           bool
	lcdActivated        bool
	currentMode         uint8
	scrollX             uint8
	scrollY             uint8
	modeClock           int
	CurrentScanline     uint8
	tileSet             [512][8][8]uint8
	vram                [0x2000]byte
	oam                 [0xA0]byte
	spriteObjectData    [40]SpriteObject
	bgPaletteMap        [4]uint8
	bgPaletteColors     [4]color.Color
	spritePaletteMap    [2][4]uint8
	spritePaletteColors [2][4]color.Color
}

func NewGpu(gb *Gameboy) *Gpu {

	g := &Gpu{
		gb: gb,
		bgPaletteColors: [4]color.Color{
			color.RGBA{255, 255, 255, 255},
			color.RGBA{192, 192, 192, 255},
			color.RGBA{96, 96, 96, 255},
			color.RGBA{0, 0, 0, 255},
		},
		spritePaletteColors: [2][4]color.Color{
			{
				color.RGBA{255, 255, 255, 255},
				color.RGBA{192, 192, 192, 255},
				color.RGBA{96, 96, 96, 255},
				color.RGBA{0, 0, 0, 255},
			},
			{
				color.RGBA{255, 255, 255, 255},
				color.RGBA{192, 192, 192, 255},
				color.RGBA{96, 96, 96, 255},
				color.RGBA{0, 0, 0, 255},
			},
		},
	}

	for i := 0; i < len(g.spriteObjectData); i++ {
		g.spriteObjectData[i] = SpriteObject{
			x: -8,
			y: -16,
		}
	}

	return g
}

func (g *Gpu) ReadByte(address uint16) (result uint8) {
	switch address & 0xFF00 {
	case 0xFE00:
		if address < 0xFEA0 {
			return g.oam[address&0xFF]
		}
		return 0
	case 0xFF00:
		switch address {
		case 0xFF40:
			var bA, sA, bSA, bM, bT, wA, wM, lA uint8
			if g.backgroundActivated {
				bA = 0x01
			}
			if g.spritesActivated {
				sA = 0x02
			}
			if g.bigSpritesActivated {
				bSA = 0x04
			}
			if g.backgroundMap {
				bM = 0x08
			}
			if !g.backgroundTile {
				bT = 0x10
			}
			if g.windowActivated {
				wA = 0x20
			}
			if g.windowMap {
				wM = 0x40
			}
			if g.lcdActivated {
				lA = 0x80
			}
			return bA | sA | bSA | bM | bT | wA | wM | lA
		case 0xFF41:
			return g.currentMode & 0x03
		case 0xFF42:
			return g.scrollY
		case 0xFF43:
			return g.scrollX
		case 0xFF44:
			return g.CurrentScanline
		case 0xFF47:
			return g.bgPaletteMap[0]&0x3 | ((g.bgPaletteMap[1] & 0x3) << 2) | ((g.bgPaletteMap[2] & 0x3) << 4) | ((g.bgPaletteMap[3] & 0x3) << 6)
		default:
			return 0x00
		}
	default:
		return 0
	}
}

func (g *Gpu) WriteByte(address uint16, value uint8) {
	switch address & 0xFF00 {
	case 0xFE00:
		if address < 0xFEA0 {
			g.oam[address&0xFF] = value
			g.updateSpriteObject(address&0xFF, value)
		}
	case 0xFF00:
		switch address {
		case 0xFF40:
			g.backgroundActivated = value&0x01 != 0
			g.spritesActivated = value&0x02 != 0
			g.bigSpritesActivated = value&0x04 != 0
			g.backgroundMap = value&0x08 != 0
			g.backgroundTile = value&0x10 == 0
			g.windowActivated = value&0x20 != 0
			g.windowMap = value&0x40 != 0
			g.lcdActivated = value&0x80 != 0
		case 0xFF41:
			g.currentMode = value & 0x03
		case 0xFF42:
			g.scrollY = value
		case 0xFF43:
			g.scrollX = value
		case 0xFF46:
			g.oamDMA(value)
		case 0xFF47:
			g.bgPaletteMap[0] = value & 0x3
			g.bgPaletteMap[1] = (value >> 2) & 0x3
			g.bgPaletteMap[2] = (value >> 4) & 0x3
			g.bgPaletteMap[3] = (value >> 6) & 0x3
		case 0xFF48:
			g.spritePaletteMap[0][0] = value & 0x3
			g.spritePaletteMap[0][1] = (value >> 2) & 0x3
			g.spritePaletteMap[0][2] = (value >> 4) & 0x3
			g.spritePaletteMap[0][3] = (value >> 6) & 0x3
		case 0xFF49:
			g.spritePaletteMap[1][0] = value & 0x3
			g.spritePaletteMap[1][1] = (value >> 2) & 0x3
			g.spritePaletteMap[1][2] = (value >> 4) & 0x3
			g.spritePaletteMap[1][3] = (value >> 6) & 0x3
		default:
		}
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

		g.tileSet[tileIndex][y][x] = uint8(lowerBit) + uint8(higherBit)
	}
}

func (g *Gpu) oamDMA(value uint8) {
	var x uint16
	for x = 0; x < 0xA0; x++ {
		g.oam[x] = g.gb.Memory.ReadByte((uint16(value) << 8) + x)
		g.updateSpriteObject(x, g.oam[x])
	}
}

func (g *Gpu) updateSpriteObject(address uint16, value uint8) {
	objectId := address >> 2
	if objectId < uint16(len(g.spriteObjectData)) {
		switch address & 0x3 {
		case 0:
			g.spriteObjectData[objectId].y = int(value) - 16
		case 1:
			g.spriteObjectData[objectId].x = int(value) - 8
		case 2:
			g.spriteObjectData[objectId].tile = value
		case 3:
			// log.Printf("Updated tile %d: %v, Change was at address 0x%04x with value 0x%02x 0b%08b", objectId, g.spriteObjectData[objectId], address&0x3, value, value)
			g.spriteObjectData[objectId].useSecondPalette = false
			g.spriteObjectData[objectId].xflip = false
			g.spriteObjectData[objectId].yflip = false
			g.spriteObjectData[objectId].priority = false
			if value&0x10 > 0 {
				g.spriteObjectData[objectId].useSecondPalette = true
			}
			if value&0x20 > 0 {
				g.spriteObjectData[objectId].xflip = true
			}
			if value&0x40 > 0 {
				g.spriteObjectData[objectId].yflip = true
			}
			if value&0x80 > 0 {
				g.spriteObjectData[objectId].priority = true
			}
		}
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
				gb.Memory.GetInterruptFlags().TriggeredFlags |= 0x01
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
	var scanrow [ScreenWidth]byte
	if g.backgroundActivated {
		scanrow = g.renderTiles(gb)
	}

	if g.spritesActivated {
		g.renderSprites(gb, scanrow)
	}
}

func (g *Gpu) renderTiles(gb *Gameboy) (scanrow [ScreenWidth]byte) {
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
		pixelPaletteColor := g.tileSet[tile][y][x]
		scanrow[i] = pixelPaletteColor
		pixelRealColor := g.bgPaletteMap[pixelPaletteColor]
		red, green, blue, _ := g.bgPaletteColors[pixelRealColor].RGBA()
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

	return scanrow
}

func (g *Gpu) renderSprites(gb *Gameboy, scanrow [ScreenWidth]byte) {
	for i := 0; i < len(g.spriteObjectData); i++ {
		spriteObject := g.spriteObjectData[i]

		if spriteObject.y <= int(g.CurrentScanline) && spriteObject.y+8 > int(g.CurrentScanline) {
			palette := g.spritePaletteMap[0]
			paletteColors := g.spritePaletteColors[0]
			if spriteObject.useSecondPalette {
				palette = g.spritePaletteMap[1]
				paletteColors = g.spritePaletteColors[1]
			}
			var tilerow [8]uint8
			if spriteObject.yflip {
				tilerow = g.tileSet[spriteObject.tile][7-(g.CurrentScanline-uint8(spriteObject.y))]
			} else {
				tilerow = g.tileSet[spriteObject.tile][g.CurrentScanline-uint8(spriteObject.y)]
			}

			for x := 0; x < 8; x++ {
				if (spriteObject.x+x >= 0 && spriteObject.x+x < ScreenWidth) && tilerow[x] > 0 && (spriteObject.priority || scanrow[spriteObject.x+x] == 0) {
					pixelPaletteColor := tilerow[x]
					if spriteObject.xflip {
						pixelPaletteColor = tilerow[7-x]
					}
					pixelRealColor := palette[pixelPaletteColor]
					red, green, blue, _ := paletteColors[pixelRealColor].RGBA()
					gb.WorkingScreen[spriteObject.x+x][g.CurrentScanline][0] = uint8(red)
					gb.WorkingScreen[spriteObject.x+x][g.CurrentScanline][1] = uint8(green)
					gb.WorkingScreen[spriteObject.x+x][g.CurrentScanline][2] = uint8(blue)
				}
			}
		}
	}
}

func (g *Gpu) Reset(gb *Gameboy) {
	g.CurrentScanline = 0
	g.currentMode = 0
	gb.clearScreen()

	g.tileSet = [512][8][8]uint8{}
	g.oam = [0xA0]uint8{}
	g.spriteObjectData = [40]SpriteObject{}

	for i := 0; i < len(g.spriteObjectData); i++ {
		g.spriteObjectData[i] = SpriteObject{
			x: -8,
			y: -16,
		}
	}
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
