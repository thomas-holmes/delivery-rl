package main

import (
	"fmt"
	"time"
)

type HUD struct {
	World  *World
	Player *Player
	XPos   int
	YPos   int
	Width  int
	Height int
}

func NewHud(player *Player, world *World, xPos int, yPos int) *HUD {
	hud := HUD{
		Player: player,
		World:  world,
		XPos:   xPos,
		YPos:   yPos,
	}

	world.MessageBus.Subscribe(&hud)

	return &hud
}

func (hud *HUD) Notify(message Message, data interface{}) {
	switch message {
	case PlayerUpdate:
	}
}
func (hud HUD) renderPlayerName(world *World) {
	world.Window.PutString(hud.XPos, hud.YPos, world.Player.Name, Yellow)
}

func (hud HUD) renderPlayerPosition(world *World) {
	position := fmt.Sprintf("(%v, %v)", hud.Player.X, hud.Player.Y)
	world.Window.PutString(hud.XPos, hud.YPos+1, position, Yellow)
}

func (hud HUD) renderPlayerHealth(world *World) {
	hpColor := Red

	pct := hud.Player.HealthPercentage()
	switch {
	case pct >= 0.8:
		hpColor = Green
	case pct >= 0.6:
		hpColor = Yellow
	case pct >= 0.4:
		hpColor = Orange
	default:
		hpColor = Red
	}

	hp := fmt.Sprintf("%v/%v", hud.Player.HP.Current, hud.Player.HP.Max)
	if hud.Player.HP.Current == 0 {
		hp += " *DEAD*"
	}

	world.Window.PutString(hud.XPos, hud.YPos+2, hp, hpColor)
}

func (hud HUD) renderPlayerLevel(world *World) {
	level := fmt.Sprintf("Level: %v", hud.Player.Level)
	world.Window.PutString(hud.XPos, hud.YPos+3, level, Yellow)
}

func (hud *HUD) Render(world *World) {
	defer timeMe(time.Now(), "HUD.Render")
	hud.renderPlayerName(world)
	hud.renderPlayerPosition(world)
	hud.renderPlayerHealth(world)
	hud.renderPlayerLevel(world)
}