package main

import (
	"math"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type FullGameLog struct {
	GameLog *GameLog

	tranFinished bool
	tranStart    uint32

	closing bool

	topPosition int

	PopMenu
}

const TransitionDuration = 250

func (pop *FullGameLog) Update(action controls.Action) {
	pop.GameLog.FullLogShown = true
	if action == controls.Cancel && !pop.closing {
		pop.closing = true

		start := pop.tranStart
		now := sdl.GetTicks()

		// Using int64 so I can subtract uint32s
		remainingAnimationTime := max64(0, int64(TransitionDuration)-(int64(now)-int64(start)))
		pop.tranStart = now - uint32(remainingAnimationTime)

		pop.tranFinished = false
	}

	if pop.tranFinished {
		pop.done = true
		pop.GameLog.FullLogShown = false
	}
}

func (pop *FullGameLog) RenderVisibleLines(window *gterm.Window) {
	messages := gl.Messages()

	yOffset := 1
	for i := 0; i < len(messages); i++ {
		idx := len(messages) - 1 - i
		if idx < 0 {
			break
		}

		message := messages[idx]
		lines := wrapText(message, 0, 2, pop.W)

		yOffset += len(lines)

		// Make sure we have enough headroom
		if yOffset >= pop.H-pop.topPosition {
			break
		}

		for j := 0; j < len(lines); j++ {
			window.PutString(pop.X, pop.Y+pop.H-yOffset+j, lines[j], White)
		}
	}

}

func (pop *FullGameLog) DrawTransition(window *gterm.Window) {
	now := sdl.GetTicks()
	if pop.tranStart == 0 {
		pop.tranStart = now
	}

	adjustment := pop.GameLog.H + 1
	_ = adjustment

	distance := pop.H - adjustment

	elapsed := now - pop.tranStart

	pctOpen := float64(elapsed) / TransitionDuration
	if pop.closing {
		pctOpen = 1 - pctOpen
	}

	pctOpen = math.Min(1, math.Max(0, pctOpen))

	actualOpen := int(float64(distance) * pctOpen)

	distance += adjustment
	actualOpen += adjustment

	if pop.closing && actualOpen == adjustment {
		pop.tranFinished = true
	}

	y := pop.Y + pop.H - actualOpen
	pop.topPosition = y
	window.ClearRegion(pop.X, y, pop.W, actualOpen)

	for x := pop.X; x < pop.X+pop.W; x++ {
		window.PutRune(x, y, horizontal, White, gterm.NoColor)
	}
	window.PutRune(pop.X+pop.W, y, vertLeftJoint, White, gterm.NoColor)

}

func (pop *FullGameLog) Render(window *gterm.Window) {
	pop.DrawTransition(window)

	pop.RenderVisibleLines(window)
}
