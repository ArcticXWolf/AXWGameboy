package cartridge

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"
)

const (
	RtcRefreshAfterCycles uint16 = 35000
)

type Mbc3RTC struct {
	Seconds                     uint8
	Minutes                     uint8
	Hours                       uint8
	DaysLower                   uint8
	DaysHigherAndControl        uint8
	LatchedSeconds              uint8
	LatchedMinutes              uint8
	LatchedHours                uint8
	LatchedDaysLower            uint8
	LatchedDaysHigherAndControl uint8
}

type Mbc3Cartridge struct {
	BaseCartridge
	Ram                   []byte
	RamEnabled            bool
	RtcEnabled            bool
	RomBank               int
	RamBank               int
	hasRTC                bool
	Rtc                   *Mbc3RTC
	RtcRegister           uint8
	RtcLatchFlagValue     uint8
	RtcLastUpdate         time.Time
	RtcRefreshCyclesCount uint16
}

func NewMbc3Cartridge(header *CartridgeHeader, data []byte) *Mbc3Cartridge {
	return &Mbc3Cartridge{
		BaseCartridge: BaseCartridge{
			Header:     header,
			BinaryData: data,
		},
		RomBank:       1,
		Ram:           make([]byte, header.RamSize),
		hasRTC:        header.Type == 0x0f || header.Type == 0x10,
		Rtc:           &Mbc3RTC{},
		RtcLastUpdate: time.Now(),
	}
}

func (c *Mbc3Cartridge) ReadByte(address uint16) uint8 {
	switch address & 0xF000 {
	case 0x0000, 0x1000, 0x2000, 0x3000:
		return c.BinaryData[address]
	case 0x4000, 0x5000, 0x6000, 0x7000:
		return c.BinaryData[c.RomBank*0x4000+int(address&0x3fff)]
	case 0xA000, 0xB000:
		if c.RamBank >= 0 {
			if c.RamEnabled {
				return c.Ram[c.RamBank*0x2000+int(address&0x1fff)]
			} else {
				return 0xFF
			}
		} else if c.hasRTC && c.RtcEnabled {
			switch c.RtcRegister {
			case 0x8:
				return c.Rtc.LatchedSeconds
			case 0x9:
				return c.Rtc.LatchedMinutes
			case 0xa:
				return c.Rtc.LatchedHours
			case 0xb:
				return c.Rtc.LatchedDaysLower
			case 0xc:
				return c.Rtc.LatchedDaysHigherAndControl
			default:
				return 0xFF
			}
		} else {
			return 0xFF
		}
	default:
		return 0
	}
}

func (c *Mbc3Cartridge) WriteByte(address uint16, value uint8) {
	switch address & 0xF000 {
	case 0x0000, 0x1000:
		switch c.Header.Type {
		case Mbc3Ram, Mbc3RamBattery:
			c.RamEnabled = (value & 0x0F) == 0x0A
		case Mbc3TimerBattery:
			c.RtcEnabled = (value & 0x0F) == 0x0A
		case Mbc3TimerRamBattery:
			c.RamEnabled = (value & 0x0F) == 0x0A
			c.RtcEnabled = (value & 0x0F) == 0x0A
		default:
			return
		}
	case 0x2000, 0x3000:
		c.RomBank = int(value)
		if int(c.RomBank)*0x4000 >= len(c.BinaryData) {
			c.RomBank = (int(c.RomBank) % (len(c.BinaryData) / 0x4000))
		}
		if c.RomBank == 0 {
			c.RomBank = 1
		}
	case 0x4000, 0x5000:

		if value >= 0x08 && value <= 0x0c {
			if c.hasRTC && c.RtcEnabled {
				c.RtcRegister = value
				c.RamBank = -1
			}
		} else if value <= 0x03 && c.RamEnabled {
			c.RamBank = int(value) & 0x3
			if len(c.Ram) <= 0 {
				c.RamBank = 0
			} else if int(c.RamBank)*0x2000 >= len(c.Ram) {
				c.RamBank = (int(c.RamBank) % (len(c.Ram) / 0x2000))
			}
		}

	case 0x6000, 0x7000:
		if c.hasRTC {
			if c.RtcLatchFlagValue == 0x00 && value == 0x01 {
				c.UpdateRTC()
				c.Rtc.LatchedSeconds = c.Rtc.Seconds
				c.Rtc.LatchedMinutes = c.Rtc.Minutes
				c.Rtc.LatchedHours = c.Rtc.Hours
				c.Rtc.LatchedDaysLower = c.Rtc.DaysLower
				c.Rtc.LatchedDaysHigherAndControl = c.Rtc.DaysHigherAndControl
			}
			c.RtcLatchFlagValue = value
		}
	case 0xA000, 0xB000:
		if c.RamBank >= 0 {
			if c.RamEnabled {
				c.Ram[c.RamBank*0x2000+int(address&0x1fff)] = value
			}
		} else if c.hasRTC && c.RtcEnabled {
			switch c.RtcRegister {
			case 0x8:
				c.Rtc.Seconds = value & 0x3F
			case 0x9:
				c.Rtc.Minutes = value & 0x3F
			case 0xa:
				c.Rtc.Hours = value & 0x1F
			case 0xb:
				c.Rtc.DaysLower = value
			case 0xc:
				c.Rtc.DaysHigherAndControl = value & 0xC1
			}
		}
	default:
		return
	}
}

func (c *Mbc3Cartridge) UpdateRTC() {
	delta := time.Since(c.RtcLastUpdate)

	if (c.Rtc.DaysHigherAndControl>>6)&0x1 == 0 && (delta >= time.Second) {
		c.RtcLastUpdate = c.RtcLastUpdate.Add(delta)
		var days uint32
		deltaSeconds := int(delta.Seconds())

		c.Rtc.Seconds += uint8(deltaSeconds % 60)
		if c.Rtc.Seconds > 59 {
			c.Rtc.Seconds -= 60
			c.Rtc.Minutes += 1
		}
		deltaSeconds /= 60

		c.Rtc.Minutes += uint8(int(deltaSeconds % 60))
		if c.Rtc.Minutes > 59 {
			c.Rtc.Minutes -= 60
			c.Rtc.Hours += 1
		}
		deltaSeconds /= 60

		c.Rtc.Hours += uint8(int(deltaSeconds % 24))
		if c.Rtc.Hours > 23 {
			c.Rtc.Hours -= 24
			days += 1
		}
		deltaSeconds /= 24

		days += uint32(deltaSeconds)
		days += uint32(c.Rtc.DaysLower)
		days += uint32(uint32(c.Rtc.DaysHigherAndControl&0x1) << 8)
		if days > 511 {
			days = days % 512
			c.Rtc.DaysHigherAndControl |= 1 << 7
		}

		c.Rtc.DaysLower = uint8(days & 0xFF)
		c.Rtc.DaysHigherAndControl = (c.Rtc.DaysHigherAndControl & 0xFE)
		if days > 0xFF {
			c.Rtc.DaysHigherAndControl |= 0x1
		}
	}
}

func (c *Mbc3Cartridge) UpdateComponentsPerCycle(cycles uint16) {
	if c.hasRTC {
		c.RtcRefreshCyclesCount += cycles
		if c.RtcRefreshCyclesCount >= RtcRefreshAfterCycles {
			c.UpdateRTC()
			c.RtcRefreshCyclesCount = 0
		}
	}
}

func (c *Mbc3Cartridge) String() string {
	return fmt.Sprintf("%v %d %d | RTC %v %v %08b %v", c.RamEnabled, c.RamBank, c.RomBank, c.RtcEnabled, c.Rtc, c.Rtc.DaysHigherAndControl, time.Since(c.RtcLastUpdate).Milliseconds())
}

func (c *Mbc3Cartridge) SaveRam(writer io.Writer) error {
	return c.SaveRamBGBFormat(writer)
}

func (c *Mbc3Cartridge) SaveRamBGBFormat(writer io.Writer) error {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.LittleEndian, c.Ram)
	if err != nil {
		return err
	}

	if c.hasRTC {
		var rtc_data = []interface{}{
			uint32(c.Rtc.Seconds),
			uint32(c.Rtc.Minutes),
			uint32(c.Rtc.Hours),
			uint32(c.Rtc.DaysLower),
			uint32(c.Rtc.DaysHigherAndControl),
			uint32(c.Rtc.LatchedSeconds),
			uint32(c.Rtc.LatchedMinutes),
			uint32(c.Rtc.LatchedHours),
			uint32(c.Rtc.LatchedDaysLower),
			uint32(c.Rtc.LatchedDaysHigherAndControl),
			uint64(c.RtcLastUpdate.Unix()),
		}

		for _, v := range rtc_data {
			err = binary.Write(buffer, binary.LittleEndian, v)
			if err != nil {
				return err
			}
		}
	}

	_, err = writer.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

type Mbc3SaveFile struct {
	Ram           []byte
	Rtc           []byte
	RtcLastUpdate time.Time
}

func (c *Mbc3Cartridge) LoadRam(reader io.Reader) error {
	buffer := new(bytes.Buffer)
	nRead, err := io.Copy(buffer, reader)
	if err != nil {
		return err
	}

	expectedSize := int64(c.Header.RamSize)
	if c.hasRTC {
		expectedSize += 48
	}
	if nRead == expectedSize {
		return c.LoadRamBGBFormat(buffer)
	}

	return c.LoadRamOldAXWFormat(buffer)
}

func (c *Mbc3Cartridge) LoadRamBGBFormat(reader io.Reader) error {
	var err error
	c.Ram, err = ioutil.ReadAll(io.LimitReader(reader, int64(c.Header.RamSize)))
	if err != nil {
		log.Printf("error during reading of ram data")
		return err
	}

	if c.hasRTC {
		var rtcData = []*uint8{
			&c.Rtc.Seconds,
			&c.Rtc.Minutes,
			&c.Rtc.Hours,
			&c.Rtc.DaysLower,
			&c.Rtc.DaysHigherAndControl,
			&c.Rtc.LatchedSeconds,
			&c.Rtc.LatchedMinutes,
			&c.Rtc.LatchedHours,
			&c.Rtc.LatchedDaysLower,
			&c.Rtc.LatchedDaysHigherAndControl,
		}

		var data []byte
		for k, v := range rtcData {
			data, err = ioutil.ReadAll(io.LimitReader(reader, 4))
			if err != nil {
				log.Printf("error during reading of rtc data part %v", k)
				return err
			}
			*v = uint8(data[0])
		}

		var rtcTimestamp []byte
		rtcTimestamp, err = ioutil.ReadAll(io.LimitReader(reader, 8))
		if err != nil {
			log.Printf("error during reading of rtc timestamp")
			return err
		}
		c.RtcLastUpdate = time.Unix(int64(binary.LittleEndian.Uint64(rtcTimestamp)), 0)
	}

	return nil
}

func (c *Mbc3Cartridge) LoadRamOldAXWFormat(reader io.Reader) error {
	var err error
	loadFile := Mbc3SaveFile{}

	d := gob.NewDecoder(reader)
	err = d.Decode(&loadFile)
	if err != nil {
		log.Printf("error during decoding of old ram format")
		return err
	}

	c.Ram = loadFile.Ram

	if c.hasRTC {
		if len(loadFile.Rtc) != 5 {
			return errors.New("cartridge has a RTC, but savefile does not include any")
		}

		c.Rtc.Seconds = loadFile.Rtc[0]
		c.Rtc.Minutes = loadFile.Rtc[1]
		c.Rtc.Hours = loadFile.Rtc[2]
		c.Rtc.DaysLower = loadFile.Rtc[3]
		c.Rtc.DaysHigherAndControl = loadFile.Rtc[4]
		c.RtcLastUpdate = loadFile.RtcLastUpdate
	}

	return err
}

func (c *Mbc3Cartridge) GetRamBank() uint8 { return uint8(c.RamBank) }
