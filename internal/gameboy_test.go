package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func executeRomUntilCompletionOrTimeout(t *testing.T, romPath string, maxExecutionCycles int, completionFunction func([]byte) (bool, bool)) {
	name := strings.TrimSuffix(filepath.Base(romPath), filepath.Ext(romPath))
	t.Run(
		name,
		func(t *testing.T) {
			result := make([]byte, 0)
			options := &GameboyOptions{
				RomPath: romPath,
				SerialOutputFunction: func(character byte) {
					result = append(result, character)
				},
				Headless: true,
			}
			gb, err := NewGameboy(options)
			if err != nil {
				t.Error(err)
			}

			var complete, success bool
			cyclecount := 0
			for cycle := 0; cycle < maxExecutionCycles; cycle++ {
				cycles := gb.Cpu.Tick(gb)
				gb.Gpu.Update(gb, cycles)
				cyclecount++

				if complete, success = completionFunction(result); complete {
					break
				}
			}

			if success {
				t.Logf("PASS: #%015d %s: %v", cyclecount, romPath, result)
			} else {
				if !complete {
					t.Errorf("TIME: #%015d %s: %v", cyclecount, romPath, result)
				} else {
					t.Errorf("ERRO: #%015d %s: %v", cyclecount, romPath, result)
				}
			}
		},
	)
}

func executeRomDirectories(t *testing.T, romDirectories []string, maxExecutionCycles int, completionFunction func([]byte) (bool, bool)) {
	for _, romDirectory := range romDirectories {
		filepath.Walk(romDirectory, func(path string, _ os.FileInfo, _ error) error {
			if filepath.Ext(path) == ".gb" {
				executeRomUntilCompletionOrTimeout(t, path, maxExecutionCycles, completionFunction)
			}
			return nil
		})
	}
}

func TestBlarggCPUInstrsRoms(t *testing.T) {
	maxExecutionCycles := 20000000
	romDirectories := []string{"../roms/blargg/cpu_instrs"}
	completionFunc := func(result []byte) (bool, bool) {
		return strings.Contains(string(result), "Passed") || strings.Contains(string(result), "Failed"), strings.Contains(string(result), "Passed")
	}
	executeRomDirectories(t, romDirectories, maxExecutionCycles, completionFunc)
}

func TestBlarggInstrTimingRoms(t *testing.T) {
	maxExecutionCycles := 20000000
	romPath := "../roms/blargg/instr_timing.gb"
	completionFunc := func(result []byte) (bool, bool) {
		return strings.Contains(string(result), "Passed") || strings.Contains(string(result), "Failed"), strings.Contains(string(result), "Passed")
	}
	executeRomUntilCompletionOrTimeout(t, romPath, maxExecutionCycles, completionFunc)
}

func TestMooneyeRoms(t *testing.T) {
	maxExecutionCycles := 20000000
	romDirectories := []string{
		"../roms/mooneye/acceptance",
		"../roms/mooneye/emulator-only/mbc1",
		"../roms/mooneye/emulator-only/mbc5",
		// "../roms/mooneye/misc",
	}
	completionFunc := func(result []byte) (bool, bool) {
		return len(result) >= 6, len(result) == 6 &&
			result[0] == 0x03 &&
			result[1] == 0x05 &&
			result[2] == 0x08 &&
			result[3] == 0x0d &&
			result[4] == 0x15 &&
			result[5] == 0x22
	}
	executeRomDirectories(t, romDirectories, maxExecutionCycles, completionFunc)
}
