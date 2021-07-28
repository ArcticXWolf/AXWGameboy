package cartridge

type BaseCartridge struct {
	BinaryData []byte
	Header     *CartridgeHeader
}

func (b *BaseCartridge) CartridgeInfo() string {
	return b.Header.String()
}

func (b *BaseCartridge) CartridgeHeader() *CartridgeHeader {
	return b.Header
}
