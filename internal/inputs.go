package internal

type Button byte

type Inputs struct {
	inputRow        [2]byte
	inputColumn     byte
	buttonsPressed  []Button
	buttonsReleased []Button
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
		buttonsPressed:  make([]Button, 4),
		buttonsReleased: make([]Button, 4),
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
	for _, button := range i.buttonsPressed {
		if button > 3 {
			i.inputRow[1] &= ^(0x1 << (button - 4))
		} else {
			i.inputRow[0] &= ^(0x1 << button)
		}
		gb.Memory.GetInterruptFlags().TriggeredFlags |= (1 << 4)
	}
	for _, button := range i.buttonsReleased {
		if button > 3 {
			i.inputRow[1] |= (0x1 << (button - 4))
		} else {
			i.inputRow[0] |= (0x1 << button)
		}
	}
}

func (i *Inputs) ClearButtonList() {
	i.buttonsPressed = nil
	i.buttonsReleased = nil
}
