package cartridge

import (
	_ "embed"
	"fmt"
	"io/ioutil"
)

type Cartridge interface {
	ReadByte(address uint16) uint8
	WriteByte(address uint16, value uint8)
	String() string
	CartridgeInfo() string
	CartridgeHeader() *CartridgeHeader
	SaveRam(filename string) error
	LoadRam(filename string) error
	UpdateComponentsPerCycle()
}

func LoadCartridgeFromPath(filename string) (Cartridge, error) {
	data, err := LoadDataFromRomFile(filename)
	if err != nil {
		return nil, err
	}
	header, err := ParseHeaderFromRomData(data)
	if err != nil {
		return nil, err
	}
	return InitializeCartridge(header, data)
}

func LoadDataFromRomFile(filepath string) (data []byte, err error) {
	data, err = ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ParseHeaderFromRomData(data []byte) (header *CartridgeHeader, err error) {
	header, err = NewCartridgeHeader(data[0x100:0x150])
	if err != nil {
		return nil, err
	}

	// if ok := header.IsGlobalChecksumValid(data); !ok {
	// 	return nil, errors.New("global checksum mismatch")
	// }

	return header, nil
}

func InitializeCartridge(header *CartridgeHeader, data []byte) (Cartridge, error) {
	switch header.Type {
	case Rom:
		return NewRomCartridge(header, data), nil
	case Mbc1, Mbc1Ram, Mbc1RamBattery:
		return NewMbc1Cartridge(header, data), nil
	case Mbc3, Mbc3Ram, Mbc3RamBattery, Mbc3TimerRamBattery, Mbc3TimerBattery:
		return NewMbc3Cartridge(header, data), nil
	case Mbc5, Mbc5Ram, Mbc5RamBattery:
		return NewMbc5Cartridge(header, data), nil
	default:
		return nil, fmt.Errorf("cartridge type %#v not implemented yet", header.Type)
	}
}
