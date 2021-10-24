package cartridge

import (
	"fmt"
	"io"
	"log"
)

type Mbc5Cartridge struct {
	BaseCartridge
	Ram        []byte
	RamEnabled bool
	RomBank    int
	RamBank    int
}

func NewMbc5Cartridge(header *CartridgeHeader, data []byte) *Mbc5Cartridge {
	c := &Mbc5Cartridge{
		BaseCartridge: BaseCartridge{
			Header:     header,
			BinaryData: data,
		},
		RomBank: 1,
		Ram:     make([]byte, header.RamSize),
	}
	log.Printf("Detected %d Rom Banks", (len(c.BinaryData) / 0x4000))
	return c
}

func (c *Mbc5Cartridge) ReadByte(address uint16) uint8 {
	switch address & 0xF000 {
	case 0x0000, 0x1000, 0x2000, 0x3000:
		return c.BinaryData[address]
	case 0x4000, 0x5000, 0x6000, 0x7000:
		return c.BinaryData[c.RomBank*0x4000+int(address&0x3fff)]
	case 0xA000, 0xB000:
		if c.RamEnabled {
			return c.Ram[c.RamBank*0x2000+int(address&0x1fff)]
		} else {
			return 0xFF
		}
	default:
		return 0
	}
}

func (c *Mbc5Cartridge) WriteByte(address uint16, value uint8) {
	switch address & 0xF000 {
	case 0x0000, 0x1000:
		switch c.Header.Type {
		case Mbc5Ram, Mbc5RamBattery:
			c.RamEnabled = value == 0x0A
		default:
			return
		}
	case 0x2000:
		c.RomBank = int((c.RomBank)&0xFF00) + int(value)
		if int(c.RomBank)*0x4000 >= len(c.BinaryData) {
			c.RomBank = int(c.RomBank) % (len(c.BinaryData) / 0x4000)
		}
		// log.Printf("LoSelected RomBank %d (0x%02x)", c.RomBank, value)
		// log.Printf("Beginning of Bank: Read 0x%02x vs Real 0x%02x", c.ReadByte(0x4000), c.BinaryData[c.RomBank*0x4000])
		// address := 0x4000
		// log.Printf("%x %x %x", c.RomBank*0x4000+int(address&0x3fff), c.RomBank*0x4000, int(address&0x3fff))
	case 0x3000:
		c.RomBank = (c.RomBank & 0x00FF) + (int(value&0x1) << 8)
		if int(c.RomBank)*0x4000 >= len(c.BinaryData) {
			c.RomBank = int(c.RomBank) % (len(c.BinaryData) / 0x4000)
		}
		// log.Printf("HiSelected RomBank %d (0x%02x)", c.RomBank, value)
	case 0x4000, 0x5000:
		c.RamBank = int(value) & 0xF
		if len(c.Ram) <= 0 {
			c.RamBank = 0
		} else if int(c.RamBank)*0x2000 >= len(c.Ram) {
			c.RamBank = int(int(c.RamBank) % (len(c.Ram) / 0x2000))
		}

	case 0x6000, 0x7000:
	case 0xA000, 0xB000:
		if c.RamEnabled {
			c.Ram[uint16(c.RamBank*0x2000+int(address&0x1fff))] = value
		}
	default:
		return
	}
}
func (c *Mbc5Cartridge) UpdateComponentsPerCycle(cycles uint16) {}

func (c *Mbc5Cartridge) String() string {
	return fmt.Sprintf("%v %d %d", c.RamEnabled, c.RamBank, c.RomBank)
}

func (c *Mbc5Cartridge) SaveRam(writer io.Writer) error {
	_, err := writer.Write(c.Ram)
	return err
}
func (c *Mbc5Cartridge) LoadRam(reader io.Reader) error {
	var err error
	_, err = reader.Read(c.Ram)
	if err != io.EOF {
		return err
	}
	return nil
}
