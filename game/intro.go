package main

import "github.com/thomas-holmes/gterm"
import "github.com/veandco/go-sdl2/sdl"

type IntroScreen struct {
	PopMenu
}

func (intro *IntroScreen) Update(input InputEvent) bool {
	switch input.Event.(type) {
	case *sdl.KeyDownEvent:
		intro.done = true
	}
	return true
}

func (intro *IntroScreen) Render(window *gterm.Window) {
	window.ClearWindow()

	content := "DeliveryRL"
	x, y := (window.Columns-len(content))/2, window.Rows/2
	window.PutString(x, y, "DeliveryRL", LightBlue)

	content = "Press any key to begin..."
	x, y = (window.Columns-len(content))/2, window.Rows/2+1
	window.PutString(x, y, "Press any key to begin...", Grey)
}
