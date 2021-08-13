package cartridge

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

// Mostly from https://gbdev.io/pandocs/The_Cartridge_Header.html

type CartridgeGBMode uint8

const (
	OnlyDMG CartridgeGBMode = iota
	SupportsCGB
	OnlyCGB
)

type CartridgeType uint8

const (
	Rom                        CartridgeType = 0x00
	Mbc1                       CartridgeType = 0x01
	Mbc1Ram                    CartridgeType = 0x02
	Mbc1RamBattery             CartridgeType = 0x03
	Mbc2                       CartridgeType = 0x05
	Mbc2Battery                CartridgeType = 0x06
	RomRam                     CartridgeType = 0x08
	RomRamBattery              CartridgeType = 0x09
	Mmm01                      CartridgeType = 0x0b
	Mmm01Ram                   CartridgeType = 0x0c
	Mmm01RamBattery            CartridgeType = 0x0d
	Mbc3TimerBattery           CartridgeType = 0x0f
	Mbc3TimerRamBattery        CartridgeType = 0x10
	Mbc3                       CartridgeType = 0x11
	Mbc3Ram                    CartridgeType = 0x12
	Mbc3RamBattery             CartridgeType = 0x13
	Mbc5                       CartridgeType = 0x19
	Mbc5Ram                    CartridgeType = 0x1a
	Mbc5RamBattery             CartridgeType = 0x1b
	Mbc5Rumble                 CartridgeType = 0x1c
	Mbc5RumbleRam              CartridgeType = 0x1d
	Mbc5RumbleRamBattery       CartridgeType = 0x1e
	Mbc6                       CartridgeType = 0x20
	Mbc7SensorRumbleRamBattery CartridgeType = 0x22
	PocketCamera               CartridgeType = 0xfc
	BandaiTama5                CartridgeType = 0xfd
	HuC3                       CartridgeType = 0xfe
	HuC1RamBattery             CartridgeType = 0xff
)

type CartridgeHeader struct {
	Title            string
	ManufacturerCode string
	CartridgeGBMode  CartridgeGBMode
	NewLicenseeCode  string
	SupportsSGB      bool
	Type             CartridgeType
	RomSize          int
	RamSize          int
	CountryCode      uint8
	OldLicenseeCode  uint8
	Version          uint8
	HeaderChecksum   uint8
	GlobalChecksum   uint16
	HeaderBinary     [0x50]byte
}

func NewCartridgeHeader(binary []byte) (*CartridgeHeader, error) {
	if len(binary) != 0x50 {
		return nil, errors.New("rom header size mismatch")
	}

	var cleanBinary [0x50]byte
	for k, v := range binary {
		cleanBinary[k] = v
	}

	ch := &CartridgeHeader{
		HeaderBinary: cleanBinary,
	}

	ch.parseGBMode()
	ch.parseTitle()
	ch.parseManufacturerCode()
	ch.parseNewLicenseeCode()
	ch.parseSGBSupport()
	ch.parseType()
	ch.parseRomSize()
	ch.parseRamSize()
	ch.parseCountryCode()
	ch.parseOldLicenseeCode()
	ch.parseVersionNumber()
	ch.parseHeaderChecksum()
	ch.parseGlobalChecksum()

	if !ch.IsHeaderChecksumValid() {
		return nil, errors.New("cartridge header checksum mismatch")
	}

	return ch, nil
}

func (ch *CartridgeHeader) parseTitle() {
	if ch.CartridgeGBMode == OnlyDMG {
		ch.Title = string(bytes.Trim(ch.HeaderBinary[0x34:0x44], "\x00"))
		return
	}
	ch.Title = string(bytes.Trim(ch.HeaderBinary[0x34:0x3F], "\x00"))
}

func (ch *CartridgeHeader) parseManufacturerCode() {
	ch.ManufacturerCode = string(ch.HeaderBinary[0x3F:0x43])
}

func (ch *CartridgeHeader) parseGBMode() {
	switch ch.HeaderBinary[0x43] {
	case 0x80:
		ch.CartridgeGBMode = SupportsCGB
	case 0xC0:
		ch.CartridgeGBMode = OnlyCGB
	default:
		ch.CartridgeGBMode = OnlyDMG
	}
}

func (ch *CartridgeHeader) parseNewLicenseeCode() {
	ch.NewLicenseeCode = string(ch.HeaderBinary[0x44:0x46])
}

func (ch *CartridgeHeader) parseSGBSupport() {
	ch.SupportsSGB = ch.HeaderBinary[0x46] == 0x03
}

func (ch *CartridgeHeader) parseType() {
	ch.Type = CartridgeType(ch.HeaderBinary[0x47])
}

func (ch *CartridgeHeader) parseRomSize() {
	ch.RomSize = int(math.Pow(2, float64(15+int(ch.HeaderBinary[0x48]))))
}

func (ch *CartridgeHeader) parseRamSize() {
	switch ch.HeaderBinary[0x49] {
	case 0x00:
		ch.RamSize = 0
	case 0x02:
		ch.RamSize = 8 * 1024
	case 0x03:
		ch.RamSize = 32 * 1024
	case 0x04:
		ch.RamSize = 128 * 1024
	case 0x05:
		ch.RamSize = 64 * 1024
	}
}

func (ch *CartridgeHeader) parseCountryCode() {
	ch.CountryCode = ch.HeaderBinary[0x4A]
}

func (ch *CartridgeHeader) parseOldLicenseeCode() {
	ch.OldLicenseeCode = ch.HeaderBinary[0x4B]
}

func (ch *CartridgeHeader) parseVersionNumber() {
	ch.Version = ch.HeaderBinary[0x4C]
}

func (ch *CartridgeHeader) parseHeaderChecksum() {
	ch.HeaderChecksum = ch.HeaderBinary[0x4D]
}

func (ch *CartridgeHeader) parseGlobalChecksum() {
	ch.GlobalChecksum = uint16(ch.HeaderBinary[0x4E])<<8 | uint16(ch.HeaderBinary[0x4F])
}

func (ch *CartridgeHeader) IsHeaderChecksumValid() bool {
	var x uint8 = 0x19
	for i := 0x34; i <= 0x4C; i++ {
		x += ch.HeaderBinary[i]
	}

	return x+ch.HeaderChecksum == 0
}

func (ch *CartridgeHeader) IsGlobalChecksumValid(romBinary []byte) bool {
	var x uint16
	for k, v := range romBinary {
		if k == 0x14e || k == 0x14f {
			continue
		}
		x += uint16(v)
	}

	return x == ch.GlobalChecksum
}

func (ch *CartridgeHeader) String() string {
	return fmt.Sprintf("Cartridge %s [0x%02x] Mode %01d | Romsize %010d | Ramsize %010d", ch.Title, ch.Type, ch.CartridgeGBMode, ch.RomSize, ch.RamSize)
}
