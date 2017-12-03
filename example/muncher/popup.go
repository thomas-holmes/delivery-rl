package main

import (
	"time"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type PopUp struct {
	XPos             int
	YPos             int
	Width            int
	Height           int
	Content          []string
	ContentColor     sdl.Color
	ContentRelativeX int
	ContentRelativeY int

	Shown bool

	Messaging
}

func NewPopUp(xPos int, yPos int, width int, height int, color sdl.Color, contents ...string) PopUp {
	contentLen := len(contents)
	maxWidth := 0

	for _, content := range contents {
		thisLen := len(content)
		if thisLen > maxWidth {
			maxWidth = thisLen
		}
	}

	centeredXOffset := (width - maxWidth) / 2
	centeredYOffset := (height - contentLen) / 2

	pop := PopUp{
		XPos:             xPos,
		YPos:             yPos,
		Width:            width,
		Height:           height,
		ContentRelativeX: centeredXOffset,
		ContentRelativeY: centeredYOffset,
		Content:          contents,
		ContentColor:     color,
	}

	return pop
}

func (pop PopUp) ClearUnderlying() {
	for y := pop.YPos; y < pop.YPos+pop.Height; y++ {
		for x := pop.XPos; x < pop.XPos+pop.Width; x++ {
			pop.Broadcast(ClearRegion, ClearRegionMessage{
				XPos:   x,
				YPos:   y,
				Width:  pop.Width,
				Height: pop.Height,
			})
		}
	}
}

func (pop *PopUp) Show() {
	pop.Shown = true
	pop.Broadcast(PopUpShown, nil)
	// Broadcast a message stopping other game systems
}

func (pop *PopUp) Hide() {
	pop.Shown = false
	pop.Broadcast(PopUpHidden, nil)
	// Broadcast message resuming other game systems
}

func (pop *PopUp) RenderBorder(window *gterm.Window) {
	leftX := pop.XPos
	rightX := pop.XPos + pop.Width - 1

	topY := pop.YPos
	bottomY := pop.YPos + pop.Height - 1

	for y := topY; y <= bottomY; y++ {
		for x := leftX; x <= rightX; x++ {
			if y == topY || y == bottomY || x == leftX || x == rightX {
				window.PutRune(x, y, '%', pop.ContentColor, gterm.NoColor)
			}
		}
	}
}

func (pop PopUp) RenderContents(window *gterm.Window) {
	xOffset := pop.XPos + pop.ContentRelativeX
	yOffset := pop.YPos + pop.ContentRelativeY

	for line, content := range pop.Content {
		window.PutString(xOffset, yOffset+line, content, pop.ContentColor)
	}
}

func (pop *PopUp) Render(window *gterm.Window) {
	defer timeMe(time.Now(), "PopUp.Render")
	pop.ClearUnderlying()
	pop.RenderBorder(window)
	pop.RenderContents(window)
}