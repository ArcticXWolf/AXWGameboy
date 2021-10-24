package internal

func instructionCompare(gb *Gameboy, value byte, otherValue byte) {
	// log.Printf("CMP: %x - %x", value, otherValue)
	substraction := otherValue - value
	gb.Cpu.Registers.SetFlagZ(substraction == 0)
	gb.Cpu.Registers.SetFlagN(true)
	gb.Cpu.Registers.SetFlagH((value & 0x0f) > (otherValue & 0x0f))
	gb.Cpu.Registers.SetFlagC(value > otherValue)
}

func instructionIncrement(gb *Gameboy, value byte) byte {
	result := value + 1
	gb.Cpu.Registers.SetFlagZ(result == 0x00)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH((value&0xF)+1 > 0xF)

	return byte(result)
}

func instructionAddition(gb *Gameboy, value byte, otherValue byte, carry bool) byte {
	// log.Printf("ADD: %x + %x", value, otherValue)
	doCarry := gb.Cpu.Registers.FlagC() && carry
	sum := int16(value) + int16(otherValue)
	sumHalf := int16(value&0x0f) + int16(otherValue&0x0f)
	if doCarry {
		sum++
		sumHalf++
	}

	gb.Cpu.Registers.SetFlagZ(byte(sum) == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(sumHalf > 0xF)
	gb.Cpu.Registers.SetFlagC(sum > 0xFF)

	return byte(sum)
}

func instructionAddition16(gb *Gameboy, value uint16, otherValue uint16) uint16 {
	// log.Printf("ADD: %x + %x", value, otherValue)
	sum := int32(value) + int32(otherValue)

	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(int32(value&0xFFF) > (sum & 0xFFF))
	gb.Cpu.Registers.SetFlagC(sum > 0xFFFF)

	return uint16(sum)
}

func instructionSubstraction(gb *Gameboy, value byte, otherValue byte, carry bool) byte {
	// log.Printf("SUB: %x + %x", value, otherValue)
	doCarry := gb.Cpu.Registers.FlagC() && carry
	substraction := int16(value) - int16(otherValue)
	substractionHalf := int16(value&0x0f) - int16(otherValue&0x0f)
	if doCarry {
		substraction--
		substractionHalf--
	}

	gb.Cpu.Registers.SetFlagZ(byte(substraction) == 0)
	gb.Cpu.Registers.SetFlagN(true)
	gb.Cpu.Registers.SetFlagH(substractionHalf < 0)
	gb.Cpu.Registers.SetFlagC(substraction < 0)

	return byte(substraction)
}

func instructionTestBit(gb *Gameboy, value byte, bitIndex uint8) {
	gb.Cpu.Registers.SetFlagZ((value>>bitIndex)&0x1 == 0x0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(true)
}

func instructionSetBit(gb *Gameboy, value byte, bitIndex uint8) byte {
	return value | (0x1 << bitIndex)
}

func instructionResetBit(gb *Gameboy, value byte, bitIndex uint8) byte {
	return value & ^(0x1 << bitIndex)
}

func instructionSwap(gb *Gameboy, value byte) byte {
	result := ((value << 4) & 0xF0) | (value >> 4)
	gb.Cpu.Registers.SetFlagZ(result == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(false)
	gb.Cpu.Registers.SetFlagC(false)

	return result
}

func instructionCBRL(gb *Gameboy, value byte) byte {
	newCarry := value >> 7
	rotation := (value << 1) & 0xFF
	if gb.Cpu.Registers.FlagC() {
		rotation |= 1
	}

	gb.Cpu.Registers.SetFlagZ(rotation == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(false)
	gb.Cpu.Registers.SetFlagC(newCarry == 1)
	return rotation
}

func instructionCBRLC(gb *Gameboy, value byte) byte {
	newCarry := value >> 7
	rotation := (value<<1)&0xFF | newCarry

	gb.Cpu.Registers.SetFlagZ(rotation == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(false)
	gb.Cpu.Registers.SetFlagC(newCarry == 1)
	return rotation
}

func instructionCBRR(gb *Gameboy, value byte) byte {
	newCarry := value & 0x1
	rotation := (value >> 1)
	if gb.Cpu.Registers.FlagC() {
		rotation |= (1 << 7)
	}

	gb.Cpu.Registers.SetFlagZ(rotation == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(false)
	gb.Cpu.Registers.SetFlagC(newCarry == 1)
	return rotation
}

func instructionCBRRC(gb *Gameboy, value byte) byte {
	newCarry := value & 0x1
	rotation := (value >> 1) | (newCarry << 7)

	gb.Cpu.Registers.SetFlagZ(rotation == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(false)
	gb.Cpu.Registers.SetFlagC(newCarry == 1)
	return rotation
}

func instructionCBSLA(gb *Gameboy, value byte) byte {
	newCarry := value >> 7
	rotation := (value << 1) & 0xFF

	gb.Cpu.Registers.SetFlagZ(rotation == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(false)
	gb.Cpu.Registers.SetFlagC(newCarry == 1)
	return rotation
}

func instructionCBSRA(gb *Gameboy, value byte) byte {
	rotation := (value & 0x80) | (value >> 1)

	gb.Cpu.Registers.SetFlagZ(rotation == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(false)
	gb.Cpu.Registers.SetFlagC(value&1 == 1)
	return rotation
}

func instructionCBSRL(gb *Gameboy, value byte) byte {
	newCarry := value & 1
	rotation := (value >> 1)

	gb.Cpu.Registers.SetFlagZ(rotation == 0)
	gb.Cpu.Registers.SetFlagN(false)
	gb.Cpu.Registers.SetFlagH(false)
	gb.Cpu.Registers.SetFlagC(newCarry == 1)
	return rotation
}

func instructionInterrupt(gb *Gameboy, interruptIndex int) bool {
	if !gb.Cpu.Registers.Ime && gb.Halted {
		gb.Halted = false
		return false
	}
	gb.Halted = false
	gb.Cpu.Registers.Ime = false
	gb.Memory.GetInterruptFlags().TriggeredFlags &= ^(1 << interruptIndex)

	gb.Cpu.Registers.Sp -= 2
	gb.Memory.WriteWord(gb.Cpu.Registers.Sp, gb.Cpu.Registers.Pc)

	switch interruptIndex {
	case 0:
		gb.Cpu.Registers.Pc = 0x40
		gb.CheatCodeManager.ApplyGameShark(gb)
	case 1:
		gb.Cpu.Registers.Pc = 0x48
	case 2:
		gb.Cpu.Registers.Pc = 0x50
	case 3:
		gb.Cpu.Registers.Pc = 0x58
	case 4:
		gb.Cpu.Registers.Pc = 0x60
	}

	return true
}
