package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
)

type CheatCodeManager struct {
	genieCodes   []*GameGenieCode
	genieEnabled bool
	sharkCodes   []*GameSharkCode
	sharkEnabled bool
}

type GameGenieCode struct {
	MemoryAddress   uint16
	CompareData     uint8
	ReplacementData uint8
	Checksum        uint8
	Code            string
}

func NewGameGenieCodeFromString(codeString string) (*GameGenieCode, error) {
	var err error
	if matched, _ := regexp.MatchString(`\A[0-9a-fA-F]{3}-[0-9a-fA-F]{3}-[0-9a-fA-F]{3}\z`, codeString); !matched {
		return nil, errors.New("invalid gamegenie code given - regex did not match")
	}

	replacementData, err := strconv.ParseUint(codeString[0:2], 16, 8)
	if err != nil {
		return nil, errors.New("invalid gamegenie code given - could not parse replacement data")
	}
	memoryAddress, err := strconv.ParseUint(codeString[6:7]+codeString[2:3]+codeString[4:6], 16, 16)
	if err != nil {
		return nil, errors.New("invalid gamegenie code given - could not parse memory address")
	}
	memoryAddress = memoryAddress ^ 0xF000

	compareData, err := strconv.ParseUint(codeString[8:9]+codeString[10:11], 16, 8)
	if err != nil {
		return nil, errors.New("invalid gamegenie code given - could not parse compare data")
	}
	compareData = uint64((uint8(compareData>>2) | uint8(compareData<<6)) ^ 0xBA)

	checksum, err := strconv.ParseUint(codeString[9:10], 16, 8)
	if err != nil {
		return nil, errors.New("invalid gamegenie code given - could not parse checksum")
	}

	return &GameGenieCode{
		MemoryAddress:   uint16(memoryAddress),
		CompareData:     uint8(compareData),
		ReplacementData: uint8(replacementData),
		Checksum:        uint8(checksum),
	}, nil
}

func (g *GameGenieCode) String() string {
	return fmt.Sprintf("Code (Address 0x%04x, replace 0x%02x with 0x%02x, CSM 0x%02x)", g.MemoryAddress, g.CompareData, g.ReplacementData, g.Checksum)
}

type GameSharkCode struct {
	MemoryAddress   uint16
	ReplacementData uint8
	Type            uint8
	Code            string
}

func NewGameSharkCodeFromString(codeString string) (*GameSharkCode, error) {
	var err error
	if matched, _ := regexp.MatchString(`\A[0-9a-fA-F]{8}\z`, codeString); !matched {
		return nil, errors.New("invalid gameshark code given - regex did not match")
	}

	replacementData, err := strconv.ParseUint(codeString[2:4], 16, 8)
	if err != nil {
		return nil, errors.New("invalid gameshark code given - could not parse replacement data")
	}
	memoryAddress, err := strconv.ParseUint(codeString[6:8]+codeString[4:6], 16, 16)
	if err != nil {
		return nil, errors.New("invalid gameshark code given - could not parse memory address")
	}

	codeType, err := strconv.ParseUint(codeString[0:2], 16, 8)
	if err != nil {
		return nil, errors.New("invalid gameshark code given - could not parse type")
	}

	return &GameSharkCode{
		MemoryAddress:   uint16(memoryAddress),
		ReplacementData: uint8(replacementData),
		Type:            uint8(codeType),
	}, nil
}

func (g *GameSharkCode) String() string {
	return fmt.Sprintf("Code (Address 0x%04x, replace with 0x%02x, type 0x%02x)", g.MemoryAddress, g.ReplacementData, g.Type)
}

func (c *CheatCodeManager) ReplaceCodeList(reader io.Reader) {
	c.genieCodes = nil
	c.sharkCodes = nil
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if len([]rune(line)) == 11 {
			code, err := NewGameGenieCodeFromString(line)
			if err != nil {
				log.Printf("error parsing code %s - %s \n", line, err)
				continue
			}
			c.genieCodes = append(c.genieCodes, code)
			log.Printf("added genie code %v", code.String())
			continue
		} else if len([]rune(line)) == 8 {
			code, err := NewGameSharkCodeFromString(line)
			if err != nil {
				log.Printf("error parsing code %s - %s \n", line, err)
				continue
			}
			c.sharkCodes = append(c.sharkCodes, code)
			log.Printf("added shark code %v", code.String())
			continue
		}

		log.Printf("error parsing code %s - no code type detected \n", line)
	}

}

func (c *CheatCodeManager) ApplyGameGenie(address uint16, originalValue uint8) uint8 {
	if !c.genieEnabled {
		return originalValue
	}

	for _, v := range c.genieCodes {
		if v.MemoryAddress != address {
			continue
		}

		if v.CompareData != originalValue {
			log.Printf("not replacing value at %x with %x (old %x != expected %x) \n", v.MemoryAddress, v.ReplacementData, originalValue, v.CompareData)
			continue
		}

		log.Printf("replacing value at %x with %x (old %x) \n", v.MemoryAddress, v.ReplacementData, v.CompareData)
		return v.ReplacementData
	}

	return originalValue
}

func (c *CheatCodeManager) ApplyGameShark(gb *Gameboy) {
	if !c.sharkEnabled {
		return
	}

	for _, v := range c.sharkCodes {
		if v.Type != 0x01 {
			continue
		}

		if v.MemoryAddress < 0xA000 || v.MemoryAddress >= 0xE000 {
			continue
		}

		gb.Memory.WriteByte(v.MemoryAddress, v.ReplacementData)
	}
}

func (c *CheatCodeManager) GetCodeList() string {
	var codeList string

	for _, v := range c.genieCodes {
		codeList = fmt.Sprintf("%s%s\n", codeList, v.Code)
	}
	for _, v := range c.sharkCodes {
		codeList = fmt.Sprintf("%s%s\n", codeList, v.Code)
	}

	return codeList
}
