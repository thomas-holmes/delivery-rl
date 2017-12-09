package main

import (
	"time"
)

type ScentMap struct {
	columns int
	rows    int
	scent   []float64
}

func (scentMap ScentMap) getScent(xPos int, yPos int) float64 {
	return scentMap.scent[yPos*scentMap.columns+xPos]
}

func (scentMap ScentMap) dirty(xPos int, yPos int, turn uint64, distance float64) {
	scentMap.scent[yPos*scentMap.columns+xPos] = float64(turn*32) - distance
}

func (scentMap ScentMap) track(turn uint64, xPos int, yPos int) []Position {
	minX := max(0, xPos-1)
	maxX := min(scentMap.columns, xPos+2)
	minY := max(0, yPos-1)
	maxY := min(scentMap.columns, yPos+2)

	candidates := make([]Position, 0, 8)
	recent := float64((turn - 50) * 32)
	strongest := recent
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			strength := scentMap.getScent(x, y)

			if strength > strongest {
				candidates = candidates[:0]
				strongest = strength
			}

			if strength == strongest {
				candidates = append(candidates,
					Position{XPos: x, YPos: y},
				)
			}
		}
	}
	return candidates
}

func (scentMap ScentMap) UpdateScents(turn uint64, world World) {
	vision := world.CurrentLevel.VisionMap
	player := world.Player
	defer timeMe(time.Now(), "ScentMap.UpdateScents")
	for y := 0; y < scentMap.rows; y++ {
		for x := 0; x < scentMap.columns; x++ {
			vision := vision.VisibilityAt(x, y)
			if vision == Visible && !world.GetTile(x, y).IsWall() {
				scentMap.dirty(x, y, turn, distance(player.X, player.Y, x, y))
			}
		}
	}
}

func NewScentMap(columns int, rows int) ScentMap {
	return ScentMap{
		columns: columns,
		rows:    rows,
		scent:   make([]float64, columns*rows, columns*rows),
	}
}
