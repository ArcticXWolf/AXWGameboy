package internal

import (
	"fmt"
	"os"

	"go.janniklasrichter.de/axwgameboy/internal/cartridge"
)

var gb_bios = [0x100]byte{
	0x31, 0xFE, 0xFF, 0xAF, 0x21, 0xFF, 0x9F, 0x32, 0xCB, 0x7C, 0x20, 0xFB, 0x21, 0x26, 0xFF, 0x0E,
	0x11, 0x3E, 0x80, 0x32, 0xE2, 0x0C, 0x3E, 0xF3, 0xE2, 0x32, 0x3E, 0x77, 0x77, 0x3E, 0xFC, 0xE0,
	0x47, 0x11, 0x04, 0x01, 0x21, 0x10, 0x80, 0x1A, 0xCD, 0x95, 0x00, 0xCD, 0x96, 0x00, 0x13, 0x7B,
	0xFE, 0x34, 0x20, 0xF3, 0x11, 0xD8, 0x00, 0x06, 0x08, 0x1A, 0x13, 0x22, 0x23, 0x05, 0x20, 0xF9,
	0x3E, 0x19, 0xEA, 0x10, 0x99, 0x21, 0x2F, 0x99, 0x0E, 0x0C, 0x3D, 0x28, 0x08, 0x32, 0x0D, 0x20,
	0xF9, 0x2E, 0x0F, 0x18, 0xF3, 0x67, 0x3E, 0x64, 0x57, 0xE0, 0x42, 0x3E, 0x91, 0xE0, 0x40, 0x04,
	0x1E, 0x02, 0x0E, 0x0C, 0xF0, 0x44, 0xFE, 0x90, 0x20, 0xFA, 0x0D, 0x20, 0xF7, 0x1D, 0x20, 0xF2,
	0x0E, 0x13, 0x24, 0x7C, 0x1E, 0x83, 0xFE, 0x62, 0x28, 0x06, 0x1E, 0xC1, 0xFE, 0x64, 0x20, 0x06,
	0x7B, 0xE2, 0x0C, 0x3E, 0x87, 0xF2, 0xF0, 0x42, 0x90, 0xE0, 0x42, 0x15, 0x20, 0xD2, 0x05, 0x20,
	0x4F, 0x16, 0x20, 0x18, 0xCB, 0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
	0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9, 0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
	0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
	0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E, 0x3c, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x3C,
	0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x20, 0xFE, 0x23, 0x7D, 0xFE, 0x34, 0x20,
	0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x20, 0xFE, 0x3E, 0x01, 0xE0, 0x50,
}

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
	bios         [0x100]byte
	Cartridge    cartridge.Cartridge
	wram         [0x2000]byte
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

func NewMemory(gb *Gameboy) (*Mmu, error) {
	var cart cartridge.Cartridge
	var err error

	if gb.Options.RomPath != "" {
		cart, err = cartridge.LoadCartridgeFromPath(gb.Options.RomPath)
		if err != nil {
			return nil, err
		}
	} else {
		cart, err = cartridge.LoadEmbeddedCartridge()
		if err != nil {
			return nil, err
		}
	}

	if gb.Options.SavePath != "" {
		if _, err := os.Stat(gb.Options.SavePath); err == nil {
			err = cart.LoadRam(gb.Options.SavePath)
		}
	}

	return &Mmu{
		inbios:    true,
		gb:        gb,
		bios:      gb_bios,
		Cartridge: cart,
		wram:      [0x2000]byte{},
		zram:      [0x80]byte{},
		isr: &InterruptFlags{
			TriggeredFlags: 0xE0,
		},
	}, err
}

func (m *Mmu) ReadByte(address uint16) (result uint8) {
	switch address & 0xF000 {
	case 0x0000: // ROM / BIOS
		if m.inbios && address < 0x0100 {
			return m.bios[address]
		} else if m.inbios && address == 0x0100 {
			m.inbios = false
		}
		return m.Cartridge.ReadByte(address)
	case 0x1000, 0x2000, 0x3000, 0x4000, 0x5000, 0x6000, 0x7000: // ROM
		return m.Cartridge.ReadByte(address)
	case 0x8000, 0x9000: // VRAM
		return m.gb.Gpu.vram[address&0x1FFF]
	case 0xA000, 0xB000: // External RAM
		return m.Cartridge.ReadByte(address)
	case 0xC000, 0xD000: // Working RAM
		return m.wram[address&0x1FFF]
	case 0xE000, 0xF000:
		if address < 0xFE00 { // Working RAM Shadow
			return m.wram[address&0x1FFF]
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
		m.gb.Gpu.vram[address&0x1FFF] = value
		m.gb.Gpu.updateTile(address)
		return
	case 0xA000, 0xB000: // External RAM
		m.Cartridge.WriteByte(address, value)
		return
	case 0xC000, 0xD000: // Working RAM
		m.wram[address&0x1FFF] = value
		return
	case 0xE000, 0xF000:
		if address < 0xFE00 { // Working RAM Shadow
			m.wram[address&0x1FFF] = value
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
