# AXWGameboy

You can try the emulator in your browser here: [https://arcticxwolf.github.io/AXWGameboy/](https://arcticxwolf.github.io/AXWGameboy/)

This is a work-in-progress gameboy emulator written in golang. So far most of the regular gameboy features are implemented. CGB support is mostly done with a few bugs still existing.

It should be compatible to macOS, FreeBSD, iOS, WebAssembly and Nintendo Switch™ thanks to cross plattform support by [Ebiten](https://github.com/hajimehoshi/ebiten) and GoMobile. For compilation support please refer to those projects first. Main target plattform is webassembly.

# Features

* CPU emulation
    * z80 specifications
    * All opcodes with instruction timing
    * CGB double speed mode
* Memory emulation
    * MBC1, MBC3, MBC5 mappers
    * Multiple ROM banks and RAM banks
    * CGB HDMA and GDMA transfers
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

# Try it online

A copy of the emulator running in your browser is available [here](https://arcticxwolf.github.io/AXWGameboy/).

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
