package ebitenprovider

import (
	"fmt"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"go.janniklasrichter.de/axwgameboy/internal"
	"go.janniklasrichter.de/axwgameboy/pkg/gameview"
	"go.janniklasrichter.de/axwgameboy/pkg/menuview"
)

type View interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
	GetExitResult() []byte
}

type AXWGameboyEbitenGame struct {
	views          []View
	GameboyOptions *internal.GameboyOptions
	activeView     int
	requestView    int
}

func NewAXWGameboyEbitenGame(options *internal.GameboyOptions) *AXWGameboyEbitenGame {
	views := []View{}

	if options.RomPath != "" {
		views = append(views, gameview.NewAXWGameboyEbitenGameView(options))
	} else {
		views = append(views, menuview.NewAXWGameboyEbitenMenuView(options))
	}

	ag := &AXWGameboyEbitenGame{
		views:          views,
		GameboyOptions: options,
		activeView:     0,
		requestView:    0,
	}

	return ag
}

func (a *AXWGameboyEbitenGame) Update() error {
	err := a.views[a.activeView].Update()

	if err == menuview.ErrGoToGame {
		a.GameboyOptions.RomPath = string(a.views[a.activeView].GetExitResult())
		a.GameboyOptions.SavePath = filepath.Join(filepath.Dir(a.GameboyOptions.RomPath), fmt.Sprintf("%s.sav", filepath.Base(a.GameboyOptions.RomPath)))
		a.views = append(a.views, gameview.NewAXWGameboyEbitenGameView(a.GameboyOptions))
		a.requestView = len(a.views) - 1
		err = nil
	}

	return err
}

func (a *AXWGameboyEbitenGame) Draw(screen *ebiten.Image) {
	a.views[a.activeView].Draw(screen)
}

func (a *AXWGameboyEbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if a.requestView != a.activeView {
		a.activeView = a.requestView
	}
	return a.views[a.activeView].Layout(outsideWidth, outsideHeight)
}
