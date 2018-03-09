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

	pizza []string
}

func (intro *IntroScreen) Update(input controls.InputEvent) {
	switch input.Event.(type) {
	case *sdl.KeyDownEvent:
		intro.done = true
	}
}

func (intro *IntroScreen) maybeLoadPizza() {
	if len(intro.pizza) > 0 {
		return
	}

	file, err := os.Open(path.Join("assets", "art", "pizza.txt"))
	if err != nil {
		log.Panicln("Could not load opening art", err)
	}

	buf := bufio.NewReader(file)

	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			intro.pizza = append(intro.pizza, line)
			break
		} else if err != nil {
			log.Panicln("Failed during reading of pizza", err)
		}
		intro.pizza = append(intro.pizza, strings.TrimRight(line, "\r\n "))
	}
}

func (intro *IntroScreen) drawPizza(window *gterm.Window) {
	x, y := 27, 2
	for _, line := range intro.pizza {
		if err := window.PutString(x, y, line, White); err != nil {
			log.Panicln("Couldn't draw pizza", err)
		}
		y++
	}
}

func (intro *IntroScreen) Render(window *gterm.Window) {
	window.ClearWindow()

	intro.maybeLoadPizza()

	intro.drawPizza(window)

	content := "DeliveryRL"
	x, y := (window.Columns-len(content))/2, window.Rows/2+7
	window.PutString(x, y, "DeliveryRL", LightBlue)

	content = "Press any key to begin..."
	x, y = (window.Columns-len(content))/2, y+1
	window.PutString(x, y, "Press any key to begin...", LightGrey)
}
