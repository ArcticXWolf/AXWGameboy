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
	vramBank         int
	cgbPalette       int
}

type TileAttributes struct {
	useBgPriorityInsteadOfOam bool
	yFlip                     bool
	xFlip                     bool
	tileVramBank              int
	bgPaletteNumber           int
}

type Gpu struct {
	gb *Gameboy

	backgroundActivated bool
	spritesActivated    bool
	bigSpritesActivated bool
	backgroundMap       bool
	backgroundTile      bool
	windowActivated     bool
	windowMap           bool
	lcdActivated        bool

	currentMode        uint8
	StatTriggerLYC     bool
	StatTriggerMode0   bool
	StatTriggerMode1   bool
	StatTriggerMode2   bool
	StatEnableMode0    bool
	StatEnableMode1    bool
	StatEnableMode2    bool
	StatEnableLYC      bool
	StatInterruptDelay bool

	scrollX              uint8
	scrollY              uint8
	windowX              uint8
	windowY              uint8
	ScanlineCompare      uint8
	CurrentScanline      uint8
	modeClock            int
	tileSet              [768][8][8]uint8
	TileAttributes       [0x800]TileAttributes
	vramBank             int
	vram                 [0x4000]byte
	oam                  [0xA0]byte
	SpriteObjectData     [40]SpriteObject
	BgPaletteMap         [4]uint8
	bgPaletteColors      [4]color.Color
	CgbBgPaletteColors   [8][4]color.Color
	SpritePaletteMap     [2][4]uint8
	spritePaletteColors  [2][4]color.Color
	CgbObjPaletteColors  [8][4]color.Color
	cgbBCPS              uint8
	cgbBCPSAutoincrement bool
	cgbOCPS              uint8
	cgbOCPSAutoincrement bool
}

func NewGpu(gb *Gameboy) *Gpu {

	g := &Gpu{
		gb:              gb,
		bgPaletteColors: getPaletteColorsByName(gb.Options.Palette),
		spritePaletteColors: [2][4]color.Color{
			getPaletteColorsByName(gb.Options.Palette),
			getPaletteColorsByName(gb.Options.Palette),
		},
	}

	for i := 0; i < len(g.CgbBgPaletteColors); i++ {
		g.CgbBgPaletteColors[i] = [4]color.Color{
			color.Black,
			color.Black,
			color.Black,
			color.Black,
		}
	}
	for i := 0; i < len(g.CgbObjPaletteColors); i++ {
		g.CgbObjPaletteColors[i] = [4]color.Color{
			color.Black,
			color.Black,
			color.Black,
			color.Black,
		}
	}

	for i := 0; i < len(g.SpriteObjectData); i++ {
		g.SpriteObjectData[i] = SpriteObject{
			x: -8,
			y: -16,
		}
	}
	for i := 0; i < len(g.TileAttributes); i++ {
		g.TileAttributes[i] = TileAttributes{}
	}

	return g
}

func (g *Gpu) ReadByte(address uint16) (result uint8) {
	switch address & 0xF000 {
	case 0x8000, 0x9000: // VRAM
		return g.vram[address&0x1FFF+uint16(g.vramBank)*0x2000]
	case 0xF000:
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
				value := g.currentMode & 0x3
				if g.StatTriggerLYC {
					value |= 0x4
				}
				if g.StatEnableMode0 {
					value |= 0x8
				}
				if g.StatEnableMode1 {
					value |= 0x10
				}
				if g.StatEnableMode2 {
					value |= 0x20
				}
				if g.StatEnableLYC {
					value |= 0x40
				}
				return value
			case 0xFF42:
				return g.scrollY
			case 0xFF43:
				return g.scrollX
			case 0xFF44:
				return g.CurrentScanline
			case 0xFF45:
				return g.ScanlineCompare
			case 0xFF47:
				return g.BgPaletteMap[0]&0x3 | ((g.BgPaletteMap[1] & 0x3) << 2) | ((g.BgPaletteMap[2] & 0x3) << 4) | ((g.BgPaletteMap[3] & 0x3) << 6)
			case 0xFF48:
				return g.SpritePaletteMap[0][0]&0x3 | ((g.SpritePaletteMap[0][1] & 0x3) << 2) | ((g.SpritePaletteMap[0][2] & 0x3) << 4) | ((g.SpritePaletteMap[0][3] & 0x3) << 6)
			case 0xFF49:
				return g.SpritePaletteMap[1][0]&0x3 | ((g.SpritePaletteMap[1][1] & 0x3) << 2) | ((g.SpritePaletteMap[1][2] & 0x3) << 4) | ((g.SpritePaletteMap[1][3] & 0x3) << 6)
			case 0xFF4A:
				return g.windowY
			case 0xFF4B:
				return g.windowX
			case 0xFF4F:
				return uint8(g.vramBank)
			case 0xFF68:
				if g.gb.cgbModeEnabled {
					if g.cgbBCPSAutoincrement {
						return (0x1 << 7) | g.cgbBCPS
					}
					return g.cgbBCPS
				}
				return 0x00
			case 0xFF69:
				if g.gb.cgbModeEnabled {
					paletteNumber := g.cgbBCPS / 8
					colorNumber := g.cgbBCPS / 2 % 4
					firstByte := g.cgbBCPS%2 == 0

					cr, cg, cb, _ := g.CgbBgPaletteColors[paletteNumber][colorNumber].RGBA()
					if firstByte {
						return ((uint8(cg>>3) & 0x7) << 5) | (uint8(cr>>3) & 0x1f)
					}
					return (uint8(cg>>6) & 0x3) | ((uint8(cb>>3) & 0x1f) << 2)
				}
				return 0x00
			case 0xFF6A:
				if g.gb.cgbModeEnabled {
					if g.cgbOCPSAutoincrement {
						return (0x1 << 7) | g.cgbOCPS
					}
					return g.cgbOCPS
				}
				return 0x00
			case 0xFF6B:
				if g.gb.cgbModeEnabled {
					paletteNumber := g.cgbOCPS / 8
					colorNumber := g.cgbOCPS / 2 % 4
					firstByte := g.cgbOCPS%2 == 0

					cr, cg, cb, _ := g.CgbObjPaletteColors[paletteNumber][colorNumber].RGBA()
					if firstByte {
						return ((uint8(cg>>3) & 0x7) << 5) | (uint8(cr>>3) & 0x1f)
					}
					return (uint8(cg>>6) & 0x3) | ((uint8(cb>>3) & 0x1f) << 2)
				}
				return 0x00
			default:
				return 0x00
			}
		default:
			return 0
		}
	default:
		return 0
	}
}

func (g *Gpu) WriteByte(address uint16, value uint8) {
	switch address & 0xF000 {
	case 0x8000, 0x9000: // VRAM
		g.vram[address&0x1FFF+uint16(g.vramBank)*0x2000] = value
		g.updateTile(address)
		g.updateTileAttribute(address)
	case 0xF000:
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
				g.StatEnableMode0 = value&0x8 != 0
				g.StatEnableMode1 = value&0x10 != 0
				g.StatEnableMode2 = value&0x20 != 0
				g.StatEnableLYC = value&0x40 != 0
			case 0xFF42:
				g.scrollY = value
			case 0xFF43:
				g.scrollX = value
			case 0xFF45:
				g.ScanlineCompare = value
			case 0xFF46:
				g.oamDMA(value)
			case 0xFF47:
				g.BgPaletteMap[0] = value & 0x3
				g.BgPaletteMap[1] = (value >> 2) & 0x3
				g.BgPaletteMap[2] = (value >> 4) & 0x3
				g.BgPaletteMap[3] = (value >> 6) & 0x3
			case 0xFF48:
				g.SpritePaletteMap[0][0] = value & 0x3
				g.SpritePaletteMap[0][1] = (value >> 2) & 0x3
				g.SpritePaletteMap[0][2] = (value >> 4) & 0x3
				g.SpritePaletteMap[0][3] = (value >> 6) & 0x3
			case 0xFF49:
				g.SpritePaletteMap[1][0] = value & 0x3
				g.SpritePaletteMap[1][1] = (value >> 2) & 0x3
				g.SpritePaletteMap[1][2] = (value >> 4) & 0x3
				g.SpritePaletteMap[1][3] = (value >> 6) & 0x3
			case 0xFF4A:
				g.windowY = value
			case 0xFF4B:
				g.windowX = value
			case 0xFF4F:
				if g.gb.cgbModeEnabled {
					g.vramBank = int(value & 0x1)
				}
			case 0xFF68:
				if g.gb.cgbModeEnabled {
					g.cgbBCPS = value & 0x3F
					g.cgbBCPSAutoincrement = value&0b10000000 > 0
				}
			case 0xFF69:
				if g.gb.cgbModeEnabled {
					paletteNumber := g.cgbBCPS / 8
					colorNumber := g.cgbBCPS / 2 % 4
					firstByte := g.cgbBCPS%2 == 0

					cr, cg, cb, ca := g.CgbBgPaletteColors[paletteNumber][colorNumber].RGBA()
					if firstByte {
						cr = uint32(value & 0x1F << 3)
						cg = uint32(uint8(cg)&0b11000000 | ((value >> 5) << 3))
					} else {
						cg = uint32(uint8(cg)&0b00111000 | ((value & 0x3) << 6))
						cb = uint32(((value >> 2) & 0x1F) << 3)
					}

					g.CgbBgPaletteColors[paletteNumber][colorNumber] = color.RGBA{uint8(cr), uint8(cg), uint8(cb), uint8(ca)}

					if g.cgbBCPSAutoincrement {
						g.cgbBCPS = (g.cgbBCPS + 1) & 0x3f
					}
				}
			case 0xFF6A:
				if g.gb.cgbModeEnabled {
					g.cgbOCPS = value & 0x3F
					g.cgbOCPSAutoincrement = value&0b10000000 > 0
				}
			case 0xFF6B:
				if g.gb.cgbModeEnabled {
					paletteNumber := g.cgbOCPS / 8
					colorNumber := g.cgbOCPS / 2 % 4
					firstByte := g.cgbOCPS%2 == 0

					cr, cg, cb, ca := g.CgbObjPaletteColors[paletteNumber][colorNumber].RGBA()
					if firstByte {
						cr = uint32(value & 0x1F << 3)
						cg = uint32(uint8(cg)&0b11000000 | ((value >> 5) << 3))
					} else {
						cg = uint32(uint8(cg)&0b00111000 | ((value & 0x3) << 6))
						cb = uint32(((value >> 2) & 0x1F) << 3)
					}

					g.CgbObjPaletteColors[paletteNumber][colorNumber] = color.RGBA{uint8(cr), uint8(cg), uint8(cb), uint8(ca)}

					if g.cgbOCPSAutoincrement {
						g.cgbOCPS = (g.cgbOCPS + 1) & 0x3f
					}
				}
			default:
			}
		default:
		}
	default:
	}
}

func (g *Gpu) updateTile(address uint16) {
	if address < 0x8000 || address >= 0x9800 {
		return
	}

	address &= 0x1FFE
	tileIndex := int(address>>4) + 384*g.vramBank
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

func (g *Gpu) updateTileAttribute(address uint16) {
	if address >= 0x9800 && address < 0xA000 {
		tileId := address - 0x9800
		ram := g.vram[tileId+0x3800]
		g.TileAttributes[tileId].useBgPriorityInsteadOfOam = (ram>>7)&0x1 > 0x0
		g.TileAttributes[tileId].yFlip = (ram>>6)&0x1 > 0x0
		g.TileAttributes[tileId].xFlip = (ram>>5)&0x1 > 0x0
		g.TileAttributes[tileId].tileVramBank = int((ram >> 3) & 0x1)
		g.TileAttributes[tileId].bgPaletteNumber = int(ram & 0x7)
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
	if objectId < uint16(len(g.SpriteObjectData)) {
		switch address & 0x3 {
		case 0:
			g.SpriteObjectData[objectId].y = int(value) - 16
		case 1:
			g.SpriteObjectData[objectId].x = int(value) - 8
		case 2:
			g.SpriteObjectData[objectId].tile = value
		case 3:
			// log.Printf("Updated tile %d: %v, Change was at address 0x%04x with value 0x%02x 0b%08b", objectId, g.spriteObjectData[objectId], address&0x3, value, value)
			g.SpriteObjectData[objectId].useSecondPalette = false
			g.SpriteObjectData[objectId].xflip = false
			g.SpriteObjectData[objectId].yflip = false
			g.SpriteObjectData[objectId].priority = false
			if value&0x10 > 0 {
				g.SpriteObjectData[objectId].useSecondPalette = true
			}
			if value&0x20 > 0 {
				g.SpriteObjectData[objectId].xflip = true
			}
			if value&0x40 > 0 {
				g.SpriteObjectData[objectId].yflip = true
			}
			if value&0x80 == 0 {
				g.SpriteObjectData[objectId].priority = true
			}
			if g.gb.cgbModeEnabled {
				g.SpriteObjectData[objectId].cgbPalette = int(value & 0x7)
				g.SpriteObjectData[objectId].vramBank = int((value >> 3) & 0x1)
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
			g.SetScanline(gb, g.CurrentScanline+1)

			if g.CurrentScanline > uint8(ScreenHeight)-1 {
				g.SetLCDMode(gb, 1)
				gb.ReadyToRender = gb.WorkingScreen
				gb.WorkingScreen = [ScreenWidth][ScreenHeight][3]uint8{}
				gb.Memory.GetInterruptFlags().TriggeredFlags |= 0x01
			} else {
				g.SetLCDMode(gb, 2)
			}
		}
	case 1: // VBlank
		if g.modeClock >= 456 {
			g.modeClock = 0
			g.SetScanline(gb, g.CurrentScanline+1)

			if g.CurrentScanline > uint8(ScreenHeight)+10-1 {
				g.SetLCDMode(gb, 2)
				g.SetScanline(gb, 0)
			}
		}
	case 2:
		if g.modeClock >= 80 {
			g.modeClock = 0
			g.SetLCDMode(gb, 3)
		}
	case 3:
		if g.modeClock >= 172 {
			g.modeClock = 0
			g.SetLCDMode(gb, 0)

			g.RenderScanline(gb)
		}
	}
	g.HandleStatInterrupt()
}

func (g *Gpu) SetLCDMode(gb *Gameboy, value uint8) {
	g.currentMode = value
	g.StatTriggerMode0 = value == 0
	g.StatTriggerMode1 = value == 1
	g.StatTriggerMode2 = value == 2
}

func (g *Gpu) SetScanline(gb *Gameboy, value uint8) {
	g.CurrentScanline = value
	g.StatTriggerLYC = g.CurrentScanline == g.ScanlineCompare
}

func (g *Gpu) RenderScanline(gb *Gameboy) {
	var scanrow [ScreenWidth]byte
	if g.backgroundActivated || g.gb.cgbModeEnabled {
		scanrow = g.renderTiles(gb)
	}

	if (g.backgroundActivated || g.gb.cgbModeEnabled) && g.windowActivated {
		scanrow = g.renderWindow(gb, scanrow)
	}

	if g.spritesActivated {
		g.renderSprites(gb, scanrow)
	}
}

func (g *Gpu) renderTiles(gb *Gameboy) (scanrow [ScreenWidth]byte) {
	var mapOffset uint16
	var xPos, yPos uint8

	mapOffset = 0x1800
	if g.backgroundMap {
		mapOffset = 0x1C00
	}

	yPos = g.CurrentScanline + g.scrollY

	tileYIndex := uint16(yPos/8) * 32

	for i := uint8(0); int(i) < ScreenWidth; i++ {
		xPos = i + g.scrollX
		tileXIndex := uint16(xPos / 8)

		tileId := uint16(g.vram[mapOffset+tileYIndex+tileXIndex])
		if g.backgroundTile && tileId < 128 {
			tileId += 256
		}
		tileAttr := g.TileAttributes[mapOffset-0x1800+tileYIndex+tileXIndex]
		tileId += uint16(384 * tileAttr.tileVramBank)

		xPixelPos := xPos % 8
		yPixelPos := yPos % 8

		if g.gb.cgbModeEnabled && tileAttr.yFlip {
			yPixelPos = 7 - yPixelPos
		}
		if g.gb.cgbModeEnabled && tileAttr.xFlip {
			xPixelPos = 7 - xPixelPos
		}

		pixelPaletteColor := g.tileSet[tileId][yPixelPos][xPixelPos]
		scanrow[i] = pixelPaletteColor
		pixelRealColor := g.BgPaletteMap[pixelPaletteColor]
		red, green, blue, _ := g.bgPaletteColors[pixelRealColor].RGBA()
		if g.gb.cgbModeEnabled {
			palette := tileAttr.bgPaletteNumber
			red, green, blue, _ = g.CgbBgPaletteColors[palette][pixelPaletteColor].RGBA()
		}
		gb.WorkingScreen[i][g.CurrentScanline][0] = uint8(red)
		gb.WorkingScreen[i][g.CurrentScanline][1] = uint8(green)
		gb.WorkingScreen[i][g.CurrentScanline][2] = uint8(blue)
	}

	return scanrow
}

func (g *Gpu) renderWindow(gb *Gameboy, scanrow [ScreenWidth]byte) [ScreenWidth]byte {
	var mapOffset uint16
	var xPos, yPos uint8

	if g.CurrentScanline < g.windowY {
		return scanrow
	}

	mapOffset = 0x1800
	if g.windowMap {
		mapOffset = 0x1C00
	}

	yPos = g.CurrentScanline - g.windowY

	tileYIndex := uint16(yPos/8) * 32

	for i := uint8(0); int(i) < ScreenWidth; i++ {
		xPos = i + g.windowX - 7
		tileXIndex := uint16(xPos / 8)

		tileId := uint16(g.vram[mapOffset+tileYIndex+tileXIndex])
		if g.backgroundTile && tileId < 128 {
			tileId += 256
		}
		tileAttr := g.TileAttributes[mapOffset-0x1800+tileYIndex+tileXIndex]
		tileId += uint16(384 * tileAttr.tileVramBank)

		xPixelPos := xPos % 8
		yPixelPos := yPos % 8

		if g.gb.cgbModeEnabled && tileAttr.yFlip {
			yPixelPos = 7 - yPixelPos
		}
		if g.gb.cgbModeEnabled && tileAttr.xFlip {
			xPixelPos = 7 - xPixelPos
		}

		pixelPaletteColor := g.tileSet[tileId][yPixelPos][xPixelPos]
		scanrow[i] = pixelPaletteColor
		pixelRealColor := g.BgPaletteMap[pixelPaletteColor]
		red, green, blue, _ := g.bgPaletteColors[pixelRealColor].RGBA()
		if g.gb.cgbModeEnabled {
			palette := tileAttr.bgPaletteNumber
			red, green, blue, _ = g.CgbBgPaletteColors[palette][pixelPaletteColor].RGBA()
		}
		gb.WorkingScreen[i][g.CurrentScanline][0] = uint8(red)
		gb.WorkingScreen[i][g.CurrentScanline][1] = uint8(green)
		gb.WorkingScreen[i][g.CurrentScanline][2] = uint8(blue)
	}

	return scanrow
}

func (g *Gpu) GetTilemapAsBytearray(vramBank int) []byte {
	var frame []byte = make([]byte, 4*ScreenHeight*ScreenWidth)

	// palette := [4]color.Color{
	// 	color.RGBA{0x00, 0x00, 0x00, 255},
	// 	color.RGBA{0xFF, 0x00, 0x00, 255},
	// 	color.RGBA{0x00, 0xFF, 0x00, 255},
	// 	color.RGBA{0x00, 0x00, 0xFF, 255},
	// }

	start := 384 * vramBank
	for tileId := start; tileId < len(g.tileSet); tileId++ {
		xBase := (tileId - start) % 16
		yBase := (tileId - start) / 16
		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				pixelPaletteColor := g.tileSet[tileId][y][x]
				red, green, blue, _ := g.bgPaletteColors[pixelPaletteColor].RGBA()
				pixelPos := (yBase*8+y)*ScreenWidth + (xBase*8 + x)
				if 4*pixelPos+3 < len(frame) {
					frame[4*pixelPos] = byte(red)
					frame[4*pixelPos+1] = byte(green)
					frame[4*pixelPos+2] = byte(blue)
					frame[4*pixelPos+3] = 0xFF
				}
			}
		}
	}

	return frame
}

func (g *Gpu) renderSprites(gb *Gameboy, scanrow [ScreenWidth]byte) {
	spriteCount := 0
	for i := 0; i < len(g.SpriteObjectData); i++ {
		var ySize int = 8
		spriteObject := g.SpriteObjectData[i]
		if g.bigSpritesActivated {
			ySize = 16
		}

		if spriteObject.y > int(g.CurrentScanline) || spriteObject.y+ySize <= int(g.CurrentScanline) {
			continue
		}
		if spriteCount >= 10 {
			break
		}
		spriteCount++

		palette := g.SpritePaletteMap[0]
		paletteColors := g.spritePaletteColors[0]
		if spriteObject.useSecondPalette {
			palette = g.SpritePaletteMap[1]
			paletteColors = g.spritePaletteColors[1]
		}

		tilerowIndex := g.CurrentScanline - uint8(spriteObject.y)
		if spriteObject.yflip {
			tilerowIndex = uint8(ySize) - tilerowIndex - 1
		}
		tilerowIndex = tilerowIndex % 8
		tileId := uint16(spriteObject.tile)
		if g.bigSpritesActivated {
			if g.CurrentScanline-uint8(spriteObject.y) < 8 {
				if spriteObject.yflip {
					tileId |= 0x01
				} else {
					tileId &= 0xFE
				}
			} else {
				if spriteObject.yflip {
					tileId &= 0xFE
				} else {
					tileId |= 0x01
				}
			}
		}
		if g.gb.cgbModeEnabled {
			tileId += uint16(384 * spriteObject.vramBank)
		}
		tilerow := g.tileSet[tileId][tilerowIndex]

		for x := 0; x < 8; x++ {
			pixelPos := spriteObject.x + x

			if pixelPos < 0 || pixelPos >= ScreenWidth {
				continue
			}

			pixelPaletteColor := tilerow[x]
			if spriteObject.xflip {
				pixelPaletteColor = tilerow[7-x]
			}
			if pixelPaletteColor > 0 && (spriteObject.priority || scanrow[pixelPos] == 0) {
				pixelRealColor := palette[pixelPaletteColor]
				red, green, blue, _ := paletteColors[pixelRealColor].RGBA()
				if g.gb.cgbModeEnabled {
					palette := spriteObject.cgbPalette
					red, green, blue, _ = g.CgbObjPaletteColors[palette][pixelPaletteColor].RGBA()
				}
				gb.WorkingScreen[spriteObject.x+x][g.CurrentScanline][0] = uint8(red)
				gb.WorkingScreen[spriteObject.x+x][g.CurrentScanline][1] = uint8(green)
				gb.WorkingScreen[spriteObject.x+x][g.CurrentScanline][2] = uint8(blue)
			}
		}
	}
}

func (g *Gpu) Reset(gb *Gameboy) {
	g.CurrentScanline = 0
	g.currentMode = 0
	gb.clearScreen()

	g.tileSet = [768][8][8]uint8{}
	g.oam = [0xA0]uint8{}
	g.SpriteObjectData = [40]SpriteObject{}

	for i := 0; i < len(g.SpriteObjectData); i++ {
		g.SpriteObjectData[i] = SpriteObject{
			x: -8,
			y: -16,
		}
	}
}

func (gb *Gameboy) clearScreen() {
	r, g, b, _ := gb.Gpu.bgPaletteColors[0].RGBA()
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			gb.WorkingScreen[x][y][0] = uint8(r)
			gb.WorkingScreen[x][y][1] = uint8(g)
			gb.WorkingScreen[x][y][2] = uint8(b)
		}
	}
	gb.ReadyToRender = gb.WorkingScreen
	gb.WorkingScreen = [ScreenWidth][ScreenHeight][3]uint8{}
}

func (g *Gpu) HandleStatInterrupt() {
	LYCInterrupt := g.StatEnableLYC && g.StatTriggerLYC
	Mode0Interrupt := g.StatEnableMode0 && g.StatTriggerMode0
	Mode1Interrupt := g.StatEnableMode1 && g.StatTriggerMode1
	Mode2Interrupt := g.StatEnableMode2 && g.StatTriggerMode2

	StatInterruptTrigger := LYCInterrupt || Mode0Interrupt || Mode1Interrupt || Mode2Interrupt

	if triggered := g.detectRisingEdge(StatInterruptTrigger); triggered {
		g.gb.Memory.GetInterruptFlags().TriggeredFlags |= 0x2
	}
}

func (g *Gpu) detectRisingEdge(signal bool) bool {
	result := signal && !g.StatInterruptDelay
	g.StatInterruptDelay = signal
	return result
}

func (gb *Gameboy) GetReadyFramebufferAsBytearray() []byte {
	var frame []byte = make([]byte, 0, 4*ScreenHeight*ScreenWidth)
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			frame = append(frame, gb.ReadyToRender[x][y][0])
			frame = append(frame, gb.ReadyToRender[x][y][1])
			frame = append(frame, gb.ReadyToRender[x][y][2])
			frame = append(frame, 0xFF)
		}
	}

	return frame
}
