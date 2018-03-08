package main

import (
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type HelpPop struct {
	PopMenu
}

func NewHelpPop(x, y, w, h int) *HelpPop {
	return &HelpPop{
		PopMenu: PopMenu{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
	}
}

func (pop *HelpPop) Update(input InputEvent) {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		k := e.Keysym.Sym
		switch {
		case k == sdl.K_ESCAPE:
			pop.done = true
		}
	}
}

func (pop HelpPop) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	pop.DrawBox(window, White)

	x := pop.X + 2
	y := pop.Y + 2

	controls := []string{
		"Left:          h",
		"Right:         l",
		"Down:          j",
		"Up:            k",
		"Inventory:     i",
		"Equip Weapon:  e",
		"Use Ability:   z",
		"Quit:          CTRL-q",
		"Game Log:      m",
		"Toggle FPS:    F12",
		"Help:          ?",
	}

	window.PutString(x, y, "DeliveryRL Controls", White)

	y += 2

	for _, s := range controls {
		window.PutString(x, y, s, White)
		y++
	}
}
