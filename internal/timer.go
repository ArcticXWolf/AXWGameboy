package internal

type Timer struct {
	lastSeenCpuCyles  int
	dividerValue      uint8
	counterForCounter uint16
	counterValue      uint16
	moduloValue       uint8
	controlFlag       uint8
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) getSpeed() uint8 {
	return t.controlFlag & 0x03
}

func (t *Timer) setSpeed(speed uint8) {
	t.controlFlag &= ^uint8(0x03)
	t.controlFlag |= speed & 0x03
}

func (t *Timer) isCounterRunning() bool {
	return t.controlFlag&0x4 > 0
}

func (t *Timer) Update(gb *Gameboy) {
	if gb.Cpu.ClockCycles-t.lastSeenCpuCyles > 16 {
		t.dividerValue++
		t.lastSeenCpuCyles += 16

		if t.isCounterRunning() {
			t.counterForCounter++
			var overflowThreshold uint16
			switch t.getSpeed() {
			case 0:
				overflowThreshold = 64
			case 1:
				overflowThreshold = 1
			case 2:
				overflowThreshold = 4
			case 3:
				overflowThreshold = 16
			}

			if t.counterForCounter >= overflowThreshold {
				t.counterForCounter = 0
				t.counterValue++

				if t.counterValue > 0xFF {
					t.counterValue = uint16(t.moduloValue)
					gb.Memory.GetInterruptFlags().TriggeredFlags |= 0x04
				}
			}
		}
	}
}
