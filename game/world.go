package main

import (
	"log"
	"math/rand"
	"path"
	"strconv"

	"github.com/MichaelTJones/pcg"

	"github.com/thomas-holmes/delivery-rl/game/monsters"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Position struct {
	X int
	Y int
}

type World struct {
	Window *gterm.Window

	turnCount uint64

	rng *pcg.PCG64

	Player *Creature

	CurrentLevel *Level

	MaxDepth     int
	Levels       []*Level
	LevelChanged bool

	CurrentUpdateTicks uint32
	CurrentTickDelta   uint32

	CameraCentered bool
	CameraWidth    int
	CameraHeight   int
	CameraOffsetX  int
	CameraOffsetY  int
	CameraX        int
	CameraY        int

	nextID      int
	nextLevelID int

	showScentOverlay bool

	bufferedInput InputEvent

	MenuStack  []Menu
	Animations []Animation

	GameOver bool
	QuitGame bool

	*GameLog
	Messaging
}

// GetNextID Generates a monotonically increasing Entity ID
func (world *World) GetNextID() int {
	world.nextID++
	return world.nextID
}

// PopInput get a queued input if there is one.
func (world *World) PopInput() InputEvent {
	input := world.bufferedInput
	world.bufferedInput = InputEvent{}

	return input
}

// AddInput queue an input for the game loop
func (world *World) AddInput(input InputEvent) {
	world.bufferedInput = input
}

// SetCurrentLevel update the worlds inner CurrentLevel pointer
func (world *World) SetCurrentLevel(index int) {
	world.CurrentLevel = world.Levels[index]
}

func (world *World) addInitialMonsters(level *Level) {
	monsters := monsters.LoadDefinitions(path.Join("assets", "definitions", "monsters.yaml"))

	for tries := 0; tries < level.MonsterDensity; tries++ {
		x := int(world.rng.Bounded(uint64(level.Columns)))
		y := int(world.rng.Bounded(uint64(level.Rows)))

		if level.CanStandOnTile(x, y) {
			index := rand.Intn(len(monsters))
			def := monsters[index]
			monster := NewMonster(x, y, def.Level, def.HP)
			monster.Name = def.Name
			monster.RenderGlyph = []rune(def.Glyph)[0]
			monster.RenderColor = def.Color.Color
			world.AddEntity(&monster, level)
		}
	}
}

// AddLevelFromCandidate constructs a real level from an intermediate level representation
func (world *World) AddLevelFromCandidate(level *CandidateLevel) {
	loadedLevel := LoadCandidateLevel(world.nextLevelID, level)
	world.nextLevelID++
	loadedLevel.Depth = len(world.Levels)

	world.Levels = append(world.Levels, &loadedLevel)

	levels := len(world.Levels)

	if levels > 1 {
		connectTwoLevels(world.Levels[levels-2], world.Levels[levels-1])
	}
	world.addInitialMonsters(&loadedLevel)
}

func (world *World) addPlayer(player *Creature, level *Level) {
	world.Player = player

	if world.CameraCentered {
		world.CameraX = player.X
		world.CameraY = player.Y
	}

	level.VisionMap.UpdateVision(world.Player.VisionDistance, world)
	level.ScentMap.UpdateScents(world)
}

func (world *World) addCreature(creature *Creature, level *Level) {
	creature.Depth = level.Depth

	if !level.CanStandOnTile(creature.X, creature.Y) {
		for _, t := range level.tiles {
			if !t.IsWall() && !(t.Creature != nil) {
				creature.X = t.X
				creature.Y = t.Y
				log.Printf("Creature position adjusted to (%v,%v)", creature.X, creature.Y)
				break
			}
		}
	}

	level.GetTile(creature.X, creature.Y).Creature = creature

	if creature.IsPlayer {
		world.addPlayer(creature, level)
	}
}

func (world *World) AddEntityToCurrentLevel(e Entity) {
	world.AddEntity(e, world.CurrentLevel)
}

func (world *World) AddEntity(e Entity, level *Level) {
	e.SetIdentity(world.GetNextID())
	log.Printf("Adding entity %+v", e)

	if n, ok := e.(Notifier); ok {
		n.SetMessageBus(world.messageBus)
	}

	if l, ok := e.(Listener); ok {
		world.messageBus.Subscribe(l)
	}

	level.Entities = append(level.Entities, e)

	if c, ok := e.(*Creature); ok {
		world.addCreature(c, level)
	}
}

func (world *World) RenderRuneAt(x int, y int, out rune, fColor sdl.Color, bColor sdl.Color) {
	err := world.Window.PutRune(x-world.CameraX+world.CameraOffsetX, y-world.CameraY+world.CameraOffsetY, out, fColor, bColor)
	if err != nil {
		log.Printf("Out of bounds %s", err)
	}
}

func (world *World) RenderStringAt(x int, y int, out string, color sdl.Color) {
	err := world.Window.PutString(x-world.CameraX+world.CameraOffsetX, y-world.CameraY+world.CameraOffsetY, out, color)
	if err != nil {
		log.Printf("Out of bounds %s", err)
	}
}
func (world *World) tidyAnimations() bool {
	insertionIndex := 0
	for _, a := range world.Animations {
		if !a.Done() {
			world.Animations[insertionIndex] = a
			insertionIndex++
		}
	}
	world.Animations = world.Animations[:insertionIndex]

	return len(world.Animations) > 0
}

func (world *World) tidyMenus() bool {
	insertionIndex := 0
	for _, menu := range world.MenuStack {
		if !menu.Done() {
			world.MenuStack[insertionIndex] = menu
			insertionIndex++
		}
	}
	world.MenuStack = world.MenuStack[:insertionIndex]

	return len(world.MenuStack) > 0
}

func (world *World) Update(input InputEvent) {
	world.AddInput(input)
	// log.Printf("Updating turn [%v]", world.turnCount)
	currentTicks := sdl.GetTicks()
	world.CurrentTickDelta = currentTicks - world.CurrentUpdateTicks
	world.CurrentUpdateTicks = currentTicks

	if world.tidyMenus() {
		currentMenu := world.MenuStack[len(world.MenuStack)-1]
		currentMenu.Update(world.PopInput())
	}

	if world.tidyAnimations() {
		for _, a := range world.Animations {
			a.Update(world.CurrentTickDelta)
		}
	}

	if world.GameOver {
		return
	}

	// Attempt at a hack
	if world.CurrentLevel.NextEntity >= len(world.CurrentLevel.Entities) {
		world.CurrentLevel.NextEntity = 0
	}
	for i := world.CurrentLevel.NextEntity; i < len(world.CurrentLevel.Entities); i++ {
		e := world.CurrentLevel.Entities[i]

		a, ok := e.(Actor)
		if !ok {
			continue
		}

		a.StartTurn()

		if a.CanAct() {
			if !a.Update(world.turnCount, world.PopInput(), world) {
				break
			}
		}

		// We've finished, so it is safe to advance
		if !a.CanAct() {
			a.EndTurn()
			world.CurrentLevel.NextEntity = i + 1
		}

		if c, ok := e.(*Creature); ok && c.IsPlayer {
			world.turnCount++
			world.CurrentLevel.VisionMap.UpdateVision(world.Player.VisionDistance, world)
			world.CurrentLevel.ScentMap.UpdateScents(world)
		}

		if world.LevelChanged {
			a.EndTurn()
			world.LevelChanged = false
			// Reset these or the player can end up in a spot where they have no energy but need input
			world.CurrentLevel.NextEntity = 0
			return // Is this the right thing to do? Or could we just break?
		}

		if world.CurrentLevel.NextEntity >= len(world.CurrentLevel.Entities) {
			world.CurrentLevel.NextEntity = 0
			i = -1
		}
	}
}

func (world *World) UpdateCamera() {
	if world.CameraCentered {
		world.CameraOffsetX = 30
		world.CameraOffsetY = 15
		world.CameraX = world.Player.X
		world.CameraY = world.Player.Y
	} else {
		world.CameraOffsetX = 0
		world.CameraOffsetY = 0
		world.CameraX = 0
		world.CameraY = 0
	}
}

// Render redrwas everything!
func (world *World) Render() {
	world.UpdateCamera()
	var minX, minY, maxX, maxY int
	if world.CameraCentered {
		minY, maxY = max(0, world.CameraY-(world.CameraHeight/2)), min(world.CurrentLevel.Rows, world.CameraY+(world.CameraHeight/2))
		minX, maxX = max(0, world.CameraX-(world.CameraWidth/2)), min(world.CurrentLevel.Columns, world.CameraX+(world.CameraWidth/2))
	} else {
		minY, maxY = 0, world.CurrentLevel.Rows
		minX, maxX = 0, world.CurrentLevel.Columns
	}
	for row := minY; row < maxY; row++ {
		for col := minX; col < maxX; col++ {
			tile := world.CurrentLevel.GetTile(col, row)

			visibility := world.CurrentLevel.VisionMap.VisibilityAt(col, row)
			tile.Render(world, visibility)
		}
	}

	world.GameLog.Render(world.Window)

	if world.showScentOverlay {
		world.OverlayScentMap()
	}

	for _, a := range world.Animations {
		a.Render(world)
	}

	// Render bottom to top
	for _, m := range world.MenuStack {
		m.Render(world.Window)
	}
}

func (world *World) OverlayVisionMap() {
	for y := 0; y < world.CurrentLevel.Rows; y++ {
		for x := 0; x < world.CurrentLevel.Columns; x++ {
			world.RenderRuneAt(x, y, []rune(strconv.Itoa(int(world.CurrentLevel.VisionMap.Map[y*world.CurrentLevel.Columns+x])))[0], Blue, gterm.NoColor)
		}
	}
}

// I'd maybe like this to be a bit better, but I cleaned up the weird coloration at the end.
// I don't really understand why it was doing what it did before but it's now more correct
// than it was.
var ScentColors = []sdl.Color{
	sdl.Color{R: 175, G: 50, B: 50, A: 200},
	sdl.Color{R: 225, G: 50, B: 25, A: 200},
	sdl.Color{R: 255, G: 0, B: 0, A: 200},
	sdl.Color{R: 100, G: 175, B: 50, A: 200},
	sdl.Color{R: 50, G: 255, B: 100, A: 200},
	sdl.Color{R: 0, G: 150, B: 175, A: 200},
	sdl.Color{R: 0, G: 50, B: 255, A: 200},
}

func (world *World) ToggleScentOverlay() {
	log.Printf("Scent Map Toggle pointer: %p", world.CurrentLevel.ScentMap)
	world.showScentOverlay = !world.showScentOverlay
}

func (world *World) OverlayScentMap() {
	for i, color := range ScentColors {
		if err := world.Window.PutRune(10+(i*2), 0, '.', White, color); err != nil {
			log.Printf("Couldn't draw overlay debug colors?")
		}
	}

	turn := world.turnCount
	for y := 0; y < world.CurrentLevel.Rows; y++ {
		for x := 0; x < world.CurrentLevel.Columns; x++ {
			scent := world.CurrentLevel.ScentMap.getScent(x, y)

			maxScent := float64((turn - 1) * 32)
			recent := float64((turn - 10) * 32)

			turnsAgo := int((maxScent - scent) / 32)
			if turnsAgo >= len(ScentColors) || turnsAgo < 0 {
				continue
			}
			distance := ((turn - uint64(turnsAgo)) * 32) - uint64(scent)

			bgColor := ScentColors[turnsAgo]

			bgColor.R /= 4
			bgColor.G /= 4
			bgColor.B /= 4

			if bgColor.R > bgColor.G && bgColor.R > bgColor.B {
				bgColor.R -= uint8(distance * 5)
			} else if bgColor.G > bgColor.B {
				bgColor.G -= uint8(distance * 5)
			} else {
				bgColor.B -= uint8(distance * 5)
			}
			if scent > 0 && scent > recent {
				world.RenderRuneAt(x, y, ' ', Purple, bgColor)
			}
		}
	}
}

func (world *World) RemoveEntity(entity Entity) {
	foundIndex := -1
	var foundEntity Entity
	for i, e := range world.CurrentLevel.Entities {
		if e.Identity() == entity.Identity() {
			foundIndex = i
			foundEntity = e
			break
		}
	}

	// This creates a but in some cases where a monster at the end of the entity list is killed.
	// This causes the game to kind of weirdly freeze forever. We need to fix how entities are removed
	// to fix it.
	if foundIndex > -1 {
		world.CurrentLevel.Entities = append(world.CurrentLevel.Entities[:foundIndex], world.CurrentLevel.Entities[foundIndex+1:]...)
	}
	if creature, ok := foundEntity.(*Creature); ok {
		world.CurrentLevel.GetTile(creature.X, creature.Y).Creature = nil
	}
}

func (world *World) MoveEntity(message MoveEntityMessage) {
	oldTile := world.CurrentLevel.GetTile(message.OldX, message.OldY)
	newTile := world.CurrentLevel.GetTile(message.NewX, message.NewY)
	newTile.Creature = oldTile.Creature
	oldTile.Creature = nil
}

// WARNING: This is has to perform a linear search which is less than ideal
// but I wanted ordered traversal, which you don't get with maps in go.
// Keep an eye on the performance of this.
func (world *World) GetEntity(id int) (Entity, bool) {
	for _, e := range world.CurrentLevel.Entities {
		if e.Identity() == id {
			return e, true
		}
	}
	return nil, false
}

func (world *World) Animating() bool {
	return len(world.Animations) > 0
}

func (world *World) UpdateAnimations() {
	currentTicks := sdl.GetTicks()
	world.CurrentTickDelta = currentTicks - world.CurrentUpdateTicks
	world.CurrentUpdateTicks = currentTicks

	if world.tidyAnimations() {
		for _, a := range world.Animations {
			a.Update(world.CurrentTickDelta)
		}
	}

}

func (world *World) AddAnimation(a Animation) {
	a.Start(world.CurrentUpdateTicks)
	world.Animations = append(world.Animations, a)
}

func (world *World) ShowEndGameMenu() {
	world.GameOver = true
	pop := NewEndGameMenu(10, 5, 40, 6, Red, "YOU ARE VERY DEAD", "I AM SO SORRY :(")
	world.Broadcast(ShowMenu, ShowMenuMessage{Menu: &pop})
}

func (world *World) Notify(message Message, data interface{}) {
	switch message {
	case ClearRegion:
		if d, ok := data.(ClearRegionMessage); ok {
			world.Window.ClearRegion(d.X, d.Y, d.W, d.H)
		}
	case MoveEntity:
		if d, ok := data.(MoveEntityMessage); ok {
			world.MoveEntity(d)
		}
	case KillEntity:
		if d, ok := data.(KillEntityMessage); ok {
			world.RemoveEntity(d.Defender)
		}
	case PlayerDead:
		world.ShowEndGameMenu()
	case PlayerFloorChange:
		if d, ok := data.(PlayerFloorChangeMessage); ok {
			if !d.Connected {
				break
			}
			world.RemoveEntity(world.Player)
			world.Player.X = d.DestX
			world.Player.Y = d.DestY
			world.CurrentLevel = world.getLevelById(d.DestLevelID)
			world.LevelChanged = true
			world.AddEntityToCurrentLevel(world.Player)
		}
	case ShowMenu:
		if d, ok := data.(ShowMenuMessage); ok {
			log.Printf("%T %+v", d.Menu, d.Menu)
			if n, ok := d.Menu.(Notifier); ok {
				n.SetMessageBus(world.messageBus)
			}
			d.Menu.Update(world.PopInput())
			world.MenuStack = append(world.MenuStack, d.Menu)
		}
	}
}

const (
	DefaultSeq uint64 = iota * 1000
)

func (world *World) BuildLevels() {
	var genFlags LevelGenFlag
	for i := 0; i < world.MaxDepth; i++ {
		switch i {
		case 0:
			genFlags = GenDownStairs
		case world.MaxDepth - 1:
			genFlags = GenUpStairs
		default:
			genFlags = GenDownStairs | GenUpStairs
		}

		level := GenLevel(world.rng, 100, 100, genFlags)
		world.AddLevelFromCandidate(level)
	}
	world.SetCurrentLevel(0)
}

func (w *World) getLevelById(id int) *Level {
	for _, l := range w.Levels {
		if l.ID == id {
			return l
		}
	}
	return nil
}

func (w *World) SaveGame() {
	w.Broadcast(SaveGame, SaveGameMessage{World: w})
}

func NewWorld(window *gterm.Window, centered bool, seed uint64) *World {
	world := &World{
		Window:         window,
		CameraCentered: centered,
		MaxDepth:       15,
		CameraX:        0,
		CameraY:        0,
		// TODO: Width/Height should probably be some function of the window dimensions
		CameraWidth:        56,
		CameraHeight:       25,
		CurrentUpdateTicks: sdl.GetTicks(),
		rng:                pcg.NewPCG64(),
	}

	world.rng.Seed(seed, DefaultSeq, seed*seed, DefaultSeq+1)

	world.messageBus = &MessageBus{}
	world.messageBus.Subscribe(world)

	world.GameLog = NewGameLog(0, window.Rows-4, 56, 3, world, world.messageBus)

	world.BuildLevels()

	return world
}
