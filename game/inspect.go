package main

import (
	"fmt"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type InspectionPop struct {
	World *World

	PopMenu

	TargetX int
	TargetY int

	targetVisible bool

	lineColor   sdl.Color
	cursorColor sdl.Color
}

func (pop *InspectionPop) adjustTarget(dX, dY int) {
	newX := pop.TargetX + dX
	newX = min(pop.World.CurrentLevel().Columns-1, max(0, newX))

	newY := pop.TargetY + dY
	newY = min(pop.World.CurrentLevel().Rows-1, max(0, newY))

	pop.targetVisible = pop.World.CurrentLevel().VisionMap.VisibilityAt(newX, newY) == Visible

	if pop.targetVisible {
		pop.cursorColor = Yellow
		pop.lineColor = Yellow
	} else {
		pop.cursorColor = Red
		pop.lineColor = Red
	}

	pop.TargetX = newX
	pop.TargetY = newY
}

func (pop *InspectionPop) Update(action controls.Action) {
	pop.CheckCancel(action)

	switch action {
	case controls.Up:
		pop.adjustTarget(0, -1)
	case controls.UpRight:
		pop.adjustTarget(1, -1)
	case controls.Right:
		pop.adjustTarget(1, 0)
	case controls.DownRight:
		pop.adjustTarget(1, 1)
	case controls.Down:
		pop.adjustTarget(0, 1)
	case controls.DownLeft:
		pop.adjustTarget(-1, 1)
	case controls.Left:
		pop.adjustTarget(-1, 0)
	case controls.UpLeft:
		pop.adjustTarget(-1, -1)
	}
}

func (pop *InspectionPop) RenderTileDescription(tile *Tile) {
	if pop.World.CurrentLevel().VisionMap.VisibilityAt(tile.X, tile.Y) == Unseen {
		return
	}
	yOffset := 1
	if c := tile.Creature; c != nil {
		xOffset := 0
		pop.World.Window.PutRune(pop.X+1+xOffset, pop.Y+yOffset, c.RenderGlyph, c.RenderColor, gterm.NoColor)
		xOffset += 2

		creatureLine1 := fmt.Sprintf("%v (%v/%v)", c.Name, c.HP.Current, c.HP.Max)
		pop.World.Window.PutString(pop.X+1+xOffset, pop.Y+yOffset, creatureLine1, White)
		xOffset += len(creatureLine1)

		if c.HasStatus(Confused) {
			status := "(Conf)"
			pop.World.Window.PutString(pop.X+1+xOffset, pop.Y+yOffset, status, Red)
			xOffset += len(status)
		}
		if c.HasStatus(Slow) {
			status := "(Slow)"
			pop.World.Window.PutString(pop.X+xOffset+1, pop.Y+yOffset, status, GarlicGrease)
			xOffset += len(status)
		}

		yOffset++

		xOffset = 0

		creatureLine2 := fmt.Sprintf("Weapon: %s", c.Equipment.Weapon.Name)
		pop.World.Window.PutString(pop.X+1+xOffset, pop.Y+yOffset, creatureLine2, White)

		yOffset++

	}
	if i := tile.Item; i != (Item{}) {
		xOffset := 0
		pop.World.Window.PutRune(pop.X+1+xOffset, pop.Y+yOffset, i.Symbol, i.Color, gterm.NoColor)

		var itemLine1 string
		if i.Power.Num > 0 {
			itemLine1 = fmt.Sprintf("- %v (%v)", i.Name, i.Power)
		} else {
			itemLine1 = fmt.Sprintf("- %v", i.Name)
		}
		yOffset += putWrappedText(pop.World.Window, itemLine1, pop.X+1, pop.Y+yOffset, 2, 4, pop.W-xOffset, White)
		yOffset++
	}
	{
		terrainLine1 := ""
		switch tile.TileKind {
		case Floor:
			terrainLine1 = "Stone floor"
		case Wall:
			terrainLine1 = "A solid rock wall"
		case UpStair:
			terrainLine1 = "Stairs leading up"
		case DownStair:
			terrainLine1 = "Stairs leading down"
		}

		if len(terrainLine1) > 0 {
			pop.World.Window.PutString(pop.X+1, pop.Y+yOffset, terrainLine1, White)
			yOffset++
		}
	}
}

// Maybe should interact with world/tiles than window directly
func (pop *InspectionPop) RenderCursor(window *gterm.Window) {
	lineColor := pop.lineColor

	positions := PlotLine(pop.World.Player.X, pop.World.Player.Y, pop.TargetX, pop.TargetY)
	lineColor.A = 50
	for _, pos := range positions {
		pop.World.RenderRuneAt(pos.X, pos.Y, ' ', gterm.NoColor, lineColor)
	}
	pop.World.RenderRuneAt(pop.TargetX, pop.TargetY, ' ', gterm.NoColor, pop.lineColor)
}

func (pop *InspectionPop) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	pop.DrawBox(window, White)

	pop.RenderCursor(window)

	tile := pop.World.CurrentLevel().GetTile(pop.TargetX, pop.TargetY)

	pop.RenderTileDescription(tile)
}
