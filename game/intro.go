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
)

type IntroScreen struct {
	PopMenu

	splash []string
}

func (intro *IntroScreen) Update(action controls.Action) {
	switch action {
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

func (intro *IntroScreen) Render(window *gterm.Window) {
	window.ClearWindow()

	intro.maybeLoadSplash()

	intro.drawSplash(window)

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
