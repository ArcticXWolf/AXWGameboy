package cartridge

import (
	"fmt"
	"io/ioutil"
)

type Cartridge interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)
	String() string
	CartridgeInfo() string
}

func LoadCartridge(filename string) (Cartridge, error) {
	header, data, err := LoadDataFromRomFile(filename)
	if err != nil {
		return nil, err
	}

	switch header.Type {
	case Rom:
		return NewRomCartridge(header, data), nil
	case Mbc1, Mbc1Ram, Mbc1RamBattery:
		return NewMbc1Cartridge(header, data), nil
	default:
		return nil, fmt.Errorf("cartridge type %#v not implemented yet", header.Type)
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

	// if ok := header.IsGlobalChecksumValid(data); !ok {
	// 	return nil, nil, errors.New("global checksum mismatch")
	// }

	return header, data, nil
}
