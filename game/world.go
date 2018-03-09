package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/items"
	m "github.com/thomas-holmes/delivery-rl/game/messages"

	"github.com/MichaelTJones/pcg"

	"github.com/thomas-holmes/delivery-rl/game/monsters"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

const MaxDepth int = 10

type Position struct {
	X int
	Y int
}

type World struct {
	Window *gterm.Window

	turnCount uint64

	rng *pcg.PCG64

	Player *Creature

	HUD HUD

	CurrentLevelID    int
	CurrentLevelIndex int

	MaxDepth     int
	Levels       []Level
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

	Input controls.InputEvent

	MenuStack  []Menu
	Animations []Animation

	GameOver bool
	QuitGame bool

	*GameLog
}

// CurrentLevel returns a pointer to the current level. Don't store this pointer!
func (world *World) CurrentLevel() *Level {
	if world.CurrentLevelIndex >= len(world.Levels) {
		log.Panicf("Have level index out of bounds. Len (%v), index (%v)", len(world.Levels), world.CurrentLevelIndex)
	}
	return &world.Levels[world.CurrentLevelIndex]
}

// GetNextID Generates a monotonically increasing Entity ID
func (world *World) GetNextID() int {
	world.nextID++
	return world.nextID
}

// SetCurrentLevel update the worlds inner CurrentLevel pointer
func (world *World) SetCurrentLevel(id int) {
	for i := range world.Levels {
		if world.Levels[i].ID == id {
			world.CurrentLevelIndex = i
			world.CurrentLevelID = id
			return
		}
	}
}

func buildMonsterFromDefinition(def monsters.Definition) *Creature {
	monster := NewMonster(0, 0, def.Level, def.HP)
	monster.Name = def.Name
	monster.RenderGlyph = []rune(def.Glyph)[0]
	monster.RenderColor = def.Color.Color

	if def.Weapon.Name != "" {
		monster.Equipment.Weapon = produceItem(items.Definition(def.Weapon))
	}

	return monster
}

func (world *World) addInitialMonsters(level *Level) {
	monsterCollection := monsters.GetCollection("monsters")

	for tries := 0; tries < level.MonsterDensity; tries++ {
		x := int(world.rng.Bounded(uint64(level.Columns)))
		y := int(world.rng.Bounded(uint64(level.Rows)))

		if level.CanStandOnTile(x, y) {
			def := monsters.GetPowerBoundedMonster(world.rng, monsterCollection, level.Depth+1)
			monster := buildMonsterFromDefinition(def)
			monster.X = x
			monster.Y = y
			world.AddEntity(monster, level.ID)
		}
	}
}

func (world *World) addDragonToLevel(level *Level) {
	dragon := buildMonsterFromDefinition(monsters.Dragon)
	dragon.Team = NeutralTeam
	dragon.IsDragon = true

	world.AddEntity(dragon, level.ID)
}

// AddLevelFromCandidate constructs a real level from an intermediate level representation
func (world *World) AddLevelFromCandidate(level *CandidateLevel) {
	loadedLevel := LoadCandidateLevel(world.nextLevelID, level)
	world.nextLevelID++
	loadedLevel.Depth = len(world.Levels)

	world.Levels = append(world.Levels, loadedLevel)

	levels := len(world.Levels)
	levelIndex := levels - 1

	if levels > 1 {
		connectTwoLevels(&world.Levels[levelIndex-1], &world.Levels[levelIndex])
	}
	log.Println("Adding monsters to level with depth", world.Levels[levelIndex].Depth)
	world.addInitialMonsters(&world.Levels[levelIndex])

	if levels == world.MaxDepth {
		world.addDragonToLevel(&world.Levels[levelIndex])
	}

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
	level.AddCreature(creature)

	if creature.IsPlayer {
		log.Printf("Adding the player")
		world.addPlayer(creature, level)
		log.Printf("%+v", level.Entities)
	}
}

func (world *World) AddEntityToCurrentLevel(e Entity) {
	log.Printf("Adding an entity (%v) to (%v)", e.Identity(), world.CurrentLevelID)
	world.AddEntity(e, world.CurrentLevelID)
}

func (world *World) AddEntity(e Entity, levelID int) {
	level := world.LevelByID(levelID)
	if e.Identity() == 0 {
		e.SetIdentity(world.GetNextID())
		log.Printf("giving it an ID! %v", e.Identity())
	}
	log.Printf("Adding entity %+v", e)

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

func (world *World) Update() {
	currentTicks := sdl.GetTicks()
	world.CurrentTickDelta = currentTicks - world.CurrentUpdateTicks
	world.CurrentUpdateTicks = currentTicks

	// Check for Game Over
	{
		if world.GameOver {
			return
		}
	}

	// Update Animations
	{
		if world.tidyAnimations() {
			for _, a := range world.Animations {
				a.Update(world.CurrentTickDelta)
			}
			return
		}
	}

	// Update Menus
	{
		if world.tidyMenus() {
			currentMenu := world.MenuStack[len(world.MenuStack)-1]
			currentMenu.Update(world.Input)
			return
		}
	}

	// Update All Actors?
	{
		for {
			entities := world.CurrentLevel().Entities
			if world.CurrentLevel().NextEntity >= len(entities) {
				world.CurrentLevel().NextEntity = 0
			}
			currentEntity := world.CurrentLevel().NextEntity
			e := entities[currentEntity]

			a, ok := e.(Actor)
			if !ok {
				world.CurrentLevel().NextEntity++
				continue
			}
			a.StartTurn()

			if !a.CanAct() {
				a.EndTurn()
				world.CurrentLevel().NextEntity++
				continue
			}

			advancedTurn := a.Update(world.turnCount, world.Input, world)

			// bit of a hack, copied from below instead of rewritten
			if world.LevelChanged {
				a.EndTurn()
				world.LevelChanged = false
				world.CurrentLevel().NextEntity = 0
				return // Is this the right thing to do? Or could we just break?
			}

			if c, ok := a.(*Creature); ok {
				if c.IsPlayer {
					if advancedTurn {
						world.turnCount++
					}
					world.CurrentLevel().VisionMap.UpdateVision(world.Player.VisionDistance, world)
					world.CurrentLevel().ScentMap.UpdateScents(world)
				}
			}

			if a.CanAct() {
				// Needs more input
				return
			} else {
				world.CurrentLevel().NextEntity++
				world.Input = controls.InputEvent{}
			}
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
		minY, maxY = max(0, world.CameraY-(world.CameraHeight/2)), min(world.CurrentLevel().Rows, world.CameraY+(world.CameraHeight/2))
		minX, maxX = max(0, world.CameraX-(world.CameraWidth/2)), min(world.CurrentLevel().Columns, world.CameraX+(world.CameraWidth/2))
	} else {
		minY, maxY = 0, world.CurrentLevel().Rows
		minX, maxX = 0, world.CurrentLevel().Columns
	}
	for row := minY; row < maxY; row++ {
		for col := minX; col < maxX; col++ {
			tile := world.CurrentLevel().GetTile(col, row)

			visibility := world.CurrentLevel().VisionMap.VisibilityAt(col, row)
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
	for y := 0; y < world.CurrentLevel().Rows; y++ {
		for x := 0; x < world.CurrentLevel().Columns; x++ {
			world.RenderRuneAt(x, y, []rune(strconv.Itoa(int(world.CurrentLevel().VisionMap.Map[y*world.CurrentLevel().Columns+x])))[0], Blue, gterm.NoColor)
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
	log.Printf("Scent Map Toggle pointer: %p", world.CurrentLevel().ScentMap)
	world.showScentOverlay = !world.showScentOverlay
}

func (world *World) OverlayScentMap() {
	for i, color := range ScentColors {
		if err := world.Window.PutRune(10+(i*2), 0, '.', White, color); err != nil {
			log.Printf("Couldn't draw overlay debug colors?")
		}
	}

	turn := world.turnCount
	for y := 0; y < world.CurrentLevel().Rows; y++ {
		for x := 0; x < world.CurrentLevel().Columns; x++ {
			scent := world.CurrentLevel().ScentMap.getScent(x, y)

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
	for i, e := range world.CurrentLevel().Entities {
		if e.Identity() == entity.Identity() {
			foundIndex = i
			foundEntity = e
			break
		}
	}

	if foundIndex > -1 {
		world.CurrentLevel().Entities = append(world.CurrentLevel().Entities[:foundIndex], world.CurrentLevel().Entities[foundIndex+1:]...)
		if foundIndex > world.CurrentLevel().NextEntity {
			world.CurrentLevel().NextEntity--
		}
	}
	if creature, ok := foundEntity.(*Creature); ok {
		world.CurrentLevel().GetTile(creature.X, creature.Y).Creature = nil
	}
}

func (world *World) MoveEntity(message MoveEntityMessage) {
	oldTile := world.CurrentLevel().GetTile(message.OldX, message.OldY)
	newTile := world.CurrentLevel().GetTile(message.NewX, message.NewY)
	newTile.Creature = oldTile.Creature
	oldTile.Creature = nil
}

// WARNING: This is has to perform a linear search which is less than ideal
// but I wanted ordered traversal, which you don't get with maps in go.
// Keep an eye on the performance of this.
func (world *World) GetEntity(id int) (Entity, bool) {
	for _, e := range world.CurrentLevel().Entities {
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
	pop := NewEndGameMenu(world, 5, 5, 50, 6, Red, "YOU ARE VERY DEAD", "I AM SO SORRY :(")
	m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: &pop}})
}

func (world *World) ShowFoodSpoiledMenu() {
	world.GameOver = true
	pop := NewEndGameMenu(world, 5, 5, 50, 6, Red, "GAME OVER", "You let the food spoil!")
	m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: &pop}})
}

func (world *World) ShowGameWonMenu() {
	world.GameOver = true
	pop := NewEndGameMenu(world, 5, 5, 50, 6, LightBlue, "You won the game!", fmt.Sprintf("Delivered with %v heat remaning", world.Player.HT.Current))
	m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: &pop}})
}

func (world *World) Notify(message m.M) {
	switch message.ID {
	case MoveEntity:
		if d, ok := message.Data.(MoveEntityMessage); ok {
			world.MoveEntity(d)
		}
	case KillEntity:
		if d, ok := message.Data.(KillEntityMessage); ok {
			world.RemoveEntity(d.Defender)
		}
	case PlayerDead:
		world.ShowEndGameMenu()
	case GameWon:
		world.ShowGameWonMenu()
	case FoodSpoiled:
		world.ShowFoodSpoiledMenu()
	case PlayerFloorChange:
		if d, ok := message.Data.(PlayerFloorChangeMessage); ok {
			if !d.Connected {
				break
			}
			world.RemoveEntity(world.Player)
			world.Player.X = d.DestX
			world.Player.Y = d.DestY
			world.LevelChanged = true
			world.SetCurrentLevel(d.DestLevelID)
			world.AddEntityToCurrentLevel(world.Player)
		}
	case TryMoveCreature:
		if d, ok := message.Data.(TryMoveCreatureMessage); ok {
			d.Creature.TryMove(d.X, d.Y, world)
		}
	case ShowMenu:
		if d, ok := message.Data.(ShowMenuMessage); ok {
			log.Printf("World: %T %+v", d.Menu, d.Menu)
			world.MenuStack = append(world.MenuStack, d.Menu)
		}
	}
}

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

		level := GenLevel(world.rng, 72, 72, genFlags)
		world.AddLevelFromCandidate(level)
	}
	world.SetCurrentLevel(0)
}

// LevelById returns a pointer to a specific level. Do not save this pointer!
func (w *World) LevelByID(id int) *Level {
	if id == w.CurrentLevelID {
		return &w.Levels[w.CurrentLevelIndex]
	}
	for i := range w.Levels {
		if w.Levels[i].ID == id {
			return &w.Levels[i]
		}
	}
	return nil
}

func NewWorld(window *gterm.Window, centered bool, rng *pcg.PCG64) *World {
	world := &World{
		Window:         window,
		CameraCentered: centered,
		MaxDepth:       MaxDepth,
		CameraX:        0,
		CameraY:        0,
		// TODO: Width/Height should probably be some function of the window dimensions
		CameraWidth:        56,
		CameraHeight:       25,
		CurrentUpdateTicks: sdl.GetTicks(),
		rng:                rng,
	}

	m.Subscribe(world.Notify)

	world.GameLog = NewGameLog(0, window.Rows-4, window.Columns, 3, world)

	world.BuildLevels()

	return world
}
