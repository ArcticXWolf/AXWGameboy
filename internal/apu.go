package internal

type Apu interface {
	ReadByte(address uint16) byte
	WriteByte(address uint16, value byte)
	WriteWaveform(address uint16, value byte)
	Buffer(cpuTicks int, speed int)
	ToggleSoundChannel(channel int)
}
