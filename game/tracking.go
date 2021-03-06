package main

import (
	"sort"
)

type ScentMap struct {
	Columns int
	Rows    int
	Scent   []float64
}

type TrackCandidate struct {
	Position
	Scent float64
}

func (scentMap ScentMap) getScent(xPos int, yPos int) float64 {
	return scentMap.Scent[yPos*scentMap.Columns+xPos]
}

func (scentMap ScentMap) dirty(xPos int, yPos int, turn uint64, distance float64) {
	scentMap.Scent[yPos*scentMap.Columns+xPos] = float64(turn*32) - distance
}

func (scentMap ScentMap) track(turn uint64, xPos int, yPos int) []TrackCandidate {
	minX := max(0, xPos-1)
	maxX := min(scentMap.Columns, xPos+2)
	minY := max(0, yPos-1)
	maxY := min(scentMap.Rows, yPos+2)

	candidates := make([]TrackCandidate, 0, 8)
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			strength := scentMap.getScent(x, y)

			candidates = append(candidates,
				TrackCandidate{Position: Position{X: x, Y: y}, Scent: strength},
			)
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Scent > candidates[j].Scent })
	return candidates
}

func (scentMap *ScentMap) UpdateScents(world *World) {
	vision := world.CurrentLevel().VisionMap
	player := world.Player
	for y := 0; y < scentMap.Rows; y++ {
		for x := 0; x < scentMap.Columns; x++ {
			vision := vision.VisibilityAt(x, y)
			if vision == Visible && !world.CurrentLevel().GetTile(x, y).IsWall() {
				scentMap.dirty(x, y, world.turnCount, euclideanDistance(player.X, player.Y, x, y))
			}
		}
	}
}

func NewScentMap(columns int, rows int) *ScentMap {
	return &ScentMap{
		Columns: columns,
		Rows:    rows,
		Scent:   make([]float64, columns*rows, columns*rows),
	}
}
