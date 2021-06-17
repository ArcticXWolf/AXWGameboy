package cartridge

import "fmt"

type RomCartridge struct {
	BaseCartridge
}

func NewRomCartridge(header *CartridgeHeader, data []byte) *RomCartridge {
	return &RomCartridge{
		BaseCartridge: BaseCartridge{
			Header:     header,
			BinaryData: data,
		},
	}
}

func (rc *RomCartridge) ReadByte(address uint16) uint8 {
	if address <= 0x8000 {
		return rc.BinaryData[address]
	}
	return 0
}

func (rc *RomCartridge) WriteByte(address uint16, value uint8) {
	return
}

func (c *RomCartridge) String() string {
	return fmt.Sprintf("ROM Cartridge")
}
