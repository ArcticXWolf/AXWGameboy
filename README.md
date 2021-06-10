# AXWGameboy

This is a work-in-progress gameboy emulator written in golang.

# Features

* CPU emulation
    * z80 specifications
    * Some opcodes implemented (171/245 and 17/256)
* Memory emulation
    * Only a simple implementation of the MMU layout
    * No I/O, memory bank switching, etc

# Working ROMs

* GB Bootstrap ROM

# Learning resources

* [Imran Nazars Blog](https://imrannazar.com/GameBoy-Emulation-in-JavaScript)
* [Gameboy CPU Manual](http://marc.rawer.de/Gameboy/Docs/GBCPUman.pdf)
* [Goboy Emulator](https://github.com/Humpheh/goboy)
* [Pastraiser Gameboy CPU instruction set](https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)
* [tomeks Blog](https://blog.rekawek.eu/2017/02/09/coffee-gb/)
