package main

import (
	"encoding/gob"
	"io"

	"github.com/MichaelTJones/pcg"
	"github.com/veandco/go-sdl2/sdl"
)

/*
Maybe a bit different due to encoding/gob
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

	CurrentLevelIndex int
	CurrentLevelID    int

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
	ID int

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

	Creatures []ExportedCreatureV0
}

func exportTile(tile Tile) ExportedTileV0 {
	et := ExportedTileV0{}

	et.X = tile.X
	et.Y = tile.Y

	et.Color = tile.Color

	et.Item = tile.Item

	et.TileGlyph = tile.TileGlyph
	et.TileKind = tile.TileKind

	return et
}

func exportTiles(tiles []Tile) []ExportedTileV0 {
	eTiles := make([]ExportedTileV0, 0, len(tiles))
	for _, t := range tiles {
		eTiles = append(eTiles, exportTile(t))
	}

	return eTiles
}

func importTile(et ExportedTileV0) Tile {
	t := Tile{}

	t.X = et.X
	t.Y = et.Y

	t.Color = et.Color

	t.Item = et.Item

	t.TileGlyph = et.TileGlyph
	t.TileKind = et.TileKind

	return t

}

func importTiles(ets []ExportedTileV0) []Tile {
	tiles := make([]Tile, 0, len(ets))

	for _, et := range ets {
		tiles = append(tiles, importTile(et))
	}

	return tiles
}

func exportCreatures(entities []Entity) []ExportedCreatureV0 {
	ecs := make([]ExportedCreatureV0, 0, len(entities))

	for _, e := range entities {
		if c, ok := e.(*Creature); ok {
			ecs = append(ecs, ExportCreature(c))
		}
	}

	return ecs
}

func exportLevel(l Level) ExportedLevelV0 {
	return ExportedLevelV0{
		ID:        l.ID,
		Columns:   l.Columns,
		Rows:      l.Rows,
		VisionMap: *l.VisionMap,
		ScentMap:  *l.ScentMap,

		Tiles: exportTiles(l.Tiles),

		Stairs: l.Stairs,

		MonsterDensity: l.MonsterDensity,

		Depth: l.Depth,

		NextEntity: l.NextEntity,

		Creatures: exportCreatures(l.Entities),
	}
}

func exportLevels(ls []Level) []ExportedLevelV0 {
	levels := make([]ExportedLevelV0, 0, len(ls))

	for _, l := range ls {
		levels = append(levels, exportLevel(l))
	}

	return levels
}

func importLevel(el ExportedLevelV0) Level {
	l := Level{}

	l.ID = el.ID
	l.Columns = el.Columns
	l.Rows = el.Rows
	l.VisionMap = &el.VisionMap
	l.ScentMap = &el.ScentMap

	l.Tiles = importTiles(el.Tiles)

	l.Stairs = el.Stairs

	l.MonsterDensity = el.MonsterDensity

	l.Depth = el.Depth

	l.NextEntity = el.NextEntity

	return l
}

func importLevels(els []ExportedLevelV0) []Level {
	levels := make([]Level, 0, len(els))

	for _, l := range els {
		levels = append(levels, importLevel(l))
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
	return e.Encode(*s)
}

func (s *SaveV0) Decode(r io.Reader) error {
	d := gob.NewDecoder(r)
	return d.Decode(s)
}

func (s *SaveV0) SaveWorld(world *World) {
	s.TurnCount = world.turnCount
	s.Rng = *world.rng
	s.TurnCount = world.turnCount
	s.MaxDepth = int8(world.MaxDepth)

	s.CurrentLevelIndex = world.CurrentLevelIndex
	s.CurrentLevelID = world.CurrentLevelID

	s.Levels = exportLevels(world.Levels)

	s.NextID = int64(world.nextID)
	s.GameLog = world.GameLog.Messages

	s.Player = ExportCreature(world.Player)
}

func (s *SaveV0) Restore(w *World) {
	w.turnCount = s.TurnCount
	w.rng = &s.Rng
	w.MaxDepth = int(s.MaxDepth)

	w.Levels = importLevels(s.Levels)
	w.CurrentLevelIndex = s.CurrentLevelIndex
	w.CurrentLevelID = s.CurrentLevelID

	w.nextID = int(s.NextID)

	w.GameLog.Messages = s.GameLog

	w.Player = importCreature(s.Player)

	for _, l := range s.Levels {
		for _, c := range l.Creatures {
			creature := importCreature(c)

			w.AddEntity(creature, l.ID)
		}
	}
}

func importCreature(e ExportedCreatureV0) *Creature {
	c := Creature{}

	c.BasicEntity = e.BasicEntity

	c.IsPlayer = e.IsPlayer
	c.VisionDistance = e.VisionDistance

	c.Experience = e.Experience

	c.RenderGlyph = e.RenderGlyph
	c.RenderColor = e.RenderColor

	c.Depth = e.Depth

	c.Team = e.Team

	c.State = e.State

	c.X = e.X
	c.Y = e.Y

	c.Energy = e.Energy

	inventory := Inventory{
		Items: e.Inventory,
	}

	c.Inventory = inventory

	c.Equipment = e.Equipment

	c.Speed = e.Speed

	c.HP = e.HP
	c.ST = e.ST

	c.HT = e.HT

	c.Spells = e.Spells

	c.Level = e.Level

	c.Name = e.Name

	return &c
}

func init() {
	gob.Register(SaveV0{})
	gob.Register(ExportedCreatureV0{})
	gob.Register(Stair{})
	gob.Register(ExportedTileV0{})
	gob.Register(VisionMap{})
	gob.Register(ScentMap{})
	gob.Register(ExportedLevelV0{})
}
