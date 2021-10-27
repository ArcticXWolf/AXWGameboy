package cartridge

import (
	"fmt"
	"io"
)

type Mbc1Cartridge struct {
	BaseCartridge
	Ram        []byte
	RamEnabled bool
	RamMode    bool
	RomBank    int
	RamBank    int
}

func NewMbc1Cartridge(header *CartridgeHeader, data []byte) *Mbc1Cartridge {
	return &Mbc1Cartridge{
		BaseCartridge: BaseCartridge{
			Header:     header,
			BinaryData: data,
		},
		RomBank: 1,
		Ram:     make([]byte, header.RamSize),
	}
}

func (c *Mbc1Cartridge) ReadByte(address uint16) uint8 {
	switch address & 0xF000 {
	case 0x0000, 0x1000, 0x2000, 0x3000:
		return c.BinaryData[address]
	case 0x4000, 0x5000, 0x6000, 0x7000:
		return c.BinaryData[c.RomBank*0x4000+int(address&0x3fff)]
	case 0xA000, 0xB000:
		if c.RamEnabled {
			if c.RamMode {
				return c.Ram[c.RamBank*0x2000+int(address&0x1fff)]
			} else {
				return c.Ram[address&0x1fff]
			}
		} else {
			return 0xFF
		}
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
			return
		}
	case 0x2000, 0x3000:
		value &= 0x1F
		c.RomBank = (c.RomBank & 0x60) + int(value)
		if int(c.RomBank)*0x4000 >= len(c.BinaryData) {
			c.RomBank = int(c.RomBank) % (len(c.BinaryData) / 0x4000)
		}
		if c.RomBank == 0 {
			c.RomBank = 1
		}
	case 0x4000, 0x5000:
		if c.RamMode {
			c.RamBank = int(value) & 0x3
			if len(c.Ram) <= 0 {
				c.RamBank = 0
			} else if int(c.RamBank)*0x2000 >= len(c.Ram) {
				c.RamBank = int(c.RamBank) % (len(c.Ram) / 0x2000)
			}
		} else {
			value = (value & 0x3) << 5
			c.RomBank = (c.RomBank & 0x1F) + int(value)
			if int(c.RomBank)*0x4000 >= len(c.BinaryData) {
				c.RomBank = int(c.RomBank) % (len(c.BinaryData) / 0x4000)
			}
			if c.RomBank == 0 {
				c.RomBank = 1
			}
		}
	case 0x6000, 0x7000:
		c.RamMode = value&0x1 == 1
	case 0xA000, 0xB000:
		if c.RamEnabled {
			if c.RamMode {
				c.Ram[c.RamBank*0x2000+int(address&0x1fff)] = value
			} else {
				c.Ram[address&0x1fff] = value
			}
		}
	default:
		return
	}
}

func (c *Mbc1Cartridge) UpdateComponentsPerCycle(cycles uint16) {}

func (c *Mbc1Cartridge) String() string {
	return fmt.Sprintf("%v %d %v %d", c.RamEnabled, c.RamBank, c.RamMode, c.RomBank)
}

func (c *Mbc1Cartridge) SaveRam(writer io.Writer) error {
	_, err := writer.Write(c.Ram)
	return err
}
func (c *Mbc1Cartridge) LoadRam(reader io.Reader) error {
	var err error
	_, err = reader.Read(c.Ram)
	if err != io.EOF {
		return err
	}
	return nil
}

func (c *Mbc1Cartridge) GetRamBank() uint8 { return uint8(c.RamBank) }
