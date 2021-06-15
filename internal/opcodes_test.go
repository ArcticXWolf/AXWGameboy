package internal

import "testing"

func Test_amountOpcodesWritten(t *testing.T) {
	implemented := 0
	for k, v := range Opcodes {
		if v != nil {
			implemented++
		} else {
			t.Logf("Missing: 0x%02x", k)
		}
	}
	t.Logf("Amount of implemented opcodes: %03d / %03d", implemented, 256)

	fillOpcodesCb(false)
	implemented = 0
	for k, v := range OpcodesCb {
		if v != nil {
			implemented++
		} else {
			t.Logf("Missing: 0x%02x", k)
		}
	}
	t.Logf("Amount of implemented CB opcodes: %03d / %03d", implemented, 256)
}
