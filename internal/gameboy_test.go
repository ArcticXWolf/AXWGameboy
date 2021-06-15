package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBlarggRoms(t *testing.T) {
	executionCycleCount := 20000000
	romPath := "../roms/blargg/cpu_instrs"
	filepath.Walk(romPath, func(path string, _ os.FileInfo, _ error) error {
		if filepath.Ext(path) == ".gb" {
			name := path[len(romPath)+1 : len(path)-3]
			t.Run(
				name,
				func(t *testing.T) {
					result := ""
					options := &GameboyOptions{
						RomPath: path,
						SerialOutputFunction: func(character byte) {
							result += string(character)
						},
						Headless: true,
					}
					gb, err := NewGameboy(options)
					if err != nil {
						t.Error(err)
					}

					cyclecount := 0
					for cycle := 0; cycle < executionCycleCount; cycle++ {
						cycles := gb.Cpu.Tick(gb)
						gb.Gpu.Update(gb, cycles)
						cyclecount++

						if strings.Contains(result, "Passed") || strings.Contains(result, "Failed") {
							break
						}
					}

					if strings.Contains(result, "Passed") {
						t.Logf("Testrom (%s) did PASS after %d executions: %v", path, cyclecount, result)
					} else {
						if cyclecount >= executionCycleCount {
							t.Errorf("Testrom (%s) timed out after %d executions: %v", path, cyclecount, result)
						} else {
							t.Errorf("Testrom (%s) did not pass after %d executions: %v", path, cyclecount, result)
						}
					}
				},
			)
		}
		return nil
	})
}

func TestMooneyeRoms(t *testing.T) {
	executionCycleCount := 20000000
	romPath := "../roms/mooneye/acceptance"
	filepath.Walk(romPath, func(path string, _ os.FileInfo, _ error) error {
		if filepath.Ext(path) == ".gb" {
			name := path[len(romPath)+1 : len(path)-3]
			t.Run(
				name,
				func(t *testing.T) {
					result := make([]byte, 0)
					options := &GameboyOptions{
						RomPath: path,
						SerialOutputFunction: func(character byte) {
							result = append(result, character)
						},
						Headless: true,
					}
					gb, err := NewGameboy(options)
					if err != nil {
						t.Error(err)
					}

					cyclecount := 0
					for cycle := 0; cycle < executionCycleCount; cycle++ {
						cycles := gb.Cpu.Tick(gb)
						gb.Gpu.Update(gb, cycles)
						cyclecount++

						if len(result) >= 6 {
							break
						}
					}

					if len(result) == 6 &&
						result[0] == 0x03 &&
						result[1] == 0x05 &&
						result[2] == 0x08 &&
						result[3] == 0x0d &&
						result[4] == 0x15 &&
						result[5] == 0x22 {
						t.Logf("Testrom (%s) did PASS after %d executions: %v", path, cyclecount, result)
					} else {
						if cyclecount >= executionCycleCount {
							t.Errorf("Testrom (%s) timed out after %d executions: %v", path, cyclecount, result)
						} else {
							t.Errorf("Testrom (%s) did not pass after %d executions: %v", path, cyclecount, result)
						}
					}
				},
			)
		}
		return nil
	})
}
