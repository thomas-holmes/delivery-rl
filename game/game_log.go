package main

import (
	"github.com/thomas-holmes/delivery-rl/game/gamelog"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
)

// This is going to result in a lot of heap allocations
type GameLog struct {
	world *World

	X int
	Y int
	W int
	H int
}

func NewGameLog(x int, y int, w int, h int, world *World) *GameLog {
	gameLog := GameLog{
		world: world,
		X:     x,
		Y:     y,
		W:     w,
		H:     h,
	}

	m.Subscribe(gameLog.Notify)

	return &gameLog
}

func (pop *GameLog) Render(window *gterm.Window) {
	y := pop.Y
	for x := pop.X; x < pop.X+pop.W; x++ {
		window.PutRune(x, y, horizontal, White, gterm.NoColor)
	}
	messages := gamelog.Messages()
	var yOffset int
	for i := 0; i < len(messages); i++ {
		idx := len(messages) - 1 - i
		if idx < 0 {
			break
		}

		message := messages[idx]
		lines := wrapText(message, 0, 2, pop.W)

		yOffset += len(lines)

		// Make sure we have enough headroom
		if yOffset >= pop.H {
			break
		}

		for j := 0; j < len(lines); j++ {
			window.PutString(pop.X, pop.Y+pop.H-yOffset+j, lines[j], White)
		}
	}
}

func (gameLog *GameLog) Notify(message m.M) {
	switch message.ID {
	case ShowFullGameLog:
		menu := &FullGameLog{PopMenu: PopMenu{X: 0, Y: 0, W: 65, H: gameLog.world.Window.Rows}, GameLog: gameLog}
		m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{menu}})
	}
}
