package internal

type MemoryDevice interface {
	ReadByte(address uint16) uint8
	ReadWord(address uint16) uint16
	WriteByte(address uint16, value uint8)
	WriteWord(address uint16, value uint16)
	GetInterruptFlags() *InterruptFlags
	String() string
}
