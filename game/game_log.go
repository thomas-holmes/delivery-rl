package main

import (
	"log"

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

func (gameLog *GameLog) Render(window *gterm.Window) {
	messages := gamelog.Messages()
	lastMessage := max(0, len(messages))
	firstMessage := max(0, lastMessage-gameLog.H)

	for i := firstMessage; i < lastMessage; i++ {
		message := messages[i]
		yOffset := lastMessage - i - 1
		cut := min(len(message), gameLog.W)
		err := window.PutString(gameLog.X, gameLog.Y+gameLog.H-yOffset, messages[i][:cut], White)
		if err != nil {
			log.Println("Failed ot render log", err)
		}
	}
}

func (gameLog *GameLog) Notify(message m.M) {
	switch message.ID {
	case ShowFullGameLog:
		menu := &FullGameLog{PopMenu: PopMenu{X: 0, Y: 0, W: gameLog.world.CameraWidth, H: gameLog.world.Window.Rows}, GameLog: gameLog}
		m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{menu}})
	}
}
