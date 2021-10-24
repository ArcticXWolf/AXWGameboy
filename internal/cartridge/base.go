package cartridge

import (
	"crypto/sha256"
)

type BaseCartridge struct {
	BinaryData []byte
	Header     *CartridgeHeader
}

func (b *BaseCartridge) CartridgeInfo() string {
	return b.Header.String()
}

func (b *BaseCartridge) CartridgeHash() [32]byte {
	return sha256.Sum256(b.BinaryData)
}

func (b *BaseCartridge) CartridgeHeader() *CartridgeHeader {
	return b.Header
}
