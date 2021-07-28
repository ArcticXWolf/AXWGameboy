package internal

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"go.janniklasrichter.de/axwgameboy/internal/cartridge"
)

//go:embed bootroms/dmg_bios.bin
var gb_bios []byte

//go:embed bootroms/cgb_bios.bin
var cgb_bios []byte

type InterruptFlags struct {
	EnableFlags    uint8
	TriggeredFlags uint8
}

func (isr *InterruptFlags) GetEnableFlagVblank() bool {
	return isr.EnableFlags&0x1 != 0
}
func (isr *InterruptFlags) GetEnableFlagLCD() bool {
	return isr.EnableFlags&0x2 != 0
}
func (isr *InterruptFlags) GetEnableFlagTimer() bool {
	return isr.EnableFlags&0x4 != 0
}
func (isr *InterruptFlags) GetEnableFlagSerial() bool {
	return isr.EnableFlags&0x8 != 0
}
func (isr *InterruptFlags) GetEnableFlagJoypad() bool {
	return isr.EnableFlags&0x10 != 0
}
func (isr *InterruptFlags) IsTriggeredVblank() bool {
	return isr.TriggeredFlags&0x1 != 0
}
func (isr *InterruptFlags) IsTriggeredLCD() bool {
	return isr.TriggeredFlags&0x2 != 0
}
func (isr *InterruptFlags) IsTriggeredTimer() bool {
	return isr.TriggeredFlags&0x4 != 0
}
func (isr *InterruptFlags) IsTriggeredSerial() bool {
	return isr.TriggeredFlags&0x8 != 0
}
func (isr *InterruptFlags) IsTriggeredJoypad() bool {
	return isr.TriggeredFlags&0x10 != 0
}

type Mmu struct {
	inbios       bool
	gb           *Gameboy
	bios         []byte
	Cartridge    cartridge.Cartridge
	wramBank     int
	wram         [0x9000]byte
	serialOutput byte
	zram         [0x80]byte
	isr          *InterruptFlags
}

func (m *Mmu) GetInterruptFlags() *InterruptFlags {
	return m.isr
}

func (m *Mmu) String() string {
	return fmt.Sprintf("Memory: %v", *m)
}

func NewMemory(gb *Gameboy) (*Mmu, bool, error) {
	var cart cartridge.Cartridge
	var err error

	if gb.Options.RomPath != "" {
		cart, err = cartridge.LoadCartridgeFromPath(gb.Options.RomPath)
		if err != nil {
			return nil, false, err
		}
	} else {
		log.Panic("no rom loaded")
	}

	romCGBEnabled := cart.CartridgeHeader().CartridgeGBMode == cartridge.OnlyCGB || cart.CartridgeHeader().CartridgeGBMode == cartridge.SupportsCGB

	if gb.Options.SavePath != "" {
		if _, err := os.Stat(gb.Options.SavePath); err == nil {
			cart.LoadRam(gb.Options.SavePath)
		}
	}

	m := &Mmu{
		inbios:    true,
		gb:        gb,
		bios:      gb_bios,
		Cartridge: cart,
		wram:      [0x9000]byte{},
		zram:      [0x80]byte{},
		isr: &InterruptFlags{
			TriggeredFlags: 0xE0,
		},
	}

	if m.gb.Options.CGBEnabled {
		m.bios = cgb_bios
	}

	return m, romCGBEnabled, err
}

func (m *Mmu) ReadByte(address uint16) (result uint8) {
	switch address & 0xF000 {
	case 0x0000: // ROM / BIOS
		if m.inbios && address < 0x0100 {
			return m.bios[address]
		}
		if m.inbios && m.gb.cgbModeEnabled && address >= 0x200 && address < 0x900 {
			return m.bios[address]
		}
		return m.Cartridge.ReadByte(address)
	case 0x1000, 0x2000, 0x3000, 0x4000, 0x5000, 0x6000, 0x7000: // ROM
		return m.Cartridge.ReadByte(address)
	case 0x8000, 0x9000: // VRAM
		return m.gb.Gpu.ReadByte(address)
	case 0xA000, 0xB000: // External RAM
		return m.Cartridge.ReadByte(address)
	case 0xC000: // Working RAM
		return m.wram[address&0x0FFF]
	case 0xD000: // Working RAM
		return m.wram[address&0x0FFF+0x1000+uint16(m.wramBank)*0x1000]
	case 0xE000, 0xF000:
		if address < 0xFE00 { // Working RAM Shadow
			if address < 0xF000 {
				return m.wram[address&0x0FFF]
			}
			return m.wram[address&0x0FFF+0x1000+uint16(m.wramBank)*0x1000]
		}
		if address < 0xFEA0 { // OAM
			return m.gb.Gpu.ReadByte(address)
		}
		if address < 0xFF00 { // Empty
			return 0
		}
		if address < 0xFF80 { // I/O
			switch address & 0x00F0 {
			case 0x00:
				if address == 0xFF00 {
					return m.gb.Inputs.ReadByte(address)
				}
				if address == 0xFF01 {
					return m.serialOutput
				}
				if address == 0xFF02 {
					return 0x00
				}
				if address >= 0xFF04 && address <= 0xFF07 {
					return m.gb.Timer.ReadByte(address)
				}
				if address == 0xFF0F {
					return m.isr.TriggeredFlags
				}
			case 0x10, 0x20, 0x30:
				return m.gb.Apu.ReadByte(address)
			case 0x40, 0x50, 0x60, 0x70:
				if address == 0xFF70 {
					if m.gb.cgbModeEnabled {
						return uint8(m.wramBank)
					}
				}
				return m.gb.Gpu.ReadByte(address)
			}
			return 0 // TODO
		}
		if address == 0xFFFF {
			return m.isr.EnableFlags
		}
		return m.zram[address&0x7F] // Highspeed Zero RAM
	default:
		return 0
	}
}

func (m *Mmu) ReadWord(address uint16) (result uint16) {
	return uint16(m.ReadByte(address)) + uint16(m.ReadByte(address+1))<<8
}

func (m *Mmu) WriteByte(address uint16, value uint8) {
	switch address & 0xF000 {
	case 0x0000: // ROM / BIOS
		m.Cartridge.WriteByte(address, value)
		return
	case 0x1000, 0x2000, 0x3000, 0x4000, 0x5000, 0x6000, 0x7000: // ROM
		m.Cartridge.WriteByte(address, value)
		return
	case 0x8000, 0x9000: // VRAM
		m.gb.Gpu.WriteByte(address, value)
		return
	case 0xA000, 0xB000: // External RAM
		m.Cartridge.WriteByte(address, value)
		return
	case 0xC000: // Working RAM
		m.wram[address&0x0FFF] = value
	case 0xD000: // Working RAM
		m.wram[address&0x0FFF+0x1000+uint16(m.wramBank)*0x1000] = value
	case 0xE000, 0xF000:
		if address < 0xFE00 { // Working RAM Shadow
			if address < 0xF000 {
				m.wram[address&0x0FFF] = value
				return
			}
			m.wram[address&0x0FFF+0x1000+uint16(m.wramBank)*0x1000] = value
			return
		}
		if address < 0xFEA0 { // OAM
			m.gb.Gpu.WriteByte(address, value)
		}
		if address < 0xFF00 { // Empty
			return
		}
		if address < 0xFF80 { // I/O
			switch address & 0x00F0 {
			case 0x00:
				if address == 0xFF00 {
					m.gb.Inputs.WriteByte(address, value)
					return
				}
				if address == 0xFF01 {
					m.serialOutput = value
					return
				}
				if address == 0xFF02 {
					sof := m.gb.Options.SerialOutputFunction
					if sof != nil {
						sof(m.serialOutput)
					}
					return
				}
				if address >= 0xFF04 && address <= 0xFF07 {
					m.gb.Timer.WriteByte(address, value)
				}
				if address == 0xFF0F {
					m.isr.TriggeredFlags = 0xE0 | value&(^uint8(0xE0))
				}
			case 0x10, 0x20:
				m.gb.Apu.WriteByte(address, value)
			case 0x30:
				m.gb.Apu.WriteWaveform(address, value)
			case 0x40, 0x50, 0x60, 0x70:
				if address == 0xFF50 {
					m.inbios = false
					return
				}
				if address == 0xFF70 {
					if m.gb.cgbModeEnabled {
						m.wramBank = int(value & 0x7)
						if m.wramBank == 0 {
							m.wramBank = 1
						}
					}
				}
				m.gb.Gpu.WriteByte(address, value)
			}
			return
		}
		if address == 0xFFFF {
			m.isr.EnableFlags = value
		}
		m.zram[address&0x7F] = value
		return
	default:
		return
	}
}

func (m *Mmu) WriteWord(address uint16, value uint16) {
	m.WriteByte(address, uint8(value&0xFF))
	m.WriteByte(address+1, uint8(value>>8))
}
