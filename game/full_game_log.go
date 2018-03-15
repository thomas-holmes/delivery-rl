package main

import (
	"log"
	"math"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	"github.com/thomas-holmes/gterm"
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

func (pop *FullGameLog) Update(action controls.Action) {
	pop.CheckCancel(action)

	switch action {
	case controls.Up:
		pop.ScrollUp(1)
	case controls.SkipUp:
		pop.ScrollUp(10)
	case controls.Top:
		pop.ScrollUp(len(gl.Messages()))
	case controls.Down:
		pop.ScrollDown(1)
	case controls.SkipDown:
		pop.ScrollDown(10)
	case controls.Bottom:
		pop.ScrollDown(len(gl.Messages()))
	}
}

func (pop *FullGameLog) RenderScrollBar(window *gterm.Window) {
	messages := gl.Messages()
	barSpace := float64(pop.H - 4)

	percentageShown := float64(min(pop.H, len(messages))) / float64(len(messages))
	scrollBarWidth := int(math.Ceil(barSpace * percentageShown))

	topOfBar := int(float64(barSpace) * float64(pop.ScrollPosition) / float64(len(messages)))

	window.PutRune(pop.X+1, pop.Y+1, upArrow, Yellow, gterm.NoColor)

	barRunesDrawn := 0
	for row := pop.Y + 2; row < (pop.Y + 1 + pop.H - 3); row++ {
		if (row >= topOfBar) && (barRunesDrawn <= scrollBarWidth) {
			window.PutRune(pop.X+1, row, fullBlock, Yellow, gterm.NoColor)
			barRunesDrawn++
		} else {
			window.PutRune(pop.X+1, row, vertical, Grey, gterm.NoColor)
		}
	}

	window.PutRune(pop.X+1, pop.Y+1+pop.H-3, downArrow, Yellow, gterm.NoColor)
}

func (pop *FullGameLog) RenderVisibleLines(window *gterm.Window) {
	messages := gl.Messages()
	messagesToRender := len(messages) - pop.ScrollPosition

	yPos := pop.Y
	for i := messagesToRender - 1; i >= 0 && yPos+1 < pop.Y+pop.H; i-- {
		message := messages[i]
		yPos += putWrappedText(window, message, pop.X+3, yPos+2, 0, 2, pop.W-4, White)
	}

}

func (pop *FullGameLog) Render(window *gterm.Window) {
	if err := window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Println("Got an error clearing FullGameLog region", err)
	}
	pop.DrawBox(window, White)
	pop.RenderScrollBar(window)

	pop.RenderVisibleLines(window)
}
