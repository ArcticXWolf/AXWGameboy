package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/pkg/viewcontainer"
)

var (
	version = "dev"
	date    = "dev"
	commit  = "dev"
)

func main() {
	log.Printf("AXWGameboy | Version %v | Builddate %v | Commit %v", version, date, commit)
	startGame()
}

func startGame() {
	ebitenGame := viewcontainer.NewAXWGameboyViewContainer()
	ebiten.SetWindowResizable(true)
	ebiten.RunGame(ebitenGame)
}
