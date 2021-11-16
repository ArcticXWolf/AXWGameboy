# AXWGameboy

You can try the emulator in your browser here: [https://arcticxwolf.github.io/AXWGameboy/](https://arcticxwolf.github.io/AXWGameboy/)

This is a work-in-progress gameboy emulator written in golang. So far most of the regular gameboy features are implemented. CGB support is mostly done with a few bugs still existing.

It should be compatible (with some additional work) to macOS, FreeBSD, iOS, WebAssembly and Nintendo Switch™ thanks to cross plattform support by [Ebiten](https://github.com/hajimehoshi/ebiten) and GoMobile. For compilation support please refer to those projects first. Main target plattform is webassembly.

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
* Cheat support with GameShark & Gamegenie
* Open-Source bootroms adapted from [Sameboy Emulator](https://github.com/LIJI32/SameBoy)

# Try it online

A copy of the emulator running in your browser is available [here](https://arcticxwolf.github.io/AXWGameboy/).

# Compilation

You might need to install some build dependencies. These are required by ebiten as our game engine. For further instructions (and install help for each OS) please refer to the ebiten documentation.

Afterwards you can build the bootroms with `make bootroms` and then use `make all` to build the emulator + start a webserver serving it on localhost:8008.

# Keyboard bindings

Keyboard button | Gameboy button
----------------|---------------
<kbd>A</kbd> | A button
<kbd>S</kbd> | B button
<kbd>Space</kbd> | Start button
<kbd>LeftAlt</kbd> | Select button
<kbd>&larr;</kbd> <kbd>&uarr;</kbd> <kbd>&darr;</kbd> <kbd>&rarr;</kbd> | Arrow buttons

Keyboard button | Misc
----------------|---------------
<kbd>LeftShift</kbd> | Speedboost (3x Speed)
<kbd>P</kbd> | Pause game
<kbd>+</kbd> | Increase volume
<kbd>-</kbd> | Decrease volume
<kbd>1</kbd> | Toggle sound channel 1
<kbd>2</kbd> | Toggle sound channel 2
<kbd>3</kbd> | Toggle sound channel 3
<kbd>4</kbd> | Toggle sound channel 4
<kbd>d</kbd> | Show debug view

# Learning resources

* [Imran Nazars Blog](https://imrannazar.com/GameBoy-Emulation-in-JavaScript)
* [Gameboy CPU Manual](http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf)
* [Goboy Emulator](https://github.com/Humpheh/goboy)
* [Pastraiser Gameboy CPU instruction set](https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
* [tomeks Blog](https://blog.rekawek.eu/2017/02/09/coffee-gb/)
* [Sameboy Emulator](https://github.com/LIJI32/SameBoy)