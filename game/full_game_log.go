package main

import (
	"log"
	"math"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type FullGameLog struct {
	GameLog *GameLog

	PopMenu

	ScrollPosition int
}

func (pop *FullGameLog) ScrollDown(distance int) {
	maxScrollPosition := max(0, len(pop.GameLog.Messages)-pop.H)
	pop.ScrollPosition = min(maxScrollPosition, pop.ScrollPosition+distance)
}

func (pop *FullGameLog) ScrollUp(distance int) {
	pop.ScrollPosition = max(0, pop.ScrollPosition-distance)
}

func (pop *FullGameLog) Update(input InputEvent) {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_ESCAPE:
			pop.done = true
		case sdl.K_k:
			fallthrough
		case sdl.K_UP:
			pop.ScrollUp(1)
		case sdl.K_PAGEUP:
			pop.ScrollUp(10)
		case sdl.K_HOME:
			pop.ScrollUp(len(pop.GameLog.Messages))
		case sdl.K_j:
			fallthrough
		case sdl.K_DOWN:
			pop.ScrollDown(1)
		case sdl.K_PAGEDOWN:
			pop.ScrollDown(10)
		case sdl.K_END:
			pop.ScrollDown(len(pop.GameLog.Messages))
		}
	}
}

func (pop *FullGameLog) RenderScrollBar(window *gterm.Window) {
	top := '^'
	bottom := 'v'

	barSpace := float64(pop.H - 2)

	percentageShown := float64(min(pop.H, len(pop.GameLog.Messages))) / float64(len(pop.GameLog.Messages))
	scrollBarWidth := int(math.Ceil(barSpace * percentageShown))

	topOfBar := int(float64(barSpace) * float64(pop.ScrollPosition) / float64(len(pop.GameLog.Messages)))

	window.PutRune(pop.X, pop.Y, top, Yellow, gterm.NoColor)

	barRunesDrawn := 0
	for row := pop.Y + 1; row < (pop.Y + pop.H - 1); row++ {
		if (row >= topOfBar) && (barRunesDrawn <= scrollBarWidth) {
			window.PutRune(pop.X, row, '#', Yellow, gterm.NoColor)
			barRunesDrawn++
		} else {
			window.PutRune(pop.X, row, '|', Grey, gterm.NoColor)
		}
	}

	window.PutRune(pop.X, pop.Y+pop.H-1, bottom, Yellow, gterm.NoColor)
}

func (pop *FullGameLog) RenderVisibleLines(window *gterm.Window) {

	messagesToRender := len(pop.GameLog.Messages) - pop.ScrollPosition

	yOffset := 0
	for i := messagesToRender - 1; i >= 0; i-- {
		message := pop.GameLog.Messages[i]
		window.PutString(pop.X+1, pop.Y+yOffset, message, White)
		yOffset++
	}

}

func (pop *FullGameLog) Render(window *gterm.Window) {
	if err := window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Println("Got an error clearing FullGameLog region", err)
	}
	pop.RenderScrollBar(window)

	pop.RenderVisibleLines(window)
}
