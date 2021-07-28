package menuview

import (
	"errors"
	"image/color"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"go.janniklasrichter.de/axwgameboy/internal"
	imagefont "golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type AXWGameboyEbitenMenuView struct {
	files          []string
	isRomfileGiven bool
	fontFace       imagefont.Face
	SelectedPath   string
}

var ErrGoToGame = errors.New("switch to gameview")

func NewAXWGameboyEbitenMenuView(options *internal.GameboyOptions) *AXWGameboyEbitenMenuView {
	files := make([]string, 0)
	filepath.Walk(GetRootPathForCurrentPlattform(), func(path string, info fs.FileInfo, err error) error {
		if !strings.Contains(path, "blargg") && !strings.Contains(path, "mooneye") {
			if filepath.Ext(path) == ".gb" || filepath.Ext(path) == ".gbc" {
				files = append(files, filepath.Clean(path))
			}
		}

		return nil
	})

	font, _ := truetype.Parse(goregular.TTF)
	fontface := truetype.NewFace(font, &truetype.Options{Size: 26})

	a := &AXWGameboyEbitenMenuView{
		isRomfileGiven: options.RomPath != "",
		files:          files,
		fontFace:       fontface,
	}
	return a
}

func (a *AXWGameboyEbitenMenuView) Update() error {
	if a.CheckForButtonPress() {
		return ErrGoToGame
	}
	return nil
}

func (a *AXWGameboyEbitenMenuView) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{40, 40, 40, 40})
	for i, path := range a.files {
		a.DrawButton(screen, path, i)
		i++
	}
}

func (a *AXWGameboyEbitenMenuView) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 576, 1024
}

func (a *AXWGameboyEbitenMenuView) GetExitResult() []byte {
	return []byte(a.SelectedPath)
}

func (a *AXWGameboyEbitenMenuView) DrawButton(screen *ebiten.Image, label string, index int) {
	ebitenutil.DrawRect(screen, 10, float64(index*50+10), 556, 40, color.RGBA{80, 80, 80, 80})
	text.Draw(screen, label, a.fontFace, 20, index*50+37, color.White)
}

func (a *AXWGameboyEbitenMenuView) CheckForButtonPress() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_, y := ebiten.CursorPosition()
		id := (y - 10) / 50
		a.SelectedPath = a.files[id]
		return true
	}

	tids := ebiten.TouchIDs()
	for _, tid := range tids {
		_, y := ebiten.TouchPosition(tid)
		id := (y - 10) / 50
		a.SelectedPath = a.files[id]
		return true
	}
	return false
}
