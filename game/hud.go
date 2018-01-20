package main

import (
	"fmt"
	"log"

	"github.com/thomas-holmes/gterm"
)

type HUD struct {
	World  *World
	Player *Creature
	XPos   int
	YPos   int
	Width  int
	Height int

	nextFreeRow int
}

func NewHud(player *Creature, world *World, xPos int, yPos int) *HUD {
	hud := HUD{
		Player:      player,
		World:       world,
		XPos:        xPos,
		YPos:        yPos,
		nextFreeRow: 0,
	}

	world.messageBus.Subscribe(&hud)

	return &hud
}

func (hud *HUD) GetNextRow() int {
	hud.nextFreeRow++
	return hud.nextFreeRow + hud.YPos - 1
}

func (hud *HUD) Notify(message Message, data interface{}) {
	switch message {
	case PlayerUpdate:
	}
}
func (hud *HUD) renderPlayerName(world *World) {
	world.Window.PutString(hud.XPos, hud.GetNextRow(), world.Player.Name, Yellow)
}

func (hud *HUD) renderPlayerPosition(world *World) {
	position := fmt.Sprintf("(%v, %v) - Level %v", hud.Player.X, hud.Player.Y, hud.Player.Depth+1)
	world.Window.PutString(hud.XPos, hud.GetNextRow(), position, Yellow)
}

func (hud *HUD) renderPlayerHealth(world *World) {
	hpColor := Red

	pct := hud.Player.HP.Percentage()
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

	label := "Health:"
	hp := fmt.Sprintf("%v/%v", hud.Player.HP.Current, hud.Player.HP.Max)
	if hud.Player.HP.Current == 0 {
		hp += " *DEAD*"
	}

	row := hud.GetNextRow()
	if err := world.Window.PutString(hud.XPos, row, label, Yellow); err != nil {
		log.Fatalln("Couldn't write HUD hp label", err)
	}

	if err := world.Window.PutString(hud.XPos+len(label)+1, row, hp, hpColor); err != nil {
		log.Fatalln("Couldn't write HUD hp", err)
	}
}

func (hud *HUD) renderPlayerStamina(world *World) {
	stColor := Red
	pct := hud.Player.ST.Percentage()
	switch {
	case pct >= 0.8:
		stColor = Green
	case pct >= 0.6:
		stColor = Yellow
	case pct >= 0.4:
		stColor = Orange
	default:
		stColor = Red
	}

	label := "Stamina:"
	st := fmt.Sprintf("%v/%v", hud.Player.ST.Current, hud.Player.ST.Max)
	if hud.Player.HP.Current == 0 {
		st += " *DEAD*"
	}

	row := hud.GetNextRow()
	if err := world.Window.PutString(hud.XPos, row, label, Yellow); err != nil {
		log.Fatalln("Couldn't write HUD st label", err)
	}

	if err := world.Window.PutString(hud.XPos+len(label)+1, row, st, stColor); err != nil {
		log.Fatalln("Couldn't write HUD st", err)
	}
}

func (hud *HUD) renderPlayerHeat(world *World) {
	htColor := Red
	pct := hud.Player.HT.Percentage()
	switch {
	case pct >= 0.8:
		htColor = Green
	case pct >= 0.6:
		htColor = Yellow
	case pct >= 0.4:
		htColor = Orange
	default:
		htColor = Red
	}

	label := "Heat:"
	ht := fmt.Sprintf("%v/%v", hud.Player.HT.Current, hud.Player.HT.Max)
	if hud.Player.HP.Current == 0 {
		ht += " *DEAD*"
	}

	row := hud.GetNextRow()
	if err := world.Window.PutString(hud.XPos, row, label, Yellow); err != nil {
		log.Fatalln("Couldn't write HUD ht label", err)
	}

	if err := world.Window.PutString(hud.XPos+len(label)+1, row, ht, htColor); err != nil {
		log.Fatalln("Couldn't write HUD ht", err)
	}
}

func (hud *HUD) renderPlayerLevel(world *World) {
	level := fmt.Sprintf("Level: %v (%v / %v)", hud.Player.Level, hud.Player.Experience, hud.Player.Level)
	world.Window.PutString(hud.XPos, hud.GetNextRow(), level, Yellow)
}

func (hud *HUD) renderTurnCount(world *World) {
	turnCount := fmt.Sprintf("Turn: %v", world.turnCount)
	world.Window.PutString(hud.XPos, hud.GetNextRow(), turnCount, Yellow)
}

func (hud *HUD) renderEquippedWeapon(world *World) {
	equipName := hud.Player.Equipment.Weapon.Name

	offsetY := hud.GetNextRow()
	offsetX := hud.XPos

	weaponStr := fmt.Sprintf("Weapon: %v", equipName)

	offsetX = hud.XPos
	hud.nextFreeRow += putWrappedText(world.Window, weaponStr, offsetX, offsetY, 0, 2, world.Window.Columns-offsetX, Yellow)
}

func (hud *HUD) renderItemDisplay(world *World) {
	hud.nextFreeRow += 3
	offsetY := hud.GetNextRow()
	offsetX := hud.XPos

	items := make([]Item, 0)
	for y := 0; y < world.CurrentLevel.Rows; y++ {
		for x := 0; x < world.CurrentLevel.Columns; x++ {
			if world.CurrentLevel.VisionMap.VisibilityAt(x, y) == Visible {
				tile := world.CurrentLevel.GetTile(x, y)
				if tile.Item != (Item{}) {
					items = append(items, tile.Item)
				}
			}
		}
	}
	for _, item := range items {
		world.Window.PutRune(hud.XPos, offsetY, item.Symbol, item.Color, gterm.NoColor)
		name := item.Name
		offsetX = hud.XPos
		hud.nextFreeRow += (putWrappedText(world.Window, name, offsetX, offsetY, 2, 4, world.Window.Columns-offsetX, Yellow) - 1)
		offsetY = hud.GetNextRow()
	}
	offsetX = hud.XPos
}

func (hud *HUD) Render(world *World) {
	hud.nextFreeRow = 0
	hud.renderPlayerName(world)
	hud.renderPlayerPosition(world)
	hud.renderPlayerHealth(world)
	hud.renderPlayerStamina(world)
	hud.renderPlayerHeat(world)
	hud.renderPlayerLevel(world)
	hud.renderTurnCount(world)
	hud.renderEquippedWeapon(world)
	hud.renderItemDisplay(world)
}
