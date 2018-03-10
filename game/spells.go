package main

import (
	"fmt"
	"log"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type SpellShape int

const (
	Square SpellShape = iota
	Line
	Cone
	Point
)

type Spell struct {
	Name        string
	Description string
	Range       int
	Power       int
	Cost        int
	Shape       SpellShape
	Hits        int
	Size        int
}

var DefaultSpells = []Spell{
	// Hits 1
	Spell{Name: "Fire Bolt", Description: "Launches a small ball of fire at an oponnent", Range: 8, Power: 2, Cost: 1, Shape: Square, Size: 0, Hits: 1},

	// Hits 1, 3 times
	Spell{Name: "Magic Missile", Description: "Fires 3 magic missiles, which are guaranteed to strike their target", Range: 8, Power: 1, Cost: 2, Shape: Square, Size: 0, Hits: 3},

	// Hits 7 wide cone
	Spell{Name: "Cone of Cold", Description: "A cone of ice extends from your hands, damaging and chilling all in front of you", Range: 1, Power: 3, Cost: 3, Shape: Cone, Size: 3, Hits: 1},

	// Hits 5x5
	Spell{Name: "Fire Ball", Description: "Hurls a large ball of fire at a group of oponnents", Range: 8, Power: 4, Cost: 4, Shape: Square, Size: 2, Hits: 1},

	// Movement Spell
	Spell{Name: "Blink", Description: "Teleport to a nearby unoccupied space you can see", Range: 8, Power: 0, Cost: 1, Shape: Point, Size: 1, Hits: 0},
}

type SpellPop struct {
	World *World

	PopMenu
}

func (pop *SpellPop) castSpell(index int) {
	if index < len(pop.World.Player.Spells) {
		spell := pop.World.Player.Spells[index]

		if pop.World.Player.CanCast(spell) {
			log.Printf("Casting spell %+v", spell)

			pop.done = true // Maybe needs to be before player.CastSpell

			m.Broadcast(m.M{ID: SpellTarget, Data: SpellTargetMessage{Spell: spell, World: pop.World}})
		}
	}
}

func (pop *SpellPop) Update(input controls.InputEvent) {
	pop.CheckCancel(input)

	switch e := input.Event.(type) {
	case sdl.KeyDownEvent:
		k := e.Keysym.Sym
		switch {
		case k >= sdl.K_a && k <= sdl.K_z:
			pop.castSpell(int(k - sdl.K_a))
		}

	}
}

func (pop SpellPop) renderItem(index int, row int, window *gterm.Window) int {
	offsetY := row
	offsetX := pop.X + 1

	spell := pop.World.Player.Spells[index]

	itemColor := Grey
	if pop.World.Player.CanCast(spell) {
		itemColor = White
	}
	selectionStr := fmt.Sprintf("%v - ", string('a'+index))

	window.PutString(offsetX, offsetY, selectionStr, itemColor)

	name := spell.Name

	offsetY += putWrappedText(window, name, offsetX, offsetY, len(selectionStr), 2, pop.W-offsetX+pop.X-1, itemColor)
	return offsetY
}

func (pop SpellPop) Render(window *gterm.Window) {
	if err := window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Printf("(%v,%v) (%v,%v)", pop.X, pop.Y, pop.W, pop.H)
		log.Println("Failed to clear region for spell menu", err)
	}

	nextRow := pop.Y + 1
	for i := 0; i < len(pop.World.Player.Spells); i++ {
		nextRow = pop.renderItem(i, nextRow, window)
	}

	pop.DrawBox(window, White)
}

// Targeting

type SpellTargeting struct {
	World *World
	Spell Spell

	TargetX int
	TargetY int

	targetVisible bool
	initialized   bool
	creatures     []*Creature
	creatureIndex int

	distance int

	cursorColor sdl.Color
	lineColor   sdl.Color

	PopMenu
}

func (pop *SpellTargeting) setInitialState() {
	if !pop.initialized {
		player := pop.World.Player
		pop.creatures = pop.World.CurrentLevel().GetVisibleCreatures(player.X, player.Y)
		pop.initialized = true
		for i, c := range pop.creatures {
			if c.Team != player.Team {
				if pop.Spell.Shape == Cone {
					positions := PlotLine(pop.World.Player.X, pop.World.Player.Y, c.X, c.Y)
					closestPosition := positions[1]
					pop.TargetX = closestPosition.X
					pop.TargetY = closestPosition.Y
					return
				} else {
					pop.TargetX = c.X
					pop.TargetY = c.Y
					pop.creatureIndex = i
					return
				}
			}
		}
	}
}

func (pop *SpellTargeting) Update(input controls.InputEvent) {
	player := pop.World.Player
	pop.setInitialState()
	newX, newY := pop.TargetX, pop.TargetY
	switch e := input.Event.(type) {
	case sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_RETURN:
			if pop.distance <= pop.Spell.Range {
				pop.done = true
				pop.World.Player.CastSpell(pop.Spell, pop.World, pop.TargetX, pop.TargetY)
			} else {
				fmt.Println("Can't cast, out of range.")
			}
		case sdl.K_ESCAPE:
			pop.done = true
		case sdl.K_h:
			newX = pop.TargetX - 1
		case sdl.K_j:
			newY = pop.TargetY + 1
		case sdl.K_k:
			newY = pop.TargetY - 1
		case sdl.K_l:
			newX = pop.TargetX + 1
		case sdl.K_b:
			newX, newY = pop.TargetX-1, pop.TargetY+1
		case sdl.K_n:
			newX, newY = pop.TargetX+1, pop.TargetY+1
		case sdl.K_y:
			newX, newY = pop.TargetX-1, pop.TargetY-1
		case sdl.K_u:
			newX, newY = pop.TargetX+1, pop.TargetY-1
		case sdl.K_EQUALS:
			pop.creatureIndex = (pop.creatureIndex + 1) % len(pop.creatures)
			newX, newY = pop.creatures[pop.creatureIndex].X, pop.creatures[pop.creatureIndex].Y
		case sdl.K_MINUS:
			pop.creatureIndex = (pop.creatureIndex - 1)
			if pop.creatureIndex < 0 {
				pop.creatureIndex = len(pop.creatures) - 1
			}
			newX, newY = pop.creatures[pop.creatureIndex].X, pop.creatures[pop.creatureIndex].Y
		}
	}

	pop.checkTargetVisibility(newX, newY)

	if (newX != pop.TargetX || newY != pop.TargetY) &&
		(newX > 0 && newX < pop.World.CurrentLevel().Columns) &&
		(newY > 0 && newY < pop.World.CurrentLevel().Rows) &&
		pop.targetVisible {
		pop.TargetX = newX
		pop.TargetY = newY
	}

	pop.distance = squareDistance(player.X, player.Y, pop.TargetX, pop.TargetY)

	if pop.distance > pop.Spell.Range || !pop.targetVisible {
		pop.cursorColor = Red
		pop.lineColor = Red
	} else {
		pop.cursorColor = Yellow
		pop.lineColor = White
	}
}

func (pop *SpellTargeting) checkTargetVisibility(newX int, newY int) {
	tile := pop.World.CurrentLevel().GetTile(newX, newY)
	if tile.IsWall() || pop.World.CurrentLevel().VisionMap.VisibilityAt(newX, newY) != Visible {
		pop.targetVisible = false
	} else {
		pop.targetVisible = true
	}
}

func (pop SpellTargeting) renderSquareCursor() {
	cursorColor := pop.cursorColor

	spell := pop.Spell

	minX := max(pop.TargetX-spell.Size, 0)
	maxX := min(pop.TargetX+spell.Size, pop.World.CurrentLevel().Columns)

	minY := max(pop.TargetY-spell.Size, 0)
	maxY := min(pop.TargetY+spell.Size, pop.World.CurrentLevel().Rows)

	cursorColor.A = 125
	for y := minY; y < maxY+1; y++ {
		for x := minX; x < maxX+1; x++ {
			pop.World.RenderRuneAt(x, y, ' ', gterm.NoColor, cursorColor)
		}
	}

	cursorColor.A = 200
	pop.World.RenderRuneAt(pop.TargetX, pop.TargetY, ' ', gterm.NoColor, cursorColor)
}

func conePositions(pX, pY, x0, y0, size int) []Position {
	endX, endY := x0, y0
	if endX > pX {
		endX += size
	} else if endX < pX {
		endX -= size
	}

	if endY > pY {
		endY += size
	} else if endY < pY {
		endY -= size
	}

	coords := PlotLine(x0, y0, endX, endY)
	octant := computeOctant(x0, y0, endX, endY)
	conePositions := make([]Position, 0)
	for i, pos := range coords {
		oX, oY := toOctantZero(octant, pos.X, pos.Y)
		for j := 0; j < i+1; j++ {
			if j == 0 {
				realX, realY := fromOctantZero(octant, oX, oY)
				conePositions = append(conePositions, Position{X: realX, Y: realY})
			} else {
				realX, realY := fromOctantZero(octant, oX, oY+j)
				conePositions = append(conePositions, Position{X: realX, Y: realY})
				realX, realY = fromOctantZero(octant, oX, oY-j)
				conePositions = append(conePositions, Position{X: realX, Y: realY})
			}
		}

	}

	return conePositions
}

func (pop SpellTargeting) renderConeCursor() {
	cursorColor := pop.cursorColor
	spell := pop.Spell
	player := pop.World.Player

	cursorColor.A = 125
	for _, pos := range conePositions(player.X, player.Y, pop.TargetX, pop.TargetY, spell.Size) {
		pop.World.RenderRuneAt(pos.X, pos.Y, ' ', gterm.NoColor, cursorColor)
	}

	cursorColor.A = 200
	pop.World.RenderRuneAt(pop.TargetX, pop.TargetY, ' ', gterm.NoColor, cursorColor)

}

func (pop SpellTargeting) renderPointCursor() {
	cursorColor := pop.cursorColor

	cursorColor.A = 200
	pop.World.RenderRuneAt(pop.TargetX, pop.TargetY, ' ', gterm.NoColor, cursorColor)
}

func (pop SpellTargeting) Render(window *gterm.Window) {
	lineColor := pop.lineColor

	lineColor.A = 50

	positions := PlotLine(pop.World.Player.X, pop.World.Player.Y, pop.TargetX, pop.TargetY)
	for _, pos := range positions {
		pop.World.RenderRuneAt(pos.X, pos.Y, ' ', gterm.NoColor, lineColor)
	}

	switch pop.Spell.Shape {
	case Cone:
		pop.renderConeCursor()
	case Line:
		fallthrough
	case Square:
		pop.renderSquareCursor()
	case Point:
		pop.renderPointCursor()
	}
}
