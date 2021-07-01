package internal

type Button byte

type InputProvider interface {
	HandleInput(gb *Gameboy)
}

type Inputs struct {
	inputRow        [2]byte
	inputColumn     byte
	ButtonsPressed  []Button
	ButtonsReleased []Button
}

const (
	ButtonA Button = iota
	ButtonB
	ButtonSelect
	ButtonStart
	ButtonRight
	ButtonLeft
	ButtonUp
	ButtonDown
)

func NewInputs() *Inputs {
	return &Inputs{
		inputRow: [2]byte{
			0xff,
			0xff,
		},
		ButtonsPressed:  make([]Button, 4),
		ButtonsReleased: make([]Button, 4),
	}
}

func (i *Inputs) ReadByte(address uint16) (result uint8) {
	switch i.inputColumn {
	case 0x10:
		return i.inputRow[0]
	case 0x20:
		return i.inputRow[1]
	default:
		return 0xFF
	}
}

func (i *Inputs) WriteByte(address uint16, value uint8) {
	i.inputColumn = value & 0x30
}

func (i *Inputs) HandleInput(gb *Gameboy) {
	for _, button := range i.ButtonsPressed {
		if button > 3 {
			i.inputRow[1] &= ^(0x1 << (button - 4))
		} else {
			i.inputRow[0] &= ^(0x1 << button)
		}
		gb.Memory.GetInterruptFlags().TriggeredFlags |= (1 << 4)
	}
	for _, button := range i.ButtonsReleased {
		if button > 3 {
			i.inputRow[1] |= (0x1 << (button - 4))
		} else {
			i.inputRow[0] |= (0x1 << button)
		}
	}
}

func (i *Inputs) ClearButtonList() {
	i.ButtonsPressed = nil
	i.ButtonsReleased = nil
}
