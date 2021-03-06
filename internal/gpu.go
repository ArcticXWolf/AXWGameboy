package internal

import (
	"fmt"
	"image/color"

	"go.janniklasrichter.de/axwgameboy/internal/cartridge"
)

type SpriteObject struct {
	X                int
	Y                int
	Tile             uint8
	UseSecondPalette bool
	Xflip            bool
	Yflip            bool
	Priority         bool
	VramBank         int
	CgbPalette       int
}

func (s *SpriteObject) String() string {
	palette := 0
	if s.UseSecondPalette {
		palette = 1
	}
	return fmt.Sprintf("Sprite X%03d Y%03d T%03d V%01d D%01d P%01d", s.X, s.Y, s.Tile, s.VramBank, palette, s.CgbPalette)
}

type TileAttributes struct {
	UseBgPriorityInsteadOfOam bool
	yFlip                     bool
	xFlip                     bool
	TileVramBank              int
	BgPaletteNumber           int
}

type Gpu struct {
	gb *Gameboy

	backgroundActivated bool
	spritesActivated    bool
	bigSpritesActivated bool
	BackgroundMap       bool
	BackgroundTile      bool
	windowActivated     bool
	WindowMap           bool
	lcdActivated        bool
	lcdCleared          bool

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

	ScrollX                    uint8
	ScrollY                    uint8
	WindowX                    uint8
	WindowY                    uint8
	WindowYInternalLineCounter uint8
	ScanlineCompare            uint8
	CurrentScanline            uint8
	modeClock                  int

	VramBank int
	Vram     [0x4000]byte

	Oam              [0xA0]byte
	SpriteObjectData [40]SpriteObject

	TileSet        [768][8][8]uint8
	TileAttributes [0x800]TileAttributes
	TileBGPriority [ScreenWidth][ScreenHeight]bool

	BgPaletteMap        [4]uint8
	BgPaletteColors     [4]color.Color
	CgbBgPaletteColors  [8][4]color.Color
	SpritePaletteMap    [2][4]uint8
	SpritePaletteColors [2][4]color.Color

	CgbObjPaletteColors  [8][4]color.Color
	cgbBCPS              uint8
	cgbBCPSAutoincrement bool
	cgbOCPS              uint8
	cgbOCPSAutoincrement bool

	cgbDMASource      uint16
	cgbDMADestination uint16
	cgbDMALength      byte
	cgbHDMAActive     bool
	cgbHDMAAborted    bool
}

func NewGpu(gb *Gameboy) *Gpu {

	g := &Gpu{
		gb:              gb,
		BgPaletteColors: getPaletteColorsByName(gb.Options.Palette),
		SpritePaletteColors: [2][4]color.Color{
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
			X: -8,
			Y: -16,
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
		return g.Vram[address&0x1FFF+uint16(g.VramBank)*0x2000]
	case 0xF000:
		switch address & 0xFF00 {
		case 0xFE00:
			if address < 0xFEA0 {
				return g.Oam[address&0xFF]
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
				if g.BackgroundMap {
					bM = 0x08
				}
				if !g.BackgroundTile {
					bT = 0x10
				}
				if g.windowActivated {
					wA = 0x20
				}
				if g.WindowMap {
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
				return g.ScrollY
			case 0xFF43:
				return g.ScrollX
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
				return g.WindowY
			case 0xFF4B:
				return g.WindowX
			case 0xFF4F:
				return uint8(g.VramBank)
			case 0xFF51:
				return 0xFF
			case 0xFF52:
				return 0xFF
			case 0xFF53:
				return 0xFF
			case 0xFF54:
				return 0xFF
			case 0xFF55:
				if g.gb.CgbModeEnabled {
					if g.cgbHDMAActive {
						return g.cgbDMALength
					} else if g.cgbHDMAAborted {
						return g.cgbDMALength | 0x80
					}
				}
				return 0xFF
			case 0xFF68:
				if g.gb.CgbModeEnabled {
					if g.cgbBCPSAutoincrement {
						return (0x1 << 7) | g.cgbBCPS
					}
					return g.cgbBCPS
				}
				return 0x00
			case 0xFF69:
				if g.gb.CgbModeEnabled {
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
				if g.gb.CgbModeEnabled {
					if g.cgbOCPSAutoincrement {
						return (0x1 << 7) | g.cgbOCPS
					}
					return g.cgbOCPS
				}
				return 0x00
			case 0xFF6B:
				if g.gb.CgbModeEnabled {
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
		g.Vram[address&0x1FFF+uint16(g.VramBank)*0x2000] = value
		g.updateTile(address)
		g.updateTileAttribute(address)
	case 0xF000:
		switch address & 0xFF00 {
		case 0xFE00:
			if address < 0xFEA0 {
				g.Oam[address&0xFF] = value
				g.updateSpriteObject(address, value)
			}
		case 0xFF00:
			switch address {
			case 0xFF40:
				g.gb.RingLogger.Printf("gpu", "Register Write %5s $%08b to $%08b on line %d", "LDLC", g.ReadByte(address), value, g.CurrentScanline)
				g.backgroundActivated = value&0x01 != 0
				g.spritesActivated = value&0x02 != 0
				g.bigSpritesActivated = value&0x04 != 0
				g.BackgroundMap = value&0x08 != 0
				g.BackgroundTile = value&0x10 == 0
				g.windowActivated = value&0x20 != 0
				g.WindowMap = value&0x40 != 0
				g.lcdActivated = value&0x80 != 0
				g.lcdCleared = false
			case 0xFF41:
				g.StatEnableMode0 = value&0x8 != 0
				g.StatEnableMode1 = value&0x10 != 0
				g.StatEnableMode2 = value&0x20 != 0
				g.StatEnableLYC = value&0x40 != 0
			case 0xFF42:
				g.gb.RingLogger.Printf("gpu", "Register Write %5s $%02x to $%02x on line %d", "SCY", g.ReadByte(address), value, g.CurrentScanline)
				g.ScrollY = value
			case 0xFF43:
				g.gb.RingLogger.Printf("gpu", "Register Write %5s $%02x to $%02x on line %d", "SCX", g.ReadByte(address), value, g.CurrentScanline)
				g.ScrollX = value
			case 0xFF45:
				g.gb.RingLogger.Printf("gpu", "Register Write %5s $%02d to $%02d on line %d", "LYC", g.ReadByte(address), value, g.CurrentScanline)
				g.ScanlineCompare = value
				g.StatTriggerLYC = g.CurrentScanline == g.ScanlineCompare
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
				g.gb.RingLogger.Printf("gpu", "Register Write %5s $%03d to $%03d on line %d", "WY", g.ReadByte(address), value, g.CurrentScanline)
				g.WindowY = value
			case 0xFF4B:
				g.gb.RingLogger.Printf("gpu", "Register Write %5s $%03d to $%03d on line %d", "WX", g.ReadByte(address), value, g.CurrentScanline)
				g.WindowX = value
			case 0xFF4F:
				if g.gb.CgbModeEnabled {
					g.VramBank = int(value & 0x1)
				}
			case 0xFF51:
				if g.gb.CgbModeEnabled {
					g.cgbDMASource = (g.cgbDMASource & 0x00FF) | (uint16(value) << 8)

				}
			case 0xFF52:
				if g.gb.CgbModeEnabled {
					g.cgbDMASource = (g.cgbDMASource & 0xFF00) | uint16(value&0xF0)
				}
			case 0xFF53:
				if g.gb.CgbModeEnabled {
					g.cgbDMADestination = (g.cgbDMADestination & 0x00FF) | (uint16(value&0x1F) << 8)
				}
			case 0xFF54:
				if g.gb.CgbModeEnabled {
					g.cgbDMADestination = (g.cgbDMADestination & 0xFF00) | uint16(value&0xF0)
				}
			case 0xFF55:
				if g.gb.CgbModeEnabled {
					g.startDMA(value)
				}
			case 0xFF68:
				if g.gb.CgbModeEnabled {
					g.cgbBCPS = value & 0x3F
					g.cgbBCPSAutoincrement = value&0b10000000 > 0
				}
			case 0xFF69:
				if g.gb.CgbModeEnabled {
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
				if g.gb.CgbModeEnabled {
					g.cgbOCPS = value & 0x3F
					g.cgbOCPSAutoincrement = value&0b10000000 > 0
				}
			case 0xFF6B:
				if g.gb.CgbModeEnabled {
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
	tileIndex := int(address>>4) + 384*g.VramBank
	y := (address >> 1) & 0x7

	for x := 0; x < 8; x++ {
		bitIndex := 1 << (7 - x)
		lowerBit := 0
		higherBit := 0

		if g.Vram[address+uint16(g.VramBank)*0x2000]&byte(bitIndex) != 0 {
			lowerBit = 1
		}
		if g.Vram[address+uint16(g.VramBank)*0x2000+1]&byte(bitIndex) != 0 {
			higherBit = 2
		}

		g.TileSet[tileIndex][y][x] = uint8(lowerBit) + uint8(higherBit)
	}
}

func (g *Gpu) updateTileAttribute(address uint16) {
	if g.VramBank != 1 || address < 0x9800 || address >= 0xA000 {
		return
	}
	tileId := address - 0x9800
	ram := g.Vram[tileId+0x3800]
	g.TileAttributes[tileId].UseBgPriorityInsteadOfOam = (ram>>7)&0x1 > 0x0
	g.TileAttributes[tileId].yFlip = (ram>>6)&0x1 > 0x0
	g.TileAttributes[tileId].xFlip = (ram>>5)&0x1 > 0x0
	g.TileAttributes[tileId].TileVramBank = int((ram >> 3) & 0x1)
	g.TileAttributes[tileId].BgPaletteNumber = int(ram & 0x7)
}

func (g *Gpu) ResetTileAttributes() {
	for i := 0; i < len(g.TileAttributes); i++ {
		g.Vram[i+0x3800] = 0x00
		g.TileAttributes[i].BgPaletteNumber = 0
		g.TileAttributes[i].TileVramBank = 0
		g.TileAttributes[i].xFlip = false
		g.TileAttributes[i].yFlip = false
		g.TileAttributes[i].UseBgPriorityInsteadOfOam = false
	}
}

func (g *Gpu) oamDMA(value uint8) {
	var x uint16
	for x = 0; x < 0xA0; x++ {
		g.Oam[x] = g.gb.Memory.ReadByte((uint16(value) << 8) + x)
		g.updateSpriteObject(0xFE00+x, g.Oam[x])
	}
	g.gb.Cpu.ClockCycles += 164
}

func (g *Gpu) startDMA(value uint8) {
	if g.cgbHDMAActive && (value>>7) == 0 {
		g.cgbHDMAActive = false
		g.cgbHDMAAborted = true
		return
	}

	bytesToCopy := (uint16(value&0x7F) + 1) * 0x10
	g.cgbHDMAAborted = false

	if (value >> 7) == 0 {
		g.performDMA(bytesToCopy)
		return
	}

	g.cgbDMALength = byte(value & 0x7F)
	g.cgbHDMAActive = true
}

func (g *Gpu) stepHDMA() {
	if !g.cgbHDMAActive {
		return
	}

	g.performDMA(0x10)

	if g.cgbDMALength > 0 {
		g.cgbDMALength--
		return
	}

	g.cgbHDMAActive = false
}

func (g *Gpu) performDMA(length uint16) {
	var x uint16
	for x = 0; x < length; x++ {
		address := 0x8000 + g.cgbDMADestination + x
		value := g.gb.Memory.ReadByte(g.cgbDMASource + x)
		g.WriteByte(address, value)
		g.updateTile(address)
		g.updateSpriteObject(address, value)
		g.updateTileAttribute(address)
	}

	g.cgbDMASource += length
	g.cgbDMADestination += length
}

func (g *Gpu) updateSpriteObject(address uint16, value uint8) {
	if address < 0xFE00 || address >= 0xFEA0 {
		return
	}
	address = address - 0xFE00

	objectId := address >> 2
	if objectId < uint16(len(g.SpriteObjectData)) {
		switch address & 0x3 {
		case 0:
			g.SpriteObjectData[objectId].Y = int(value) - 16
		case 1:
			g.SpriteObjectData[objectId].X = int(value) - 8
		case 2:
			g.SpriteObjectData[objectId].Tile = value
		case 3:
			g.SpriteObjectData[objectId].UseSecondPalette = false
			g.SpriteObjectData[objectId].Xflip = false
			g.SpriteObjectData[objectId].Yflip = false
			g.SpriteObjectData[objectId].Priority = false
			if value&0x10 > 0 {
				g.SpriteObjectData[objectId].UseSecondPalette = true
			}
			if value&0x20 > 0 {
				g.SpriteObjectData[objectId].Xflip = true
			}
			if value&0x40 > 0 {
				g.SpriteObjectData[objectId].Yflip = true
			}
			if value&0x80 == 0 {
				g.SpriteObjectData[objectId].Priority = true
			}
			if g.gb.CgbModeEnabled {
				g.SpriteObjectData[objectId].CgbPalette = int(value & 0x7)
				g.SpriteObjectData[objectId].VramBank = int((value >> 3) & 0x1)
			}
		}
	}
}

func (g *Gpu) ResetOAM() {
	var x uint16
	for x = 0; x < 0xA0; x++ {
		g.Oam[x] = 0x00
		g.updateSpriteObject(x, g.Oam[x])
	}
}

func (g *Gpu) Update(gb *Gameboy, cyclesUsed int) {
	if !g.lcdActivated {
		g.gb.RingLogger.Printf("gpu", "Clear LCD on line %d", g.CurrentScanline)
		g.CurrentScanline = 0
		g.currentMode = 0
		g.WindowYInternalLineCounter = 0
		if !g.lcdCleared {
			gb.clearScreen()
			g.lcdCleared = true
		}
		return
	}

	g.modeClock += cyclesUsed / int(gb.GetSpeedMultiplier())

	switch g.currentMode {
	case 0: // HBlank
		if g.modeClock >= 204 {
			g.modeClock = 0
			g.SetScanline(gb, g.CurrentScanline+1)
			g.stepHDMA()

			if g.CurrentScanline > uint8(ScreenHeight)-1 {
				g.SetLCDMode(gb, 1)
				gb.ReadyToRender = gb.WorkingScreen
				gb.WorkingScreen = [ScreenWidth][ScreenHeight][3]uint8{}
				g.WindowYInternalLineCounter = 0
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
	gb.Gpu.HandleStatInterrupt()
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
	if g.backgroundActivated || g.gb.CgbModeEnabled {
		scanrow = g.renderTiles(gb)
	} else {
		scanrow = g.renderEmptyLine(gb)
	}

	if (g.backgroundActivated || g.gb.CgbModeEnabled) && g.windowActivated {
		scanrow = g.renderWindow(gb, scanrow)
	}

	if g.spritesActivated {
		g.renderSprites(gb, scanrow)
	}
}

func (g *Gpu) renderEmptyLine(gb *Gameboy) (scanrow [ScreenWidth]byte) {
	pixelRealColor := g.BgPaletteMap[0]
	red, green, blue, _ := g.BgPaletteColors[pixelRealColor].RGBA()
	for i := uint8(0); int(i) < ScreenWidth; i++ {
		scanrow[i] = 0
		gb.WorkingScreen[i][g.CurrentScanline][0] = uint8(red)
		gb.WorkingScreen[i][g.CurrentScanline][1] = uint8(green)
		gb.WorkingScreen[i][g.CurrentScanline][2] = uint8(blue)
	}
	return scanrow
}

func (g *Gpu) renderTiles(gb *Gameboy) (scanrow [ScreenWidth]byte) {
	var mapOffset uint16
	var xPos, yPos uint8

	mapOffset = 0x1800
	if g.BackgroundMap {
		mapOffset = 0x1C00
	}

	yPos = g.CurrentScanline + g.ScrollY

	tileYIndex := uint16(yPos/8) * 32

	for i := uint8(0); int(i) < ScreenWidth; i++ {
		xPos = i + g.ScrollX
		tileXIndex := uint16(xPos / 8)

		tileId := uint16(g.Vram[mapOffset+tileYIndex+tileXIndex])
		if g.BackgroundTile && tileId < 128 {
			tileId += 256
		}
		tileAttr := g.TileAttributes[mapOffset-0x1800+tileYIndex+tileXIndex]
		tileId += uint16(384 * tileAttr.TileVramBank)

		xPixelPos := xPos % 8
		yPixelPos := yPos % 8

		if g.gb.CgbModeEnabled && tileAttr.yFlip {
			yPixelPos = 7 - yPixelPos
		}
		if g.gb.CgbModeEnabled && tileAttr.xFlip {
			xPixelPos = 7 - xPixelPos
		}

		pixelPaletteColor := g.TileSet[tileId][yPixelPos][xPixelPos]
		scanrow[i] = pixelPaletteColor
		pixelRealColor := g.BgPaletteMap[pixelPaletteColor]
		red, green, blue, _ := g.BgPaletteColors[pixelRealColor].RGBA()
		if g.gb.CgbModeEnabled {
			palette := tileAttr.BgPaletteNumber
			red, green, blue, _ = g.CgbBgPaletteColors[palette][pixelPaletteColor].RGBA()
			g.TileBGPriority[i][g.CurrentScanline] = tileAttr.UseBgPriorityInsteadOfOam
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

	if g.CurrentScanline < g.WindowY {
		return scanrow
	} else if g.WindowX > uint8(ScreenWidth) {
		return scanrow
	} else if g.WindowY > uint8(ScreenHeight) {
		return scanrow
	}

	mapOffset = 0x1800
	if g.WindowMap {
		mapOffset = 0x1C00
	}

	yPos = g.WindowYInternalLineCounter

	tileYIndex := uint16((yPos)/8) * 32

	for i := uint8(0); int(i) < ScreenWidth; i++ {
		if i < g.WindowX-7 {
			continue
		}

		xPos = i - (g.WindowX - 7)
		tileXIndex := uint16(xPos / 8)

		tileId := uint16(g.Vram[mapOffset+tileYIndex+tileXIndex])
		if g.BackgroundTile && tileId < 128 {
			tileId += 256
		}
		tileAttr := g.TileAttributes[mapOffset-0x1800+tileYIndex+tileXIndex]
		tileId += uint16(384 * tileAttr.TileVramBank)

		xPixelPos := xPos % 8
		yPixelPos := yPos % 8

		if g.gb.CgbModeEnabled && tileAttr.yFlip {
			yPixelPos = 7 - yPixelPos
		}
		if g.gb.CgbModeEnabled && tileAttr.xFlip {
			xPixelPos = 7 - xPixelPos
		}

		pixelPaletteColor := g.TileSet[tileId][yPixelPos][xPixelPos]
		scanrow[i] = pixelPaletteColor
		pixelRealColor := g.BgPaletteMap[pixelPaletteColor]
		red, green, blue, _ := g.BgPaletteColors[pixelRealColor].RGBA()
		if g.gb.CgbModeEnabled {
			palette := tileAttr.BgPaletteNumber
			red, green, blue, _ = g.CgbBgPaletteColors[palette][pixelPaletteColor].RGBA()
			g.TileBGPriority[i][g.CurrentScanline] = tileAttr.UseBgPriorityInsteadOfOam
		}
		gb.WorkingScreen[i][g.CurrentScanline][0] = uint8(red)
		gb.WorkingScreen[i][g.CurrentScanline][1] = uint8(green)
		gb.WorkingScreen[i][g.CurrentScanline][2] = uint8(blue)
	}
	g.gb.RingLogger.Printf("gpu", "Rendering Window L%d on S%d", g.WindowYInternalLineCounter, g.CurrentScanline)
	g.WindowYInternalLineCounter++

	return scanrow
}

const spriteXPixelComparisonOffset int = 10

func (g *Gpu) renderSprites(gb *Gameboy, scanrow [ScreenWidth]byte) {
	var spriteXPerScreenPixel [ScreenWidth]int
	spriteCount := 0

	for i := 0; i < len(g.SpriteObjectData); i++ {
		var ySize int = 8
		spriteObject := g.SpriteObjectData[i]
		if g.bigSpritesActivated {
			ySize = 16
		}

		if spriteObject.Y > int(g.CurrentScanline) || spriteObject.Y+ySize <= int(g.CurrentScanline) {
			continue
		}
		if spriteCount >= 10 {
			break
		}
		spriteCount++

		palette := g.SpritePaletteMap[0]
		paletteColors := g.SpritePaletteColors[0]
		if spriteObject.UseSecondPalette {
			palette = g.SpritePaletteMap[1]
			paletteColors = g.SpritePaletteColors[1]
		}

		tilerowIndex := g.CurrentScanline - uint8(spriteObject.Y)
		if spriteObject.Yflip {
			tilerowIndex = uint8(ySize) - tilerowIndex - 1
		}
		tilerowIndex = tilerowIndex % 8
		tileId := uint16(spriteObject.Tile)
		if g.bigSpritesActivated {
			if g.CurrentScanline-uint8(spriteObject.Y) < 8 {
				if spriteObject.Yflip {
					tileId |= 0x01
				} else {
					tileId &= 0xFE
				}
			} else {
				if spriteObject.Yflip {
					tileId &= 0xFE
				} else {
					tileId |= 0x01
				}
			}
		}
		if g.gb.CgbModeEnabled {
			tileId += uint16(384 * spriteObject.VramBank)
		}
		tilerow := g.TileSet[tileId][tilerowIndex]

		for x := 0; x < 8; x++ {
			// skip pixels out of bounds
			pixelPos := spriteObject.X + x
			if pixelPos < 0 || pixelPos >= ScreenWidth {
				continue
			}

			pixelPaletteColor := tilerow[x]
			if spriteObject.Xflip {
				pixelPaletteColor = tilerow[7-x]
			}

			// skip transparent pixels
			if pixelPaletteColor == 0 {
				continue
			}

			// skip pixels without priority that are hidden by BG
			if !g.gb.CgbModeEnabled || g.backgroundActivated {
				if !(spriteObject.Priority && !g.TileBGPriority[pixelPos][g.CurrentScanline]) && scanrow[pixelPos] != 0 {
					continue
				}
			}

			if g.gb.CgbModeEnabled {
				// skip if pixel is occupied by sprite with lower spriteObjectID
				if spriteXPerScreenPixel[pixelPos] != 0 {
					continue
				}
			} else {
				// skip if pixel is occupied by sprite with lower x coordinate
				if spriteXPerScreenPixel[pixelPos] != 0 && spriteXPerScreenPixel[pixelPos] <= spriteObject.X+spriteXPixelComparisonOffset {
					continue
				}
			}

			pixelRealColor := palette[pixelPaletteColor]
			red, green, blue, _ := paletteColors[pixelRealColor].RGBA()
			if g.gb.CgbModeEnabled {
				cgbPalette := spriteObject.CgbPalette
				if g.gb.Memory.Cartridge.CartridgeHeader().CartridgeGBMode == cartridge.OnlyDMG {
					cgbPalette = 0
					if spriteObject.UseSecondPalette {
						cgbPalette = 1
					}
				}
				red, green, blue, _ = g.CgbObjPaletteColors[cgbPalette][pixelPaletteColor].RGBA()
			}
			gb.WorkingScreen[pixelPos][g.CurrentScanline][0] = uint8(red)
			gb.WorkingScreen[pixelPos][g.CurrentScanline][1] = uint8(green)
			gb.WorkingScreen[pixelPos][g.CurrentScanline][2] = uint8(blue)

			spriteXPerScreenPixel[pixelPos] = spriteObject.X + spriteXPixelComparisonOffset
		}
	}
}

func (g *Gpu) Reset(gb *Gameboy) {
	g.CurrentScanline = 0
	g.currentMode = 0
	gb.clearScreen()

	g.TileSet = [768][8][8]uint8{}
	g.Oam = [0xA0]uint8{}
	g.SpriteObjectData = [40]SpriteObject{}

	for i := 0; i < len(g.SpriteObjectData); i++ {
		g.SpriteObjectData[i] = SpriteObject{
			X: -8,
			Y: -16,
		}
	}
}

func (gb *Gameboy) clearScreen() {
	gb.WorkingScreen = [ScreenWidth][ScreenHeight][3]uint8{}
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
