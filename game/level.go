package main

import (
	"log"
	"sort"
)

type DistanceCandidate struct {
	Distance float64
	Position
}

func (level *Level) AddCreature(c *Creature) {
	c.Depth = level.Depth

	if !level.CanStandOnTile(c.X, c.Y) {
		for _, t := range level.Tiles {
			if !t.IsWall() && !(t.Creature != nil) {
				c.X = t.X
				c.Y = t.Y
				log.Printf("Creature position adjusted to (%v,%v)", c.X, c.Y)
				break
			}
		}
	}

	level.Entities = append(level.Entities, c)

	level.GetTile(c.X, c.Y).Creature = c
}

func (level Level) getStair(x int, y int) (Stair, bool) {
	for _, s := range level.Stairs {
		if s.X == x && s.Y == y {
			return s, true
		}
	}
	return Stair{}, false
}

func (level Level) GetTile(x int, y int) *Tile {
	return &level.Tiles[y*level.Columns+x]
}

func (level *Level) IsTileOccupied(x int, y int) bool {
	return level.GetTile(x, y).Creature != nil
}

func (level *Level) CanStandOnTile(column int, row int) bool {
	if level == nil {
		log.Panicf("Wtf is going on")
	}
	return !level.GetTile(column, row).IsWall() && !level.IsTileOccupied(column, row)
}

func (level *Level) GetCreatureAtTile(xPos int, yPos int) (*Creature, bool) {
	if creature := level.GetTile(xPos, yPos).Creature; creature != nil {
		return creature, true
	}
	return nil, false
}

// GetVisibleCreatures returns a slice of creatures sorted so that the first is the closest
// based on euclidean distance.
func (level *Level) GetVisibleCreatures(originX int, originY int) []*Creature {
	candidates := make([]DistanceCandidate, 0, 8)
	for y := 0; y < level.VisionMap.Rows; y++ {
		for x := 0; x < level.VisionMap.Columns; x++ {
			if level.VisionMap.VisibilityAt(x, y) == Visible {
				candidates = append(candidates, DistanceCandidate{Position: Position{X: x, Y: y}, Distance: euclideanDistance(originX, originY, x, y)})
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Distance < candidates[j].Distance })

	creatures := make([]*Creature, 0, len(candidates))
	for _, candidate := range candidates {
		if creature, ok := level.GetCreatureAtTile(candidate.X, candidate.Y); ok {
			creatures = append(creatures, creature)
		}
	}

	return creatures
}

// ClampedXY takes an origin point and a riuds to expand. Clamps to level dimensions
// Because I've written this 10000 times and I hate it
func (l *Level) ClampedXY(x, y, radius int) []Position {
	cols, rows := l.Columns, l.Rows
	minX, maxX := max(0, x-radius), min(cols-1, x+radius)
	minY, maxY := max(0, y-radius), min(rows-1, y+radius)

	positions := make([]Position, 0, radius*radius)

	for iy := minY; iy <= maxY; iy++ {
		for ix := minX; ix <= maxX; ix++ {
			pos := Position{X: ix, Y: iy}
			positions = append(positions, pos)
		}
	}

	return positions
}

type Stair struct {
	Down bool
	X    int
	Y    int

	Connected   bool
	DestX       int
	DestY       int
	DestLevelID int
}

type Level struct {
	ID int

	Columns   int
	Rows      int
	VisionMap *VisionMap
	ScentMap  *ScentMap
	Tiles     []Tile
	Stairs    []Stair

	MonsterDensity int

	Depth int // One Indexed

	NextEntity int
	Entities   []Entity
}

// connectTwoLevels connects multiple levels arbitrarily. If there is an uneven number
// of stair cases you will end up with a dead stair.
func connectTwoLevels(upper *Level, lower *Level) {
	for i, downStair := range upper.Stairs {
		if !downStair.Down || downStair.Connected {
			continue
		}

		for j, upStair := range lower.Stairs {
			if upStair.Down || upStair.Connected {
				continue
			}

			upper.Stairs[i].DestLevelID = lower.ID
			upper.Stairs[i].DestX = upStair.X
			upper.Stairs[i].DestY = upStair.Y
			upper.Stairs[i].Connected = true

			lower.Stairs[j].DestLevelID = upper.ID
			lower.Stairs[j].DestX = downStair.X
			lower.Stairs[j].DestY = downStair.Y
			lower.Stairs[j].Connected = true

			break
		}
	}
}

func LoadCandidateLevel(id int, candidate *CandidateLevel) Level {
	level := Level{ID: id, Depth: candidate.depth}

	tiles := make([]Tile, 0, len(candidate.tiles))

	var stairs []Stair

	for y := 0; y < candidate.H; y++ {
		for x := 0; x < candidate.W; x++ {
			tile, cTile := NewTile(x, y), candidate.tiles[y*candidate.W+x]
			tile.TileKind = cTile.TileKind
			tile.TileGlyph = TileKindToGlyph(cTile.TileKind)
			tile.Item = cTile.Item

			switch tile.TileKind {
			case UpStair:
				stair := Stair{
					X:    x,
					Y:    y,
					Down: false,
				}
				stairs = append(stairs, stair)
			case DownStair:
				stair := Stair{
					X:    x,
					Y:    y,
					Down: true,
				}
				stairs = append(stairs, stair)
			}
			tiles = append(tiles, tile)
		}
	}

	level.Columns = candidate.W
	level.Rows = candidate.H
	level.Tiles = tiles
	level.Stairs = stairs
	level.MonsterDensity = 12

	level.VisionMap = NewVisionMap(level.Columns, level.Rows)
	level.ScentMap = NewScentMap(level.Columns, level.Rows)

	return level
}
