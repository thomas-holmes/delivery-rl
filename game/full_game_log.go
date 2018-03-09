package main

import (
	"log"
	"math"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type FullGameLog struct {
	GameLog *GameLog

	PopMenu

	ScrollPosition int
}

func (pop *FullGameLog) ScrollDown(distance int) {
	messages := gl.Messages()
	maxScrollPosition := max(0, len(messages)-pop.H)
	pop.ScrollPosition = min(maxScrollPosition, pop.ScrollPosition+distance)
}

func (pop *FullGameLog) ScrollUp(distance int) {
	pop.ScrollPosition = max(0, pop.ScrollPosition-distance)
}

func (pop *FullGameLog) Update(input controls.InputEvent) {
	pop.CheckCancel(input)

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
			pop.ScrollUp(len(gl.Messages()))
		case sdl.K_j:
			fallthrough
		case sdl.K_DOWN:
			pop.ScrollDown(1)
		case sdl.K_PAGEDOWN:
			pop.ScrollDown(10)
		case sdl.K_END:
			pop.ScrollDown(len(gl.Messages()))
		}
	}
}

func (pop *FullGameLog) RenderScrollBar(window *gterm.Window) {
	messages := gl.Messages()
	barSpace := float64(pop.H - 2)

	percentageShown := float64(min(pop.H, len(messages))) / float64(len(messages))
	scrollBarWidth := int(math.Ceil(barSpace * percentageShown))

	topOfBar := int(float64(barSpace) * float64(pop.ScrollPosition) / float64(len(messages)))

	window.PutRune(pop.X, pop.Y, upArrow, Yellow, gterm.NoColor)

	barRunesDrawn := 0
	for row := pop.Y + 1; row < (pop.Y + pop.H - 1); row++ {
		if (row >= topOfBar) && (barRunesDrawn <= scrollBarWidth) {
			window.PutRune(pop.X, row, fullBlock, Yellow, gterm.NoColor)
			barRunesDrawn++
		} else {
			window.PutRune(pop.X, row, vertical, Grey, gterm.NoColor)
		}
	}

	window.PutRune(pop.X, pop.Y+pop.H-1, downArrow, Yellow, gterm.NoColor)
}

func (pop *FullGameLog) RenderVisibleLines(window *gterm.Window) {
	messages := gl.Messages()
	messagesToRender := len(messages) - pop.ScrollPosition

	yPos := pop.Y
	for i := messagesToRender - 1; i >= 0 && yPos+1 < pop.Y+pop.H; i-- {
		message := messages[i]
		yPos += putWrappedText(window, message, pop.X+1, yPos, 0, 2, pop.W-2, White)
		/*
			window.PutString(pop.X+1, pop.Y+yOffset, message, White)
			yOffset++
		*/
	}

}

func (pop *FullGameLog) Render(window *gterm.Window) {
	if err := window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Println("Got an error clearing FullGameLog region", err)
	}
	pop.RenderScrollBar(window)

	pop.RenderVisibleLines(window)
}
