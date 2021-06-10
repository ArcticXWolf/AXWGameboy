package internal

import "log"

func instructionCompare(c *Cpu, value byte, otherValue byte) {
	log.Printf("CMP: %x - %x", value, otherValue)
	substraction := otherValue - value
	c.Registers.SetFlagZ(substraction == 0)
	c.Registers.SetFlagN(true)
	c.Registers.SetFlagH((value & 0x0f) > (otherValue & 0x0f))
	c.Registers.SetFlagC(value > otherValue)
}

func instructionAddition(c *Cpu, value byte, otherValue byte, carry bool) byte {
	log.Printf("ADD: %x + %x", value, otherValue)
	doCarry := c.Registers.FlagC() && carry
	sum := int16(value) + int16(otherValue)
	sumHalf := int16(value&0x0f) + int16(otherValue&0x0f)
	if doCarry {
		sum++
		sumHalf++
	}

	c.Registers.SetFlagZ(byte(sum) == 0)
	c.Registers.SetFlagN(false)
	c.Registers.SetFlagH(sumHalf > 0xF)
	c.Registers.SetFlagC(sum > 0xFF)

	return byte(sum)
}

func instructionSubstraction(c *Cpu, value byte, otherValue byte, carry bool) byte {
	log.Printf("SUB: %x + %x", value, otherValue)
	doCarry := c.Registers.FlagC() && carry
	substraction := int16(value) - int16(otherValue)
	substractionHalf := int16(value&0x0f) - int16(otherValue&0x0f)
	if doCarry {
		substraction--
		substractionHalf--
	}

	c.Registers.SetFlagZ(byte(substraction) == 0)
	c.Registers.SetFlagN(true)
	c.Registers.SetFlagH(substractionHalf < 0)
	c.Registers.SetFlagC(substraction < 0)

	return byte(substraction)
}
