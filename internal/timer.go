package internal

// Cyclecorrect implementation thanks to
// /u/gogos-venge
// https://www.reddit.com/r/EmuDev/comments/acsu62/question_regarding_rapid_togglegb_test_rom_game/

type Timer struct {
	gb *Gameboy

	DividerValue uint16
	CounterValue uint8
	CounterCarry bool
	ModuloValue  uint8
	ControlFlag  uint8

	lastSeenCpuCycles int
	overflowing       bool
	releaseOverflow   bool
	fallingEdgeDelay  bool
}

func NewTimer(gb *Gameboy) *Timer {
	return &Timer{gb: gb}
}

func (t *Timer) isTimerEnabled() bool {
	return t.ControlFlag&0x4 > 0
}

func (t *Timer) getMultiplexerMask() uint16 {
	switch t.ControlFlag & 0x3 {
	case 0:
		return 0x200
	case 3:
		return 0x80
	case 2:
		return 0x20
	case 1:
		return 0x8
	}
	return 0x0
}

func (t *Timer) ReadByte(address uint16) (result uint8) {
	switch address {
	case 0xFF04:
		return uint8(t.DividerValue >> 8)
	case 0xFF05:
		return t.CounterValue
	case 0xFF06:
		return t.ModuloValue
	case 0xFF07:
		return t.ControlFlag
	default:
		return 0x00
	}
}

func (t *Timer) WriteByte(address uint16, value uint8) {
	switch address {
	case 0xFF04:
		t.DividerValue = 0
		t.Update(t.gb)
	case 0xFF05:
		if t.releaseOverflow {
			return
		}
		t.CounterValue = value
		t.CounterCarry = false
		t.overflowing = false
		t.releaseOverflow = false
	case 0xFF06:
		if t.releaseOverflow {
			t.CounterValue = value
		}
		t.ModuloValue = value
	case 0xFF07:
		t.ControlFlag = value & 0x7
		t.Update(t.gb)
		t.Update(t.gb)
	default:
	}
}

func (t *Timer) Update(gb *Gameboy) {
	cycles := gb.Cpu.ClockCycles - t.lastSeenCpuCycles
	t.lastSeenCpuCycles = gb.Cpu.ClockCycles

	for i := 0; i*4 < cycles; i++ {
		t.DividerValue += uint16(4)

		var timersignal bool = (t.DividerValue&t.getMultiplexerMask()) == t.getMultiplexerMask() && t.isTimerEnabled()

		if t.releaseOverflow {
			t.overflowing = false
			t.releaseOverflow = false
		}

		if t.overflowing {
			t.CounterValue = t.ModuloValue
			gb.Memory.GetInterruptFlags().TriggeredFlags |= 0x4
			t.CounterCarry = false
			t.releaseOverflow = true
		}

		if t.detectFallingEdge(timersignal) {
			t.CounterValue++
			if t.CounterValue == 0x0 && t.CounterCarry {
				t.overflowing = true
			} else if t.CounterValue == 0xFF {
				t.CounterCarry = true
			}
		}
	}
}

func (t *Timer) detectFallingEdge(signal bool) bool {
	result := !signal && t.fallingEdgeDelay
	t.fallingEdgeDelay = signal
	return result
}
