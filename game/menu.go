package main

import (
	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type PopMenu struct {
	done bool

	X int
	Y int
	W int
	H int
}

func (pop PopMenu) Done() bool {
	return pop.done
}

func (pop PopMenu) DrawBox(window *gterm.Window, color sdl.Color) {
	window.PutRune(pop.X, pop.Y, topLeft, color, gterm.NoColor)
	for x := pop.X + 1; x < pop.X+pop.W-1; x++ {
		window.PutRune(x, pop.Y, horizontal, color, gterm.NoColor)
	}
	window.PutRune(pop.X+pop.W-1, pop.Y, topRight, color, gterm.NoColor)

	for y := pop.Y + 1; y < pop.Y+pop.H-1; y++ {
		window.PutRune(pop.X, y, vertical, color, gterm.NoColor)

		window.PutRune(pop.X+pop.W-1, y, vertical, color, gterm.NoColor)
	}

	window.PutRune(pop.X, pop.Y+pop.H-1, bottomLeft, color, gterm.NoColor)
	for x := pop.X + 1; x < pop.X+pop.W-1; x++ {
		window.PutRune(x, pop.Y+pop.H-1, horizontal, color, gterm.NoColor)
	}
	window.PutRune(pop.X+pop.W-1, pop.Y+pop.H-1, bottomRight, color, gterm.NoColor)
}

type Menu interface {
	Update(controls.InputEvent)
	Render(window *gterm.Window)
	Done() bool
}
