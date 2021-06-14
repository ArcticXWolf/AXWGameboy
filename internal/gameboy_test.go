package internal

import (
	"fmt"
	"strings"
	"testing"
)

func TestBlarggRoms(t *testing.T) {
	type romtest struct {
		filename            string
		executionCycleCount int
	}
	var blarggRoms [11]romtest = [11]romtest{
		{"../roms/01.gb", 4000000},
		{"../roms/02.gb", 4000000},
		{"../roms/03.gb", 4000000},
		{"../roms/04.gb", 4000000},
		{"../roms/05.gb", 5000000},
		{"../roms/06.gb", 4000000},
		{"../roms/07.gb", 4000000},
		{"../roms/08.gb", 6000000},
		{"../roms/09.gb", 7000000},
		{"../roms/10.gb", 10000000},
		{"../roms/11.gb", 10000000},
	}
	for i := 0; i < len(blarggRoms); i++ {
		t.Run(
			fmt.Sprintf("Blargg ROM %02d", i+1),
			func(t *testing.T) {
				result := ""
				options := &GameboyOptions{
					RomPath: blarggRoms[i].filename,
					SerialOutputFunction: func(character byte) {
						result += string(character)
					},
					Headless: true,
				}
				gb, err := NewGameboy(options)
				if err != nil {
					t.Error(err)
				}

				for cycle := 0; cycle < blarggRoms[i].executionCycleCount; cycle++ {
					cycles := gb.Cpu.Tick(gb)
					gb.Gpu.Update(gb, cycles)
				}

				if !strings.Contains(result, "Passed") {
					t.Errorf("Blargg Testrom #%02d did not pass after %d executions: %s", i+1, blarggRoms[i].executionCycleCount, result)
				}
			},
		)
	}
}
