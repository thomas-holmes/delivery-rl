package main

import (
	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type EndGameMenu struct {
	world *World

	Content          []string
	ContentColor     sdl.Color
	ContentRelativeX int
	ContentRelativeY int

	PopMenu
}

func (pop *EndGameMenu) Update(action controls.Action) {
	pop.CheckCancel(action)
	switch action {
	case controls.Quit:
		fallthrough
	case controls.Confirm:
		pop.done = true
	}

	if pop.done {
		pop.world.QuitGame = true
	}
}

func (pop *EndGameMenu) RenderBorder(window *gterm.Window) {
	pop.DrawBox(window, pop.ContentColor)
}

func (pop *EndGameMenu) RenderContents(window *gterm.Window) {
	xOffset := pop.X + pop.ContentRelativeX
	yOffset := pop.Y + pop.ContentRelativeY

	for line, content := range pop.Content {
		window.PutString(xOffset, yOffset+line, content, pop.ContentColor)
	}
}

func (pop *EndGameMenu) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	pop.RenderBorder(window)
	pop.RenderContents(window)
}

func NewEndGameMenu(world *World, x int, y int, w int, h int, color sdl.Color, contents ...string) EndGameMenu {
	contentLen := len(contents)
	maxWidth := 0

	for _, content := range contents {
		thisLen := len(content)
		if thisLen > maxWidth {
			maxWidth = thisLen
		}
	}

	centeredXOffset := (w - maxWidth) / 2
	centeredYOffset := (h - contentLen) / 2

	pop := EndGameMenu{
		world: world,
		PopMenu: PopMenu{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
		ContentRelativeX: centeredXOffset,
		ContentRelativeY: centeredYOffset,
		Content:          contents,
		ContentColor:     color,
	}

	return pop
}
