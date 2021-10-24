package internal

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func saveCurrentDisplayToImage(gb *Gameboy, filename string) {
	image := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{ScreenWidth, ScreenHeight}})
	for y := 0; y < ScreenHeight; y++ {
		for x := 0; x < ScreenWidth; x++ {
			image.SetRGBA(x, y, color.RGBA{
				R: gb.ReadyToRender[x][y][0],
				G: gb.ReadyToRender[x][y][1],
				B: gb.ReadyToRender[x][y][2],
				A: 255,
			})
		}
	}
	f, _ := os.Create(filename)
	defer f.Close()
	png.Encode(f, image)
}

func executeRomUntilCompletionOrTimeout(t *testing.T, romPath string, romReader io.Reader, maxExecutionCycles int, completionFunction func([]byte) (bool, bool)) {
	name := strings.TrimSuffix(filepath.Base(romPath), filepath.Ext(romPath))
	t.Run(
		name,
		func(t *testing.T) {
			result := make([]byte, 0)
			options := &GameboyOptions{
				RomReader: romReader,
				SerialOutputFunction: func(character byte) {
					result = append(result, character)
				},
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

			imagepath := filepath.Join("..", "build", "testresults", filepath.Base(filepath.Dir(romPath)))
			os.MkdirAll(imagepath, os.ModePerm)

			saveCurrentDisplayToImage(gb, filepath.Join(imagepath, fmt.Sprintf("%s.png", name)))
		},
	)
}

func executeRomDirectories(t *testing.T, romDirectories []string, maxExecutionCycles int, completionFunction func([]byte) (bool, bool)) {
	for _, romDirectory := range romDirectories {
		filepath.Walk(romDirectory, func(path string, _ os.FileInfo, _ error) error {
			if filepath.Ext(path) == ".gb" {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()
				executeRomUntilCompletionOrTimeout(t, path, file, maxExecutionCycles, completionFunction)
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
	file, err := os.Open(romPath)
	if err != nil {
		return
	}
	defer file.Close()

	executeRomUntilCompletionOrTimeout(t, romPath, file, maxExecutionCycles, completionFunc)
}

func TestMooneyeRoms(t *testing.T) {
	maxExecutionCycles := 20000000
	romDirectories := []string{
		"../roms/mooneye/acceptance",
		"../roms/mooneye/emulator-only/mbc1",
		"../roms/mooneye/emulator-only/mbc5",
		"../roms/mooneye/misc",
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
