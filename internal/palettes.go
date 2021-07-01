package internal

import "image/color"

func getPaletteColorsByName(name string) [4]color.Color {
	switch name {
	case "Red":
		return [4]color.Color{
			color.RGBA{0xff, 0xf6, 0xd3, 255},
			color.RGBA{0xf9, 0xa8, 0x75, 255},
			color.RGBA{0xeb, 0x6b, 0x6f, 255},
			color.RGBA{0x7c, 0x3f, 0x58, 255},
		}
	case "DMG":
		return [4]color.Color{
			color.RGBA{0xe0, 0xf8, 0xd0, 255},
			color.RGBA{0x88, 0xc0, 0x70, 255},
			color.RGBA{0x34, 0x68, 0x56, 255},
			color.RGBA{0x08, 0x18, 0x20, 255},
		}
	default:
		return [4]color.Color{
			color.RGBA{255, 255, 255, 255},
			color.RGBA{192, 192, 192, 255},
			color.RGBA{96, 96, 96, 255},
			color.RGBA{0, 0, 0, 255},
		}
	}
}
