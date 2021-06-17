package cartridge

type BaseCartridge struct {
	BinaryData []byte
	Header     *CartridgeHeader
}
