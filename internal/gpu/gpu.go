package gpu

type Gpu struct {
	backgroundActivated bool
	backgroundMap       bool
	backgroundTile      bool
	lcdActivated        bool
	scrollX             uint8
	scrollY             uint8
	CurrentScanline     uint8
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
		if g.backgroundTile {
			bT = 0x10
		}
		if g.lcdActivated {
			lA = 0x80
		}
		return bA | bM | bT | lA
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
		g.backgroundTile = value&0x10 != 0
		g.lcdActivated = value&0x80 != 0
	case 0xFF42:
		g.scrollY = value
	case 0xFF43:
		g.scrollX = value
	case 0xFF47:
	default:
	}
}
