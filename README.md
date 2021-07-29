# AXWGameboy

This is a work-in-progress gameboy emulator written in golang. So far most of the regular gameboy features are implemented. CGB support is currently in progress.

Binaries are available for 64-bit Windows, Linux and Android. It should be compatible to macOS, FreeBSD, iOS, WebAssembly and Nintendo Switch™ thanks to cross plattform support by [Ebiten](https://github.com/hajimehoshi/ebiten) and GoMobile. For compilation support please refer to those projects first.

# Features

* CPU emulation
    * z80 specifications
    * All opcodes with instruction timing
* Memory emulation
    * MBC1, MBC3, MBC5 mappers
    * Multiple ROM banks and RAM banks
* PPU
    * Background layer
    * Window layer
    * Sprites
    * CGB Palettes
    * CGB TileAttributes
* APU
    * All 4 channels with all features
    * Big parts are adapted from [Goboy Emulator](https://github.com/Humpheh/goboy)
    * Changed audio backend to the ebiten native one
* Simple debugger
* Hardcoded keybindings

# Usage

Just download the binaries and place them in your path (or in a place of your choosing).

Afterwards you can doubleclick on the binary to open it into a file selector (showing all roms in the current folder).

Otherwise you can use the following flags from commandline:

```
$ axwgameboy-windows-amd64.exe --help
Usage of axwgameboy-windows-amd64.exe:
  -color
        Defaults to true. If set to false, it forces all games to be in non-color-mode. (default true)
  -cui
        Disable normal console output and show console debug gui instead
  -osb
        Enable on-screen-button display.
  -palette string
        For non-color-mode: Specify which color palette shall be used. Currently available: dmg, red, white (default "DMG")
  -rom string
        Set to the path of the rom, which will be used. If not specified, then a file selector will be shown with all roms in the current folder.
  -save string
        Enables RAM persistence and has to contain the path to the desired savefile.
  -serial
        Print bytes of serial output to console as ASCII characters.
  -sound float
        Sets the starting master volume. Specify between 0 and 1. (default 0.5)
```

If a rom is chosen via the file selector, this automatically creates a RAM savefile in the same location.

# Compilation

To compile the source for yourself, just clone the repo and use the provided makefile. You will also need to provide a dump/rom of the GB and GB color boot rom. Place these under `internal/bootroms/dmg_bios.bin` and `internal/bootroms/cgb_bios.bin`.

You might need to install some build dependencies. These are required by ebiten as our game engine. For further instructions (and install help for each OS) please refer to the ebiten documentation.

# Keyboard bindings

Keyboard button | Gameboy button
----------------|---------------
<kbd>Y</kbd>/<kbd>Z</kbd> | A button
<kbd>X</kbd> | B button
<kbd>Space</kbd> | Start button
<kbd>LeftAlt</kbd> | Select button
<kbd>&larr;</kbd> <kbd>&uarr;</kbd> <kbd>&darr;</kbd> <kbd>&rarr;</kbd> | Arrow buttons

Keyboard button | Misc
----------------|---------------
<kbd>LeftShift</kbd> | Speedboost (3x Speed)
<kbd>P</kbd> | Pause game
<kbd>Escape</kbd> | Close emulator
<kbd>+</kbd> | Increase volume
<kbd>-</kbd> | Decrease volume
<kbd>1</kbd> | Toggle sound channel 1
<kbd>2</kbd> | Toggle sound channel 2
<kbd>3</kbd> | Toggle sound channel 3
<kbd>4</kbd> | Toggle sound channel 4
<kbd>d</kbd> | Pause execution and enter debug mode in console
<kbd>t</kbd> | Toggle tilemap0
<kbd>Z</kbd>/<kbd>Y</kbd> | Toggle tilemap1

# Debugger

This emulator contains a simple debugger allowing you to step through the execution of the rom. After initiating a break (see keybindings) execution will pause and some default information about the state of the cpu will be written to console. Afterwards the following commands can be used:

Command | Description
--------|------------
c | Disable debugger and continue execution from here
step | Enter step-mode. Each empty command will advance execution by one cpu instruction
b#### | Set the next breakpoint to address 0x#### and continue execution until the address is reached
gpu | Print some simple information about the state of the GPU
ipl | Only in non-color-mode: Identify the different layers by coloring each layer a different color: BG red, OBJ0 green, OBJ1 blue
t### | Print tile with id ### to console
ws | Print current screen to console
sp | Show current stack and stack pointer
mem | Print the whole memory to console
vmem | Print the current VRAM bank
m####!!!! | Print the memory range from 0x#### to 0x!!!!
cart | Print information about the ROM
q | Exit emulator

# Learning resources

* [Imran Nazars Blog](https://imrannazar.com/GameBoy-Emulation-in-JavaScript)
* [Gameboy CPU Manual](http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf)
* [Goboy Emulator](https://github.com/Humpheh/goboy)
* [Pastraiser Gameboy CPU instruction set](https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
* [tomeks Blog](https://blog.rekawek.eu/2017/02/09/coffee-gb/)
