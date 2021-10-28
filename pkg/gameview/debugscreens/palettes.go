package debugscreens

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
)

type PaletteList struct {
	X int
	Y int
}

func (t *PaletteList) GetXPos() int   { return t.X }
func (t *PaletteList) GetYPos() int   { return t.Y }
func (t *PaletteList) GetWidth() int  { return 4 * 10 }
func (t *PaletteList) GetHeight() int { return 19 * 11 }

func (t *PaletteList) GetPaletteListAsBytearray(gb *internal.Gameboy) []byte {
	var frame []byte = make([]byte, 4*t.GetWidth()*t.GetHeight())

	// DMG
	paletteId := 0
	for c := 0; c < 4; c++ {
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				pixelPos := (paletteId*11+y)*t.GetWidth() + (c*10 + x)
				pixelPalette := gb.Gpu.BgPaletteMap[c]
				red, green, blue, _ := gb.Gpu.BgPaletteColors[pixelPalette].RGBA()
				if 4*pixelPos+3 < len(frame) {
					frame[4*pixelPos] = byte(red)
					frame[4*pixelPos+1] = byte(green)
					frame[4*pixelPos+2] = byte(blue)
					frame[4*pixelPos+3] = 0xFF
				}
			}
		}
	}
	paletteId++
	for c := 0; c < 4; c++ {
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				pixelPos := (paletteId*11+y)*t.GetWidth() + (c*10 + x)
				pixelPalette := gb.Gpu.SpritePaletteMap[0][c]
				red, green, blue, _ := gb.Gpu.SpritePaletteColors[0][pixelPalette].RGBA()
				if 4*pixelPos+3 < len(frame) {
					frame[4*pixelPos] = byte(red)
					frame[4*pixelPos+1] = byte(green)
					frame[4*pixelPos+2] = byte(blue)
					frame[4*pixelPos+3] = 0xFF
				}
			}
		}
	}
	paletteId++
	for c := 0; c < 4; c++ {
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				pixelPos := (paletteId*11+y)*t.GetWidth() + (c*10 + x)
				pixelPalette := gb.Gpu.SpritePaletteMap[1][c]
				red, green, blue, _ := gb.Gpu.SpritePaletteColors[1][pixelPalette].RGBA()
				if 4*pixelPos+3 < len(frame) {
					frame[4*pixelPos] = byte(red)
					frame[4*pixelPos+1] = byte(green)
					frame[4*pixelPos+2] = byte(blue)
					frame[4*pixelPos+3] = 0xFF
				}
			}
		}
	}

	// CGB
	for p := 0; p < len(gb.Gpu.CgbBgPaletteColors); p++ {
		paletteId++
		for c := 0; c < 4; c++ {
			for x := 0; x < 10; x++ {
				for y := 0; y < 10; y++ {
					pixelPos := (paletteId*11+y)*t.GetWidth() + (c*10 + x)
					red, green, blue, _ := gb.Gpu.CgbBgPaletteColors[p][c].RGBA()
					if 4*pixelPos+3 < len(frame) {
						frame[4*pixelPos] = byte(red)
						frame[4*pixelPos+1] = byte(green)
						frame[4*pixelPos+2] = byte(blue)
						frame[4*pixelPos+3] = 0xFF
					}
				}
			}
		}
	}
	for p := 0; p < len(gb.Gpu.CgbObjPaletteColors); p++ {
		paletteId++
		for c := 0; c < 4; c++ {
			for x := 0; x < 10; x++ {
				for y := 0; y < 10; y++ {
					pixelPos := (paletteId*11+y)*t.GetWidth() + (c*10 + x)
					red, green, blue, _ := gb.Gpu.CgbObjPaletteColors[p][c].RGBA()
					if 4*pixelPos+3 < len(frame) {
						frame[4*pixelPos] = byte(red)
						frame[4*pixelPos+1] = byte(green)
						frame[4*pixelPos+2] = byte(blue)
						frame[4*pixelPos+3] = 0xFF
					}
				}
			}
		}
	}

	return frame
}

func (t *PaletteList) Draw(gb *internal.Gameboy, screen *ebiten.Image) {
	subscreen := screen.SubImage(image.Rect(t.GetXPos(), t.GetYPos(), t.GetXPos()+t.GetWidth(), t.GetYPos()+t.GetHeight())).(*ebiten.Image)
	subscreen.ReplacePixels(t.GetPaletteListAsBytearray(gb))
}
