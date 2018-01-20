package main

import (
	"encoding/gob"
	"io"

	"github.com/MichaelTJones/pcg"
	"github.com/veandco/go-sdl2/sdl"
)

/*
|-4 Bytes-||--ArbitraryBytes--||--ArbitraryBytes--||--ArbitraryBytes--|
|--Magic--||----FieldBytes----||----FieldBytes----||----FieldBytes----|
*/
type SaveVersion struct {
	Magic int32
}

type SaveV0 struct {
	SaveVersion

	// World Fields
	TurnCount uint64    // World
	Rng       pcg.PCG64 // World
	MaxDepth  int8      // World. Down-Convert to int8

	// Perhaps we store RNG state for level generation and regenerate? Need versioned level gen code then :(
	Levels []ExportedLevelV0

	NextID int64

	GameLog []string

	Player ExportedCreatureV0
}

func NewSaveV0() SaveV0 {
	return SaveV0{
		SaveVersion: SaveVersion{Magic: 0},
	}
}

type ExportedCreatureV0 struct {
	BasicEntity

	IsPlayer       bool
	VisionDistance int

	Experience int

	RenderGlyph rune
	RenderColor sdl.Color

	Depth int

	Team

	State MonsterBehavior

	X int
	Y int

	Energy

	Inventory []Item

	Equipment Equipment

	Speed int

	HP Resource // health
	ST Resource // stamina

	HT Resource // heat

	Spells []Spell

	Level int

	Name string
}

type ExportedTileV0 struct {
	X int
	Y int

	Color sdl.Color

	Item Item

	TileGlyph rune
	TileKind
}

type ExportedLevelV0 struct {
	Columns   int
	Rows      int
	VisionMap VisionMap
	ScentMap  ScentMap

	Tiles  []ExportedTileV0
	Stairs []Stair

	MonsterDensity int

	Depth int

	NextEntity int

	/* Figure this out. Fortunately I think we only have Creatures here
	Entities   []Entity
	*/
}

func exportTiles(tiles []Tile) []ExportedTileV0 {
	eTiles := make([]ExportedTileV0, 0, len(tiles))

	return eTiles
}

func exportLevel(l *Level) ExportedLevelV0 {
	return ExportedLevelV0{
		Columns:   l.Columns,
		Rows:      l.Rows,
		VisionMap: *l.VisionMap,
		ScentMap:  *l.ScentMap,

		Tiles: exportTiles(l.tiles),

		Stairs: l.stairs,

		MonsterDensity: l.MonsterDensity,

		Depth: l.Depth,

		NextEntity: l.NextEntity,

		/* Figure this out. Fortunately I think we only have Creatures here
		Entities   []Entity
		*/
	}
}

func exportLevels(ls []*Level) []ExportedLevelV0 {
	levels := make([]ExportedLevelV0, 0, len(ls))

	for _, l := range ls {
		levels = append(levels, exportLevel(l))
	}

	return levels
}

func (s *ExportedLevelV0) Encode(w io.Writer) error {
	e := gob.NewEncoder(w)
	return e.Encode(s)
}

func (s *ExportedLevelV0) Decode(r io.Reader) error {
	d := gob.NewDecoder(r)
	return d.Decode(s)
}

func ExportCreature(c *Creature) ExportedCreatureV0 {
	items := make([]Item, 0, len(c.Items))
	for _, i := range c.Items {
		items = append(items, i)
	}

	return ExportedCreatureV0{
		BasicEntity: c.BasicEntity,

		IsPlayer:       c.IsPlayer,
		VisionDistance: c.VisionDistance,

		Experience: c.Experience,

		RenderGlyph: c.RenderGlyph,
		RenderColor: c.RenderColor,

		Depth: c.Depth,

		Team: c.Team,

		State: c.State,

		X: c.X,
		Y: c.Y,

		Energy: c.Energy,

		Inventory: items,

		Equipment: c.Equipment,

		Speed: c.Speed,

		HP: c.HP,
		ST: c.ST,

		HT: c.HT,

		Spells: c.Spells,

		Level: c.Level,

		Name: c.Name,
	}
}

func (s *ExportedCreatureV0) Encode(w io.Writer) error {
	e := gob.NewEncoder(w)
	return e.Encode(s)
}

func (s *ExportedCreatureV0) Decode(r io.Reader) error {
	d := gob.NewDecoder(r)
	return d.Decode(s)
}

func (s *SaveV0) Encode(w io.Writer) error {
	e := gob.NewEncoder(w)
	return e.Encode(s)
}

func (s *SaveV0) SaveWorld(world *World) {
	s.TurnCount = world.turnCount
	s.Rng = *world.rng
	s.TurnCount = world.turnCount
	s.MaxDepth = int8(world.MaxDepth)

	s.Levels = exportLevels(world.Levels)

	s.NextID = int64(world.nextID)
	s.GameLog = world.GameLog.Messages

	s.Player = ExportCreature(world.Player)
}

func init() {
	gob.Register(SaveV0{})
	gob.Register(ExportedCreatureV0{})
	gob.Register(ExportedLevelV0{})
}
