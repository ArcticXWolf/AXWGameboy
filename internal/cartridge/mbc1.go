package cartridge

import (
	"fmt"
)

type Mbc1Cartridge struct {
	BaseCartridge
	Ram        [0x8000]byte
	RamEnabled bool
	RamMode    bool
	RomBank    uint16
	RamBank    uint16
}

func NewMbc1Cartridge(header *CartridgeHeader, data []byte) *Mbc1Cartridge {
	return &Mbc1Cartridge{
		BaseCartridge: BaseCartridge{
			Header:     header,
			BinaryData: data,
		},
		RomBank: 1,
	}
}

func (c *Mbc1Cartridge) ReadByte(address uint16) uint8 {
	switch address & 0xF000 {
	case 0x0000, 0x1000, 0x2000, 0x3000:
		return c.BinaryData[address]
	case 0x4000, 0x5000, 0x6000, 0x7000:
		return c.BinaryData[c.RomBank*0x4000+(address&0x3fff)]
	case 0xA000, 0xB000:
		return c.Ram[c.RamBank*0x2000+(address&0x1fff)]
	default:
		return 0
	}
}

func (c *Mbc1Cartridge) WriteByte(address uint16, value uint8) {
	switch address & 0xF000 {
	case 0x0000, 0x1000:
		switch c.Header.Type {
		case Mbc1Ram, Mbc1RamBattery:
			c.RamEnabled = (value & 0x0F) == 0x0A
		default:
			c.RamEnabled = false
		}
	case 0x2000, 0x3000:
		value &= 0x1F
		if value == 0 {
			value = 1
		}
		c.RomBank = (c.RomBank & 0x60) + uint16(value)
	case 0x4000, 0x5000:
		if c.RamMode {
			c.RamBank = uint16(value) & 0x3
		} else {
			value = (value & 0x3) << 5
			c.RomBank = (c.RomBank & 0x1F) + uint16(value)
		}
	case 0x6000, 0x7000:
		c.RamMode = value&0x1 == 1
	case 0xA000, 0xB000:
		c.Ram[c.RamBank*0x2000+(address&0x1fff)] = value
	default:
		return
	}
}

func (c *Mbc1Cartridge) String() string {
	return fmt.Sprintf("%v %d %v %d", c.RamEnabled, c.RamBank, c.RamMode, c.RomBank)
}
