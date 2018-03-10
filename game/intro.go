package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type IntroScreen struct {
	PopMenu

	splash []string
}

func (intro *IntroScreen) Update(input controls.InputEvent) {
	switch input.Event.(type) {
	case sdl.KeyDownEvent:
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

func (intro *IntroScreen) Render(window *gterm.Window) {
	window.ClearWindow()

	intro.maybeLoadSplash()

	intro.drawSplash(window)

	content := "DeliveryRL"
	x, y := (window.Columns-len(content))/2-3, window.Rows-5
	window.PutString(x, y, "DeliveryRL", LightBlue)
	log.Printf("Drawing at %d", x)

	content = "Press any key to begin..."
	x, y = (window.Columns-len(content))/2-3, y+1
	window.PutString(x, y, "Press any key to begin...", LightGrey)
}
