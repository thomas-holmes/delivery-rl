package main

import (
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type TileKind int

const (
	Wall TileKind = iota
	Floor
	UpStair
	DownStair
)

type TileEffect int

const (
	None TileEffect = iota
	Greasy
	Peppery
)

const (
	WallGlyph      = '█'
	FloorGlyph     = '∙'
	UpStairGlyph   = '<'
	DownStairGlyph = '>'
)

func TileKindToGlyph(kind TileKind) rune {
	switch kind {
	case Wall:
		return WallGlyph
	case Floor:
		return FloorGlyph
	case UpStair:
		return UpStairGlyph
	case DownStair:
		return DownStairGlyph
	}

	return WallGlyph // Default to a wall for now, I guess.
}

func NewTile(x int, y int) Tile {
	return Tile{
		X:     x,
		Y:     y,
		Color: White,
	}
}

type Tile struct {
	X int
	Y int

	Color sdl.Color

	Creature *Creature
	Item     Item

	TileGlyph rune
	TileKind

	TileEffect
}

func (tile Tile) IsWall() bool {
	return tile.TileKind == Wall
}

func (tile *Tile) Render(world *World, visibility Visibility) {
	if visibility == Unseen {
		return
	}

	tile.RenderBackground(world, visibility)
	if tile.Creature != nil && visibility == Visible {
		tile.Creature.Render(world)
	}
}

func (tile Tile) RenderBackground(world *World, visibility Visibility) {
	var glyph rune
	var color sdl.Color
	renderFloor := true

	if tile.Creature != nil && visibility == Visible {
		renderFloor = false
	}

	if visibility == Visible {
		switch tile.TileEffect {
		case Greasy:
			glyph = grease
			color = GarlicGrease
			color.A /= 2
			renderFloor = false
			world.RenderRuneAt(tile.X, tile.Y, glyph, color, gterm.NoColor)
		}
	}

	if tile.Item.Symbol != 0 && tile.Creature == nil {
		glyph = tile.Item.Symbol
		color = tile.Item.Color
		renderFloor = true
	} else {
		glyph = tile.TileGlyph
		if tile.TileKind == UpStair || tile.TileKind == DownStair {
			color = Orange
		} else {
			color = tile.Color
		}

	}

	if visibility == Seen {
		color.R /= 2
		color.G /= 2
		color.B /= 2
	}

	if renderFloor {
		world.RenderRuneAt(tile.X, tile.Y, glyph, color, gterm.NoColor)
	}
}
