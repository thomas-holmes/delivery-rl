package main

import (
	"github.com/thomas-holmes/delivery-rl/game/controls"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

const WarpCost int = 4

type WarpPop struct {
	*World

	targetVisible bool
	TargetX       int
	TargetY       int

	lineColor   sdl.Color
	cursorColor sdl.Color

	PopMenu
}

func NewWarpPop(world *World) Menu {
	return &WarpPop{
		World:   world,
		TargetX: world.Player.X,
		TargetY: world.Player.Y,
		PopMenu: PopMenu{X: 65, Y: 32, W: 34, H: 3},
	}
}

func (pop *WarpPop) warp() {
	m.Broadcast(m.M{ID: PlayerWarp, Data: PlayerWarpMessage{World: pop.World, TargetX: pop.TargetX, TargetY: pop.TargetY}})
	pop.done = true
}

func (pop *WarpPop) adjustTarget(dX, dY int) {
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

func (pop *WarpPop) Update(input controls.InputEvent) {
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
		pop.warp()
	}
}

func (pop *WarpPop) drawCursor(window *gterm.Window) {
	lineColor := pop.lineColor

	positions := PlotLine(pop.World.Player.X, pop.World.Player.Y, pop.TargetX, pop.TargetY)
	lineColor.A = 50
	for _, pos := range positions {
		pop.World.RenderRuneAt(pos.X, pos.Y, ' ', gterm.NoColor, lineColor)
	}

	pop.World.RenderRuneAt(pop.TargetX, pop.TargetY, ' ', gterm.NoColor, pop.lineColor)
}

func (pop *WarpPop) RenderTooltip(window *gterm.Window) {
	window.PutString(pop.X+1, pop.Y+1, "Casting Warp...", White)
}

func (pop *WarpPop) Render(window *gterm.Window) {
	pop.drawCursor(window)

	pop.DrawBox(window, White)
	pop.RenderTooltip(window)
}
