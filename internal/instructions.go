package internal

func instructionCompare(gb *Gameboy, value byte, otherValue byte) {
	// log.Printf("CMP: %x - %x", value, otherValue)
	substraction := otherValue - value
	gb.Cpu.Registers.SetFlagZ(substraction == 0)
	gb.Cpu.Registers.SetFlagN(true)
	gb.Cpu.Registers.SetFlagH((value & 0x0f) > (otherValue & 0x0f))
	gb.Cpu.Registers.SetFlagC(value > otherValue)
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
