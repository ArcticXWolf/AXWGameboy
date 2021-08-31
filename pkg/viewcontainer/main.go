package viewcontainer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/pkg/gameview"
)

type ViewContainerStage int

const (
	GameSelection ViewContainerStage = iota
	GameRunning
	SettingsOpen
)

type AXWGameboyViewContainer struct {
	GameView     *gameview.AXWGameboyEbitenGameView
	currentStage ViewContainerStage
}

func NewAXWGameboyViewContainer() *AXWGameboyViewContainer {
	ag := &AXWGameboyViewContainer{
		currentStage: GameSelection,
	}

	ag.createGameSelector()

	return ag
}

func (a *AXWGameboyViewContainer) Update() error {
	if a.currentStage == GameRunning {
		return a.GameView.Update()
	}
	return nil
}

func (a *AXWGameboyViewContainer) Draw(screen *ebiten.Image) {
	if a.currentStage == GameRunning {
		a.GameView.Draw(screen)
	}
}

func (a *AXWGameboyViewContainer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if a.currentStage == GameRunning {
		return a.GameView.Layout(outsideWidth, outsideHeight)
	}
	return 1, 1
}

func (a *AXWGameboyViewContainer) GetExitResult() []byte {
	return []byte{}
}
