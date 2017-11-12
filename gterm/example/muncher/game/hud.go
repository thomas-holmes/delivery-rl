package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type HUD struct {
	World  *World
	Player *Player
	XPos   int
	YPos   int
	Width  int
	Height int
	Dirty  bool
}

func NewHud(player *Player, world *World, xPos int, yPos int) *HUD {
	hud := HUD{
		Player: player,
		World:  world,
		XPos:   xPos,
		YPos:   yPos,
		Dirty:  true,
	}

	world.MessageBus.Subscribe(&hud)

	return &hud
}

func (hud *HUD) Notify(message Message, data interface{}) {
	switch message {
	case PlayerUpdate:
		hud.Dirty = true
	}
}

var Green = sdl.Color{R: 0, G: 255, B: 0, A: 255}
var Yellow = sdl.Color{R: 255, G: 255, B: 0, A: 255}
var Orange = sdl.Color{R: 255, G: 192, B: 0, A: 255}
var Red = sdl.Color{R: 255, G: 0, B: 0, A: 255}

func (hud HUD) renderPlayerName(world *World) {
	world.Window.ClearCell(hud.XPos, hud.YPos)
	world.Window.AddToCell(hud.XPos, hud.YPos, hud.Player.Name, Yellow)
}

func (hud HUD) renderPlayerPosition(world *World) {
	world.Window.ClearCell(hud.XPos, hud.YPos+1)
	position := fmt.Sprintf("(%v, %v)", hud.Player.xPos, hud.Player.yPos)
	world.Window.AddToCell(hud.XPos, hud.YPos+1, position, Yellow)
}

func (hud HUD) renderPlayerHealth(world *World) {
	world.Window.ClearCell(hud.XPos, hud.YPos+2)

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

	world.Window.AddToCell(hud.XPos, hud.YPos+2, hp, hpColor)
}

func (hud *HUD) Render(world *World) {
	if hud.Dirty {
		hud.renderPlayerName(world)
		hud.renderPlayerPosition(world)
		hud.renderPlayerHealth(world)

		hud.Dirty = false
	}
}
