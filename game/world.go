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

	nextLevelID int

	showScentOverlay bool

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
	if def.Armour.Name != "" {
		monster.Equipment.Armour = produceItem(items.Definition(def.Armour))
	}

	return monster
}

func (world *World) addInitialMonsters(level *Level) {
	monsterCollection := monsters.GetCollection("monsters")

	for tries := 0; tries < level.MonsterDensity+level.Depth; tries++ {
		x := int(world.rng.Bounded(uint64(level.Columns)))
		y := int(world.rng.Bounded(uint64(level.Rows)))

		if level.CanStandOnTile(x, y) {
			def := monsters.GetPowerBoundedMonster(world.rng, monsterCollection, level.Depth)
			monster := buildMonsterFromDefinition(def)
			monster.X = x
			monster.Y = y
			world.AddCreature(monster, level.ID)
		}
	}
}

func (world *World) addDragonToLevel(level *Level) {
	dragon := buildMonsterFromDefinition(monsters.Dragon)
	dragon.Team = NeutralTeam
	dragon.IsDragon = true

	world.AddCreature(dragon, level.ID)
}

// AddLevelFromCandidate constructs a real level from an intermediate level representation
func (world *World) AddLevelFromCandidate(level *CandidateLevel) {
	loadedLevel := LoadCandidateLevel(world.nextLevelID, level)
	world.nextLevelID++

	world.Levels = append(world.Levels, loadedLevel)

	levels := len(world.Levels)
	levelIndex := levels - 1

	if levels > 1 {
		connectTwoLevels(&world.Levels[levelIndex-1], &world.Levels[levelIndex])
	}
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
		world.addPlayer(creature, level)
	}
}

func (world *World) AddCreatureToCurrentLevel(c *Creature) {
	world.AddCreature(c, world.CurrentLevelID)
}

func (world *World) AddCreature(c *Creature, levelID int) {
	level := world.LevelByID(levelID)

	world.addCreature(c, level)
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

func (world *World) Update(action controls.Action) {
	currentTicks := sdl.GetTicks()
	world.CurrentTickDelta = currentTicks - world.CurrentUpdateTicks
	world.CurrentUpdateTicks = currentTicks

	// Update Menus
	{
		if world.tidyMenus() {
			currentMenu := world.MenuStack[len(world.MenuStack)-1]
			currentMenu.Update(action)
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

	// Update All Actors?
	{
		for {
			if world.GameOver {
				break
			}
			creatures := world.CurrentLevel().Creatures
			if world.CurrentLevel().NextCreature >= len(creatures) {
				world.CurrentLevel().NextCreature = 0
			}
			currentCreature := world.CurrentLevel().NextCreature
			creature := creatures[currentCreature]

			creature.StartTurn(world)

			if !creature.CanAct() {
				creature.EndTurn()
				world.CurrentLevel().NextCreature++
				continue
			}

			acted := creature.Update(world.turnCount, action, world)

			// bit of a hack, copied from below instead of rewritten
			if world.LevelChanged {
				creature.EndTurn()
				world.LevelChanged = false
				world.CurrentLevel().NextCreature = 0
				return // Is this the right thing to do? Or could we just break?
			}

			if creature.IsPlayer {
				if acted {
					world.turnCount++
				}
				world.CurrentLevel().VisionMap.UpdateVision(world.Player.VisionDistance, world)
				world.CurrentLevel().ScentMap.UpdateScents(world)
			}

			if creature.CanAct() {
				// Needs more input
				return
			} else {
				creature.EndTurn()
				world.CurrentLevel().NextCreature++
				action = controls.None
			}
		}
	}
}

func (world *World) UpdateCamera() {
	if world.CameraCentered {
		world.CameraOffsetX = 30
		world.CameraOffsetY = 25
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
	world.UpdateAnimations()
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

func (world *World) RemoveCreature(creature *Creature) {
	foundIndex := -1
	for i, c := range world.CurrentLevel().Creatures {
		if c == creature {
			foundIndex = i
			break
		}
	}

	if foundIndex > -1 {
		world.CurrentLevel().Creatures = append(world.CurrentLevel().Creatures[:foundIndex], world.CurrentLevel().Creatures[foundIndex+1:]...)
		if foundIndex > world.CurrentLevel().NextCreature {
			world.CurrentLevel().NextCreature--
		}
	}

	world.CurrentLevel().GetTile(creature.X, creature.Y).Creature = nil
}

func (world *World) MoveCreature(message MoveCreatureMessage) {
	oldTile := world.CurrentLevel().GetTile(message.OldX, message.OldY)
	newTile := world.CurrentLevel().GetTile(message.NewX, message.NewY)
	newTile.Creature = oldTile.Creature
	oldTile.Creature = nil
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
	pop := NewEndGameMenu(world, 5, 5, 50, 6, LightBlue, "You won the game!", fmt.Sprintf("Delivered in %d turns with %d heat remaning.", world.turnCount, world.Player.HT.Current))
	m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: &pop}})
}

func (world *World) Notify(message m.M) {
	switch message.ID {
	case MoveCreature:
		if d, ok := message.Data.(MoveCreatureMessage); ok {
			world.MoveCreature(d)
		}
	case KillCreature:
		if d, ok := message.Data.(KillCreatureMessage); ok {
			world.RemoveCreature(d.Defender)
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

			world.RemoveCreature(world.Player)
			world.Player.X = d.DestX
			world.Player.Y = d.DestY
			world.LevelChanged = true
			world.SetCurrentLevel(d.DestLevelID)
			world.AddCreatureToCurrentLevel(world.Player)
		}
	case TryMoveCreature:
		if d, ok := message.Data.(TryMoveCreatureMessage); ok {
			d.Creature.TryMove(d.X, d.Y, world)
		}
	case SplashGrease:
		if d, ok := message.Data.(SplashGreaseMessage); ok {
			world.SplashGrease(d.Item, d.X, d.Y)
		}
	case ShowMenu:
		if d, ok := message.Data.(ShowMenuMessage); ok {
			world.MenuStack = append(world.MenuStack, d.Menu)
		}
	}
}

func (world *World) PlaceItem(item Item, x, y int) bool {
	tile := world.CurrentLevel().GetTile(x, y)
	if tile.TileKind == Floor {
		if tile.Item.Kind == items.Unknown {
			tile.Item = item
			return true
		} else if item.Stacks && tile.Item.Name == item.Name {
			tile.Item.Count += item.Count
			return true
		}
	}

	return false
}

func (world *World) PlaceItemAround(item Item, x, y int) bool {
	if world.PlaceItem(item, x, y) {
		return true
	}
	// Try adjacent space
	cols, rows := world.CurrentLevel().Columns, world.CurrentLevel().Rows
	minX, maxX := max(0, x-1), min(cols-1, x+1)
	minY, maxY := max(0, y-1), min(rows-1, y+1)

	for iy := minY; iy <= maxY; iy++ {
		for ix := minX; ix <= maxX; ix++ {
			if world.PlaceItem(item, ix, iy) {
				return true
			}
		}
	}

	return false
}

func (world *World) SplashGrease(item Item, x, y int) {
	var greasedPositions []Position
	for _, pos := range world.CurrentLevel().ClampedXY(x, y, 1) {
		tile := world.CurrentLevel().GetTile(pos.X, pos.Y)
		if tile.TileKind == Floor {
			tile.TileEffect = Greasy
			greasedPositions = append(greasedPositions, pos)
		}
	}

	for {
		if len(greasedPositions) == 0 {
			break
		}

		greasedPos := greasedPositions[0]
		greasedPositions = greasedPositions[1:]

		for _, pos := range world.CurrentLevel().ClampedXY(greasedPos.X, greasedPos.Y, 1) {
			tile := world.CurrentLevel().GetTile(pos.X, pos.Y)
			if tile.TileKind == Floor && tile.TileEffect != Greasy {
				if world.rng.Bounded(8) == 0 {
					tile.TileEffect = Greasy
					greasedPositions = append(greasedPositions, pos)
				}
			}
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

		level := GenLevel(world.rng, 72, 72, i+1, genFlags)
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
		Window:             window,
		CameraCentered:     centered,
		MaxDepth:           MaxDepth,
		CameraX:            0,
		CameraY:            0,
		CameraWidth:        56,
		CameraHeight:       50,
		CurrentUpdateTicks: sdl.GetTicks(),
		rng:                rng,
	}

	m.Subscribe(world.Notify)

	world.GameLog = NewGameLog(0, window.Rows-9, window.Columns-35, 8, world)

	world.BuildLevels()

	return world
}
