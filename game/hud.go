package main

import (
	"fmt"
	"log"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
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
	hud := &HUD{
		Player:      player,
		World:       world,
		XPos:        xPos,
		YPos:        yPos,
		Width:       world.Window.Columns - xPos,
		Height:      world.Window.Rows - yPos,
		nextFreeRow: 0,
	}

	return hud
}

func (hud *HUD) GetNextRow() int {
	hud.nextFreeRow++
	return hud.nextFreeRow + hud.YPos - 1
}

func (hud *HUD) renderPlayerName(world *World) {
	world.Window.PutString(hud.XPos, hud.GetNextRow(), world.Player.Name, Yellow)
}

func (hud *HUD) renderPlayerPosition(world *World) {
	position := fmt.Sprintf("(%v, %v) - Level %v", hud.Player.X, hud.Player.Y, hud.Player.Depth+1)
	world.Window.PutString(hud.XPos, hud.GetNextRow(), position, Yellow)
}

var partial = rune('▌')
var partialRight = rune('▐')
var block = rune('█')

func drawBar(window *gterm.Window, pct float64, width int, x int, y int, color sdl.Color) {
	dimColor := color
	dimColor.R /= 2
	dimColor.G /= 2
	dimColor.B /= 2

	chunks := width * 2
	filledChunks := int(pct * float64(chunks))

	for i := 0; i < width; i++ {
		// Render dim
		leftChunk, rightChunk := i*2, i*2+1

		if rightChunk < filledChunks {
			window.PutRune(x+i, y, block, color, gterm.NoColor)
		} else {
			if leftChunk < filledChunks {
				window.PutRune(x+i, y, partial, color, gterm.NoColor)
			} else {
				window.PutRune(x+i, y, partial, dimColor, gterm.NoColor)
			}
			window.PutRune(x+i, y, partialRight, dimColor, gterm.NoColor)
		}
	}
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

	width := 20
	xOffset := hud.XPos + hud.Width - width - 1
	playerHP := hud.Player.HP.Percentage()

	drawBar(world.Window, playerHP, width, xOffset, row, hpColor)
}

func (hud *HUD) renderPlayerStamina(world *World) {
	stColor := Red
	pct := hud.Player.ST.Percentage()
	switch {
	case pct >= 0.8:
		stColor = Purple
	case pct >= 0.6:
		stColor = Blue
	case pct >= 0.4:
		stColor = Red
	default:
		stColor = Orange
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

	width := 20
	xOffset := hud.XPos + hud.Width - width - 1
	playerST := hud.Player.ST.Percentage()

	drawBar(world.Window, playerST, width, xOffset, row, stColor)
}

func (hud *HUD) renderPlayerHeat(world *World) {
	htColor := Red
	pct := hud.Player.HT.Percentage()
	switch {
	case pct >= 0.8:
		htColor = Red
	case pct >= 0.6:
		htColor = Orange
	case pct >= 0.4:
		htColor = Yellow
	default:
		htColor = Blue
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

	width := 20
	xOffset := hud.XPos + hud.Width - width - 1
	playerHT := hud.Player.HT.Percentage()

	drawBar(world.Window, playerHT, width, xOffset, row, htColor)
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
	for y := 0; y < world.CurrentLevel().Rows; y++ {
		for x := 0; x < world.CurrentLevel().Columns; x++ {
			if world.CurrentLevel().VisionMap.VisibilityAt(x, y) == Visible {
				tile := world.CurrentLevel().GetTile(x, y)
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
