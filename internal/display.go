package internal

type DisplayProvider interface {
	Render(gb *Gameboy)
}
