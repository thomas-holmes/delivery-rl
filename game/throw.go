package main

import (
	"github.com/thomas-holmes/delivery-rl/game/controls"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

// This whole thing is going to duplicate spells pretty badly

func NewThrowPop(item Item, world *World) Menu {
	return &ThrowPop{
		Item:    item,
		World:   world,
		TargetX: world.Player.X,
		TargetY: world.Player.Y,
	}
}

type ThrowPop struct {
	Item

	*World

	TargetX int
	TargetY int

	targetVisible bool

	cursorColor sdl.Color
	lineColor   sdl.Color

	PopMenu
}

func (pop *ThrowPop) throwItem() {
	m.Broadcast(m.M{ID: PlayerThrowItem, Data: PlayerThrowItemMessage{World: pop.World, Item: pop.Item, TargetX: pop.TargetX, TargetY: pop.TargetY}})
	pop.done = true
}

func (pop *ThrowPop) adjustTarget(dX, dY int) {
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

func (pop *ThrowPop) Update(input controls.InputEvent) {
	pop.CheckCancel(input)

	switch input.Action() {
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
	case controls.Confirm:
		pop.throwItem()
	}
}

func (pop *ThrowPop) drawCursor(window *gterm.Window) {
	lineColor := pop.lineColor

	positions := PlotLine(pop.World.Player.X, pop.World.Player.Y, pop.TargetX, pop.TargetY)
	lineColor.A = 50
	for _, pos := range positions {
		pop.World.RenderRuneAt(pos.X, pos.Y, ' ', gterm.NoColor, lineColor)
	}

	pop.World.RenderRuneAt(pop.TargetX, pop.TargetY, ' ', gterm.NoColor, pop.lineColor)

}
func (pop *ThrowPop) Render(window *gterm.Window) {
	pop.drawCursor(window)
}
