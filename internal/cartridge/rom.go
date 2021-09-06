package cartridge

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

func (rc *RomCartridge) WriteByte(address uint16, value uint8) {}

func (c *RomCartridge) UpdateComponentsPerCycle(cycles uint16) {}

func (c *RomCartridge) String() string {
	return "ROM Cartridge"
}

func (c *RomCartridge) SaveRam(filename string) error { return nil }
func (c *RomCartridge) LoadRam(filename string) error { return nil }
