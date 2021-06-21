# AXWGameboy

This is a work-in-progress gameboy emulator written in golang.

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
* Simple debugger
* Hardcoded keybindings

# Learning resources

* [Imran Nazars Blog](https://imrannazar.com/GameBoy-Emulation-in-JavaScript)
* [Gameboy CPU Manual](http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf)
* [Goboy Emulator](https://github.com/Humpheh/goboy)
* [Pastraiser Gameboy CPU instruction set](https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
* [tomeks Blog](https://blog.rekawek.eu/2017/02/09/coffee-gb/)
