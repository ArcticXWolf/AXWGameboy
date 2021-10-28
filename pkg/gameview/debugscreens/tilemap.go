package debugscreens

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type Tilemap struct {
	X int
	Y int
}

func (t *Tilemap) GetXPos() int   { return t.X }
func (t *Tilemap) GetYPos() int   { return t.Y }
func (t *Tilemap) GetWidth() int  { return 16 * 8 }
func (t *Tilemap) GetHeight() int { return 48 * 8 }

func (t *Tilemap) GetTilemapAsBytearray(gb *internal.Gameboy) []byte {
	var frame []byte = make([]byte, 4*t.GetWidth()*t.GetHeight())

	palette := [4]color.Color{
		color.RGBA{0x00, 0x00, 0x00, 255},
		color.RGBA{0x66, 0x66, 0x66, 255},
		color.RGBA{0xAA, 0xAA, 0xAA, 255},
		color.RGBA{0xFF, 0xFF, 0xFF, 255},
	}

	for tileId := 0; tileId < len(gb.Gpu.TileSet); tileId++ {
		xBase := tileId % (t.GetWidth() / 8)
		yBase := tileId / (t.GetWidth() / 8)
		for x := 0; x < 8; x++ {
			for y := 0; y < 8; y++ {
				pixelPaletteColor := gb.Gpu.TileSet[tileId][y][x]
				red, green, blue, _ := palette[pixelPaletteColor].RGBA()
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

	return frame
}

func (t *Tilemap) Draw(gb *internal.Gameboy, screen *ebiten.Image) {
	subscreen := screen.SubImage(image.Rect(t.GetXPos(), t.GetYPos(), t.GetXPos()+t.GetWidth(), t.GetYPos()+t.GetHeight())).(*ebiten.Image)
	subscreen.ReplacePixels(t.GetTilemapAsBytearray(gb))
}
