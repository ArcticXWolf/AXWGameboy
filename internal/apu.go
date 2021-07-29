package internal

type Apu interface {
	ReadByte(address uint16) byte
	WriteByte(address uint16, value byte)
	WriteWaveform(address uint16, value byte)
	Buffer(cpuTicks int, speed int, cpuSpeedBoost float64)
	ToggleSoundChannel(channel int)
	ChangeVolume(increment float64)
}
