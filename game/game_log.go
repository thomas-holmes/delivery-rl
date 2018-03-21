package main

import (
	"math"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

// This is going to result in a lot of heap allocations
type GameLog struct {
	world *World

	FullLogShown bool

	tranFinished bool
	tranStart    uint32

	closing bool

	topPosition int

	PopMenu
}

func NewGameLog(x int, y int, w int, h int, world *World) *GameLog {
	gameLog := GameLog{
		world: world,
		PopMenu: PopMenu{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
		closing:   true,
		tranStart: sdl.GetTicks() - TransitionDuration,
	}

	m.Subscribe(gameLog.Notify)

	return &gameLog
}

func (pop *GameLog) Update(action controls.Action) {
	pop.FullLogShown = true

	if action == controls.Cancel && !pop.closing {
		pop.closePartWayOpen()
	}

	if pop.tranFinished {
		pop.closeFullGameLog()
	}
}

const TransitionDuration = 200

func (pop *GameLog) DrawTransition(window *gterm.Window) {
	now := sdl.GetTicks()
	if pop.tranStart == 0 {
		pop.tranStart = now
	}

	adjustment := pop.H + 1
	_ = adjustment

	distance := window.Rows - adjustment

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

	y := window.Rows - actualOpen
	pop.topPosition = y
	window.ClearRegion(pop.X, y, pop.W, actualOpen)

	for x := pop.X; x < pop.X+pop.W; x++ {
		window.PutRune(x, y, horizontal, White, gterm.NoColor)
	}
	window.PutRune(pop.X+pop.W, y, vertLeftJoint, White, gterm.NoColor)

}

func (pop *GameLog) RenderVisibleLines(window *gterm.Window) {
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
		if yOffset >= window.Rows-pop.topPosition {
			break
		}

		for j := 0; j < len(lines); j++ {
			window.PutString(pop.X, window.Rows-yOffset+j, lines[j], White)
		}
	}

}

func (pop *GameLog) Render(window *gterm.Window) {
	pop.DrawTransition(window)

	pop.RenderVisibleLines(window)
}

func (pop *GameLog) closePartWayOpen() {
	pop.closing = true

	start := pop.tranStart
	now := sdl.GetTicks()

	// Using int64 so I can subtract uint32s
	remainingAnimationTime := max64(0, int64(TransitionDuration)-(int64(now)-int64(start)))
	pop.tranStart = now - uint32(remainingAnimationTime)

	pop.tranFinished = false
}

func (pop *GameLog) closeFullGameLog() {
	pop.done = true
	pop.FullLogShown = false
	pop.closing = true
	pop.tranFinished = false
}

func (pop *GameLog) openFullLog() {
	pop.done = false
	pop.closing = false
	pop.tranFinished = false
	pop.tranStart = sdl.GetTicks()
	m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: pop}})
}

func (pop *GameLog) Notify(message m.M) {
	switch message.ID {
	case ShowFullGameLog:
		pop.openFullLog()
	}
}
