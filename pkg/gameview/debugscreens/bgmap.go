package debugscreens

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type BgMap struct {
	X     int
	Y     int
	MapId int
}

func (t *BgMap) GetXPos() int   { return t.X }
func (t *BgMap) GetYPos() int   { return t.Y }
func (t *BgMap) GetWidth() int  { return 32 * 8 }
func (t *BgMap) GetHeight() int { return 32 * 8 }

func (t *BgMap) GetAsBytearray(gb *internal.Gameboy) []byte {
	var frame []byte = make([]byte, 4*t.GetWidth()*t.GetHeight())

	mapOffset := 0x1800
	if t.MapId == 1 {
		mapOffset = 0x1C00
	}

	for tileMapId := 0; tileMapId < 32*32; tileMapId++ {
		xBase := tileMapId % (t.GetWidth() / 8)
		yBase := tileMapId / (t.GetWidth() / 8)
		tileId := uint16(gb.Gpu.Vram[mapOffset+tileMapId])
		if gb.Gpu.BackgroundTile && tileId < 128 {
			tileId += 256
		}
		tileAttr := gb.Gpu.TileAttributes[mapOffset-0x1800+tileMapId]
		tileId += uint16(384 * tileAttr.TileVramBank)
		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				pixelPaletteColor := gb.Gpu.TileSet[tileId][y][x]
				pixelRealColor := gb.Gpu.BgPaletteMap[pixelPaletteColor]
				red, green, blue, _ := gb.Gpu.BgPaletteColors[pixelRealColor].RGBA()

				if gb.CgbModeEnabled {
					palette := tileAttr.BgPaletteNumber
					red, green, blue, _ = gb.Gpu.CgbBgPaletteColors[palette][pixelPaletteColor].RGBA()
				}
				pixelPos := (yBase*8+y)*t.GetWidth() + (xBase*8 + x)
				if 4*pixelPos+3 < len(frame) {
					frame[4*pixelPos] = byte(red)
					frame[4*pixelPos+1] = byte(green)
					frame[4*pixelPos+2] = byte(blue)
					frame[4*pixelPos+3] = 0xFF
				}
			}
		}
	}

	if (gb.Gpu.BackgroundMap && t.MapId == 1) || (!gb.Gpu.BackgroundMap && t.MapId == 0) {
		for x := 0; x < t.GetWidth(); x++ {
			pixelPos := int(gb.Gpu.ScrollY)*t.GetWidth() + x
			frame[4*pixelPos] = 0xDD
			frame[4*pixelPos+1] = 0x00
			frame[4*pixelPos+2] = 0x00
			frame[4*pixelPos+3] = 0xFF
		}
		for y := 0; y < t.GetHeight(); y++ {
			pixelPos := y*t.GetWidth() + int(gb.Gpu.ScrollX)
			frame[4*pixelPos] = 0xDD
			frame[4*pixelPos+1] = 0x00
			frame[4*pixelPos+2] = 0x00
			frame[4*pixelPos+3] = 0xFF
		}
	}

	return frame
}

func (t *BgMap) Draw(gb *internal.Gameboy, screen *ebiten.Image) {
	subscreen := screen.SubImage(image.Rect(t.GetXPos(), t.GetYPos(), t.GetXPos()+t.GetWidth(), t.GetYPos()+t.GetHeight())).(*ebiten.Image)
	subscreen.ReplacePixels(t.GetAsBytearray(gb))
}
