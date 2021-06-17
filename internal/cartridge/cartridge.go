package cartridge

import (
	"errors"
	"io/ioutil"
)

type Cartridge interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)
}

func LoadCartridge(filename string) (Cartridge, error) {
	header, data, err := LoadDataFromRomFile(filename)
	if err != nil {
		return nil, err
	}

	switch header.Type {
	case Rom:
		return NewRomCartridge(header, data), nil
	default:
		return nil, errors.New("cartridge type not implemented yet")
	}
}

func LoadDataFromRomFile(filepath string) (header *CartridgeHeader, data []byte, err error) {
	data, err = ioutil.ReadFile(filepath)
	if err != nil {
		return nil, nil, err
	}

	header, err = NewCartridgeHeader(data[0x100:0x150])
	if err != nil {
		return nil, nil, err
	}

	if ok := header.IsGlobalChecksumValid(data); !ok {
		return nil, nil, errors.New("global checksum mismatch")
	}

	return header, data, nil
}
