package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"path"
	"strings"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/gterm"
)

type IntroScreen struct {
	PopMenu

	activeChoice int
	timeFloor    uint32

	splash []string
}

func (intro *IntroScreen) adjustSelectionWrap(delta int) {
	intro.timeFloor = sdl.GetTicks()
	intro.activeChoice += delta
	if intro.activeChoice < 0 {
		intro.activeChoice = max(0, len(menuChoices)-1)
	} else if intro.activeChoice >= len(menuChoices) {
		intro.activeChoice = 0
	}
}

func (intro *IntroScreen) Update(action controls.Action) {
	switch action {
	case controls.Up:
		intro.adjustSelectionWrap(-1)
	case controls.Down:
		intro.adjustSelectionWrap(1)
	case controls.Confirm:
		fallthrough
	case controls.Cancel:
		intro.done = true
	}
}

func (intro *IntroScreen) maybeLoadSplash() {
	if len(intro.splash) > 0 {
		return
	}

	file, err := os.Open(path.Join("assets", "art", "splash.txt"))
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

var menuChoices = []string{
	"Start Game",
	"Quit Game",
}

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
	for i, content := range menuChoices {
		var bg sdl.Color
		if i == intro.activeChoice {
			bg = activeBg
		} else {
			bg = inactiveBg
		}

		x := (window.Columns - len(content)) / 2

		window.PutStringBg(x, y, content, White, bg)
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

	content = "A 2018 7DRL by Keipra"
	x, y = (window.Columns-len(content))/2, y+1
	window.PutString(x, y, content, LightGrey)
}
