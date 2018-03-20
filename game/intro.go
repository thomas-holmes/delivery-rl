package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/gterm"
)

var StartMenuChoice simpleMenu = simpleMenu("Start Game")
var QuitMenuChoice simpleMenu = simpleMenu("End Game")

func NewIntroScreen(window *gterm.Window) *IntroScreen {
	i := &IntroScreen{window: window}

	i.fonts = getAvailableFonts()

	i.setInitialFont()

	i.menuChoices = []menuChoice{
		StartMenuChoice,
		i.fonts,
		QuitMenuChoice,
	}

	return i
}

type IntroScreen struct {
	PopMenu

	window *gterm.Window

	fonts *availableFonts

	activeChoice    int
	timeFloor       uint32
	checkedForFonts bool

	menuChoices []menuChoice

	splash []string
}

func (intro *IntroScreen) setInitialFont() {
	fonts := intro.fonts

	configuredFontPath, err := filepath.Abs(path.Join(AssetRoot, Font))
	if err != nil {
		log.Panicln("Your file system is broken", err)
	}

	for i, f := range fonts.fonts {
		thisPath, err := filepath.Abs(f.Path)
		if err != nil {
			log.Panicln("Your file system is broken", err)
		}
		if thisPath == configuredFontPath {
			fonts.selectedFont = i
			return
		}
	}
}

func (intro *IntroScreen) adjustSelectionWrap(delta int) {
	intro.timeFloor = sdl.GetTicks() - 4000
	intro.activeChoice += delta
	if intro.activeChoice < 0 {
		intro.activeChoice = max(0, len(intro.menuChoices)-1)
	} else if intro.activeChoice >= len(intro.menuChoices) {
		intro.activeChoice = 0
	}
}

func (intro *IntroScreen) Update(action controls.Action) {
	switch action {
	case controls.Up:
		intro.adjustSelectionWrap(-1)
	case controls.Down:
		intro.adjustSelectionWrap(1)
	case controls.Left:
		if a, ok := intro.menuChoices[intro.activeChoice].(*availableFonts); ok {
			a.adjust(-1)
			font := a.fonts[a.selectedFont]
			log.Printf("selected new font at: %s", font.Path)
			intro.window.ChangeFont(font.Path, font.W, font.H)
		}
	case controls.Right:
		if a, ok := intro.menuChoices[intro.activeChoice].(*availableFonts); ok {
			a.adjust(1)
			font := a.fonts[a.selectedFont]
			log.Printf("selected new font at: %s", font.Path)
			intro.window.ChangeFont(font.Path, font.W, font.H)
		}
	case controls.Confirm:
		switch intro.menuChoices[intro.activeChoice] {
		case StartMenuChoice:
			fallthrough
		case QuitMenuChoice:
			intro.done = true
		}
	}
}

type font struct {
	Name string
	Path string
	W    int
	H    int
}

func (a availableFonts) Label() string {
	return fmt.Sprintf("Font: %s", a.fonts[a.selectedFont].Name)
}

func (a *availableFonts) adjust(delta int) {
	a.selectedFont += delta
	if a.selectedFont < 0 {
		a.selectedFont = max(0, len(a.fonts)-1)
	} else if a.selectedFont >= len(a.fonts) {
		a.selectedFont = 0
	}
}

func getAvailableFonts() *availableFonts {
	fontPath := path.Join(AssetRoot, "assets", "font")
	files, err := ioutil.ReadDir(fontPath)
	if err != nil {
		log.Printf("Could not find any fonts at %s, %s", fontPath, err)
	}

	var fontNames []font
	for _, f := range files {
		var name string
		var w, h int
		_, err := fmt.Sscanf(f.Name(), "%5s_%2dx%2d.png", &name, &w, &h)
		if err != nil {
			log.Printf("Found incorrectly named file %s", f.Name())
			continue
		}

		currentAbs, err := filepath.Abs(path.Join(fontPath, f.Name()))
		if err != nil {
			log.Panicln("Things have gone horribly wrong with your filesystem", err)
		}

		foundFont := font{
			Name: fmt.Sprintf("%s_%dx%d", name, w, h),
			Path: currentAbs,
			W:    w,
			H:    h,
		}

		fontNames = append(fontNames, foundFont)
		sort.Slice(fontNames, func(i, j int) bool { return fontNames[i].W < fontNames[j].W })
	}

	dfp, err := filepath.Abs(DefaultFontPath)
	if err != nil {
		log.Panicln("Things have gone horribly wrong with your filesystem", err)
	}
	yfp, err := filepath.Abs(Font)
	if err != nil {
		log.Panicln("Things have gone horribly wrong with your filesystem", err)
	}

	if yfp != dfp {

		var yfpIsProvidedFont bool

		for _, font := range fontNames {
			if font.Path == yfp {
				yfpIsProvidedFont = true
			}
		}

		if !yfpIsProvidedFont {
			fileName := filepath.Base(yfp)
			customFont := font{
				Name: fileName,
				Path: yfp,
				W:    FontH,
				H:    FontW,
			}
			fontNames = append([]font{customFont}, fontNames...)
		}
	}

	return &availableFonts{
		fonts: fontNames,
	}
}

func (intro *IntroScreen) maybeLoadSplash() {
	if len(intro.splash) > 0 {
		return
	}

	file, err := os.Open(path.Join(AssetRoot, "assets", "art", "splash.txt"))
	if err != nil {
		log.Panicln("Could not load opening art", err)
	}

	buf := bufio.NewReader(file)

	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			intro.splash = append(intro.splash, line)
			break
		} else if err != nil {
			log.Panicln("Failed during reading of splash", err)
		}
		intro.splash = append(intro.splash, strings.TrimRight(line, "\r\n "))
	}
}

func (intro *IntroScreen) drawSplash(window *gterm.Window) {
	x, y := 0, 0
	for _, line := range intro.splash {
		if err := window.PutString(x, y, line, White); err != nil {
			log.Panicln("Couldn't draw splash", err)
		}
		y++
	}
}

func (intro *IntroScreen) StartGame() bool {
	return intro.Done() && intro.menuChoices[intro.activeChoice] == StartMenuChoice
}

func (intro *IntroScreen) QuitGame() bool {
	return intro.Done() && intro.menuChoices[intro.activeChoice] == QuitMenuChoice
}

type availableFonts struct {
	fonts        []font
	selectedFont int
}

type menuChoice interface {
	Label() string
}

type simpleMenu string

func (s simpleMenu) Label() string { return string(s) }

func (intro *IntroScreen) drawMenuItems(window *gterm.Window) {
	ticks := sdl.GetTicks()
	if intro.timeFloor == 0 {
		intro.timeFloor = ticks
	}
	time := (ticks - intro.timeFloor) / 10.0
	angle := float64(time%360) * (math.Pi / 180.0)
	cos := math.Cos(angle)

	r, g, b := 127, 127, 127

	r = int(float64(r) + (float64(127 * cos)))
	g = int(float64(g) + (float64(127 * cos)))
	b = int(float64(b) + (float64(127 * cos)))

	activeBg := sdl.Color{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	inactiveBg := gterm.NoColor

	y := window.Rows - 10
	for i, choice := range intro.menuChoices {
		var bg sdl.Color
		if i == intro.activeChoice {
			bg = activeBg
		} else {
			bg = inactiveBg
		}

		x := (window.Columns - len(choice.Label())) / 2

		window.PutStringBg(x, y, choice.Label(), White, bg)
		y++
	}

}

func (intro *IntroScreen) Render(window *gterm.Window) {
	window.ClearWindow()

	intro.maybeLoadSplash()

	intro.drawSplash(window)

	intro.drawMenuItems(window)

	content := "DeliveryRL"
	x, y := (window.Columns-len(content))/2, window.Rows-5
	window.PutString(x, y, content, LightBlue)

	content = "Press any key to begin..."
	x, y = (window.Columns-len(content))/2, y+1
	window.PutString(x, y, content, LightGrey)

	content = "A 2018 7DRL by"
	name := "keipra"
	x, y = (window.Columns-(len(content)+len(name)))/2, y+1
	window.PutString(x, y, content, LightGrey)

	window.PutString(x+len(content)+1, y, name, KeipraPurple)
}

func (intro *IntroScreen) Reset() {
	intro.done = false
	intro.activeChoice = 0
	intro.timeFloor = 0
}
