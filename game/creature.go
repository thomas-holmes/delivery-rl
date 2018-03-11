package main

import (
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/thomas-holmes/delivery-rl/game/dice"

	"github.com/thomas-holmes/delivery-rl/game/items"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Team int

const (
	NeutralTeam Team = iota
	PlayerTeam
	MonsterTeam
)

func ParseTeam(t string) Team {
	switch strings.ToLower(t) {
	case "neutral":
		return NeutralTeam
	case "player":
		return PlayerTeam
	case "monster":
		return MonsterTeam
	}

	return NeutralTeam
}

func (creature *Creature) Regen() {
	creature.HP.Tick()
	creature.ST.Tick()

	creature.HT.Tick()
}

type Creature struct {
	BasicEntity

	CompletedExternalAction bool

	IsPlayer       bool
	IsDragon       bool
	VisionDistance int

	Experience int

	RenderGlyph rune
	RenderColor sdl.Color

	Depth int

	Team

	State MonsterBehavior

	X int
	Y int

	CurrentlyActing bool
	GainedEnergy    bool
	Energy

	Effects map[StatusEffect]int

	Inventory

	Equipment

	BaseSpeed     int
	SpeedModifier int

	HP Resource // health
	ST Resource // stamina

	HT Resource // heat

	Level int

	Name string

	m.Unsubscribe
}

func (c *Creature) CheckTerrain(world *World) {
	tile := world.CurrentLevel().GetTile(c.X, c.Y)
	switch tile.TileEffect {
	case None:
		c.SpeedModifier = 1
	case Greasy:
		c.SpeedModifier = 4
	}
}

func (c *Creature) StartTurn() {
	if !c.GainedEnergy {
		c.Energy.AddEnergy(100)
		c.GainedEnergy = true

		if c.Energy.Current >= c.BaseSpeed {
			c.CurrentlyActing = true
		}
	}
}

func (c *Creature) EndTurn() {
	c.CurrentlyActing = false
	c.GainedEnergy = false
	c.Regen()
	c.TickEffects()
}

func (c Creature) IsDead() bool {
	return c.HP.Current <= 0
}

func (c Creature) CanAct() bool {
	return !c.IsDragon && !c.IsDead() && (c.CurrentlyActing && c.Energy.Current >= c.BaseSpeed)
}

func (c Creature) XPos() int {
	return c.X
}

func (c Creature) YPos() int {
	return c.Y
}

func (c Creature) Speed() int {
	return c.BaseSpeed * c.SpeedModifier
}

func (c *Creature) Damage(damage int) {
	c.HP.Current = max(0, c.HP.Current-damage)
}

func (c *Creature) TryMove(newX int, newY int, world *World) (MoveResult, interface{}) {
	confused := c.HasStatus(Confused)

	tile := world.CurrentLevel().GetTile(newX, newY)
	if !c.IsPlayer && tile.Item.Kind == items.Food {
		return MoveIsFood, nil
	}

	if world.CurrentLevel().CanStandOnTile(newX, newY) {
		return MoveIsSuccess, nil
	}

	if defender, ok := world.CurrentLevel().GetCreatureAtTile(newX, newY); ok {
		if c.IsPlayer && defender.IsDragon {
			return MoveIsVictory, nil
		}
		// Attack friendlies if confused
		if confused || (c.Team != NeutralTeam) && (c.Team != defender.Team) {
			// Check if I still need to get entity, I think this isn't necessary any more
			a, aOk := world.GetEntity(c.ID)
			d, dOk := world.GetEntity(defender.ID)
			if aOk && dOk {
				return MoveIsEnemy, MoveEnemy{Attacker: a, Defender: d}
			}
		}
	}

	return MoveIsInvalid, nil
}

func (player *Creature) TryWarp(world *World, newX, newY int) bool {
	if newX != player.X || newY != player.Y {
		result, _ := player.TryMove(newX, newY, world)
		switch result {
		case MoveIsInvalid:
			return false
		case MoveIsSuccess:
			oldX := player.X
			oldY := player.Y
			player.X = newX
			player.Y = newY
			player.ST.Current -= WarpCost
			m.Broadcast(m.M{ID: MoveEntity, Data: MoveEntityMessage{ID: player.ID, OldX: oldX, OldY: oldY, NewX: newX, NewY: newY}})
		case MoveIsEnemy:
			return false
		case MoveIsVictory:
			m.Broadcast(m.M{ID: GameWon})
		}
	}
	player.CompletedExternalAction = true
	return true
}

func NewCreature(level int, maxHP int) *Creature {
	return &Creature{
		Level: level,
		Team:  NeutralTeam,
		Energy: Energy{
			Current: 100,
			Max:     100,
		},
		HP:        Resource{Current: maxHP, Max: maxHP, RegenRate: 0.05},
		ST:        Resource{Current: 2, Max: 2, RegenRate: 0.15},
		Effects:   make(map[StatusEffect]int),
		BaseSpeed: 100,
		Equipment: NewEquipment(),
	}
}

func (p *Creature) AddStartingItems() {
	if def, ok := items.GetCollection("consumeables").GetByName("Hand Warmer"); ok {
		it := produceItem(def)
		it.Count = 5
		p.Inventory.Add(it)
	}

	if def, ok := items.GetCollection("consumeables").GetByName("Chicken Wing"); ok {
		it := produceItem(def)
		it.Count = 5
		p.Inventory.Add(it)
	}

	if def, ok := items.GetCollection("consumeables").GetByName("Garlic Butter"); ok {
		it := produceItem(def)
		it.Count = 5
		p.Inventory.Add(it)
	}

	if def, ok := items.GetCollection("consumeables").GetByName("Red Pepper Flakes"); ok {
		it := produceItem(def)
		it.Count = 5
		p.Inventory.Add(it)
	}

	if def, ok := items.GetCollection("consumeables").GetByName("Breadstick"); ok {
		it := produceItem(def)
		it.Count = 5
		p.Inventory.Add(it)
	}
}

func NewPlayer() *Creature {
	player := NewCreature(1, 20)
	player.Team = PlayerTeam
	player.RenderGlyph = '@'
	player.RenderColor = Red
	player.IsPlayer = true
	player.VisionDistance = 12
	player.HP.RegenRate = 0.15
	player.ST = Resource{Current: 4, Max: 4, RegenRate: 0.10}
	player.HT = Resource{Current: 125, Max: 125, RegenRate: -0.1}

	player.Unsubscribe = m.Subscribe(player.Notify)

	player.AddStartingItems()

	return player
}

func NewMonster(xPos int, yPos int, level int, hp int) *Creature {
	monster := NewCreature(level, hp)

	monster.X = xPos
	monster.Y = yPos
	monster.Team = MonsterTeam
	monster.RenderColor = Green
	monster.RenderGlyph = []rune(strconv.Itoa(monster.Level))[0]

	return monster
}

func (player *Creature) LevelUp() {
	player.Experience = max(0, player.Experience-player.NextLevelCost())
	player.Level++
	player.HP.Max = player.HP.Max + max(1, int(float64(player.HP.Max)*0.1))
	player.HP.Current = player.HP.Max
	player.ST.Max = player.ST.Max + max(1, int(float64(player.ST.Max)*0.1))
	player.ST.Current = player.ST.Max
	gl.Append("You are now level %d", player.Level)
}

const HasLevelUps bool = false

func (player *Creature) GainExp(exp int) {
	if HasLevelUps {
		player.Experience += exp
		if player.Experience >= player.NextLevelCost() {
			player.LevelUp()
		}
	}
}

func (player *Creature) NextLevelCost() int {
	return player.Level * 10
}

func (player *Creature) Heal(amount int) {
	amount = max(amount, 0)

	newHp := min(player.HP.Current+amount, player.HP.Max)
	player.HP.Current = newHp
}

func (player *Creature) BoostMaxHP(amount int) {
	player.HP.Max += amount
	player.HP.Current += amount
}

func (player *Creature) BoostMaxST(amount int) {
	player.ST.Max += amount
	player.ST.Current += amount
}

func (player *Creature) RestoreHeat(amount int) {
	amount = max(amount, 0)

	player.HT.Current = min(player.HT.Current+amount, player.HT.Max)
}

func (player *Creature) PickupItem(world *World) bool {
	tile := world.CurrentLevel().GetTile(player.X, player.Y)
	a := Item{}
	if tile.Item == a {
		return false
	}

	player.Inventory.Add(tile.Item)
	tile.Item = Item{}
	return true
}

func (creature *Creature) IsFoodRuined() bool {
	return creature.HT.Current <= 0
}

func (creature *Creature) EndGame() {
	m.Broadcast(m.M{ID: FoodSpoiled})
}

// Update returns true if an action that would constitute advancing the turn took place
func (creature *Creature) Update(turn uint64, input controls.InputEvent, world *World) bool {
	creature.CheckTerrain(world)

	success := false
	if creature.IsPlayer {
		if creature.IsFoodRuined() {
			creature.EndGame()
			return true
		}
		success = creature.HandleInput(input, world)
	} else {
		success = creature.Pursue(turn, world)
	}

	if success {
		creature.Energy.Current -= creature.Speed()
		return true
	}

	return false
}

func (creature *Creature) Quaff(potion Item) {
	if !potion.CanQuaff() {
		log.Fatalf("Asked to quaff unquaffable item. %+v", potion)
		return
	}

	QuaffPotion(creature, potion)
	creature.CompletedExternalAction = true
	creature.Inventory.RemoveItem(potion)
}

func (creature *Creature) ActivateItem(item Item) {
	if !item.CanActivate() {
		log.Fatalf("Asked to activate non-activateable item. %+v", item)
	}

	ActivateItem(creature, item)
	creature.CompletedExternalAction = true
	creature.Inventory.RemoveItem(item)
}

func (creature *Creature) DropItem(item Item, world *World) {
	if world.PlaceItemAround(item, creature.X, creature.Y) {
		creature.CompletedExternalAction = true
		gl.Append("Dropped %d %s", item.Count, item.Name)
		creature.Inventory.RemoveAllItem(item)
	} else {
		gl.Append("Could not drop %s, there was no room", item.Name)
	}
}
func (creature *Creature) ThrowItem(throwMessage PlayerThrowItemMessage) {
	if ThrowItem(creature, throwMessage.World, throwMessage.Item, throwMessage.TargetX, throwMessage.TargetY) {
		log.Printf("Throwing this thing %+v", throwMessage.Item)
		creature.CompletedExternalAction = true
	}
}

func (creature *Creature) TickEffects() {
	for k, v := range creature.Effects {
		if v-1 == 0 {
			delete(creature.Effects, k)
		} else {
			creature.Effects[k] = v - 1
		}
	}
}

func (creature *Creature) HasStatus(effect StatusEffect) bool {
	if _, ok := creature.Effects[effect]; ok {
		return ok
	}
	return false
}

func (creature *Creature) ApplyStatusEffect(effect StatusEffect) {
	if remaining, ok := creature.Effects[effect]; ok {
		creature.Effects[effect] = remaining + 5
	} else {
		creature.Effects[effect] = 5
	}
}

// HandleInput updates player position based on user input
func (player *Creature) HandleInput(input controls.InputEvent, world *World) bool {
	newX := player.X
	newY := player.Y

	if player.CompletedExternalAction {
		player.CompletedExternalAction = false
		return true
	}

	switch input.Action() {
	case controls.Ascend:
		tile := world.CurrentLevel().GetTile(player.X, player.Y)
		if tile.TileKind == UpStair {
			if stair, ok := world.CurrentLevel().getStair(player.X, player.Y); ok {
				m.Broadcast(m.M{ID: PlayerFloorChange, Data: PlayerFloorChangeMessage{
					Stair: stair,
				}})
				return true
			} else {
				return false
			}
		}
	case controls.Descend:
		tile := world.CurrentLevel().GetTile(player.X, player.Y)
		if tile.TileKind == DownStair {
			if stair, ok := world.CurrentLevel().getStair(player.X, player.Y); ok {
				m.Broadcast(m.M{ID: PlayerFloorChange, Data: PlayerFloorChangeMessage{
					Stair: stair,
				}})
				return true
			} else {
				return false
			}
		}
	case controls.Help:
		m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: NewHelpPop(2, 2, 50, world.Window.Rows-4)}})
		return false
	case controls.Left:
		newX = player.X - 1
	case controls.Down:
		newY = player.Y + 1
	case controls.Up:
		newY = player.Y - 1
	case controls.Right:
		newX = player.X + 1
	case controls.DownLeft:
		newX, newY = player.X-1, player.Y+1
	case controls.DownRight:
		newX, newY = player.X+1, player.Y+1
	case controls.UpLeft:
		newX, newY = player.X-1, player.Y-1
	case controls.UpRight:
		newX, newY = player.X+1, player.Y-1
	case controls.Wait:
		break
	case controls.Get:
		return player.PickupItem(world)
	case controls.Inventory:
		menu := &InventoryPop{World: world, PopMenu: PopMenu{X: 2, Y: 2, W: 30, H: world.Window.Rows - 4}, Inventory: player.Inventory}
		m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
		return false
	case controls.Examine:
		menu := &InspectionPop{PopMenu: PopMenu{X: 65, Y: 32, W: 34, H: 26}, World: world, InspectX: player.X, InspectY: player.Y}
		m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
		return false
	case controls.Warp:
		if player.ST.Current >= WarpCost {
			menu := NewWarpPop(world)
			m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
		} else {
			gl.Append("Costs %d ST to Warp.", player.ST.Current)
		}
		return false
	case controls.Messages:
		m.Broadcast(m.M{ID: ShowFullGameLog})
		return false
	case controls.Quit:
		world.QuitGame = true
		world.GameOver = true
		return false
	default:
		return false
	}

	if newX != player.X || newY != player.Y {
		result, data := player.TryMove(newX, newY, world)
		switch result {
		case MoveIsInvalid:
			return false
		case MoveIsSuccess:
			oldX := player.X
			oldY := player.Y
			player.X = newX
			player.Y = newY
			m.Broadcast(m.M{ID: MoveEntity, Data: MoveEntityMessage{ID: player.ID, OldX: oldX, OldY: oldY, NewX: newX, NewY: newY}})
		case MoveIsEnemy:
			if data, ok := data.(MoveEnemy); ok {
				m.Broadcast(m.M{ID: AttackEntity, Data: AttackEntityMesasge{
					Attacker: data.Attacker,
					Defender: data.Defender,
				}})
			}
		case MoveIsVictory:
			m.Broadcast(m.M{ID: GameWon})
			return false
		}
	}
	return true
}

func computeExperience(attacker *Creature, defender *Creature) int {
	axp := attacker.Level * attacker.Level
	dxp := defender.Level * defender.Level
	diff := max(defender.Level, dxp-axp)
	return diff
}

func (creature *Creature) Notify(message m.M) {
	switch message.ID {
	case KillEntity:
		if d, ok := message.Data.(KillEntityMessage); ok {
			attacker, ok := d.Attacker.(*Creature)
			if !ok {
				return
			}
			defender, ok := d.Defender.(*Creature)
			if !ok {
				return
			}

			if defender.ID == creature.ID {
				m.Broadcast(m.M{ID: PlayerDead})
				return
			}
			if attacker.ID != creature.ID {
				return
			}
			expGain := computeExperience(attacker, defender)
			gl.Append("Gained %d experience for killing %s", expGain, defender.Name)
			attacker.GainExp(computeExperience(attacker, defender))
		}
	case EquipItem:
		if d, ok := message.Data.(EquipItemMessage); ok {
			creature.CompletedExternalAction = true

			// Put it back in my inventory
			if creature.Equipment.Weapon.Name != "Bare Hands" {
				creature.Inventory.Add(creature.Equipment.Weapon)
			}

			creature.Equipment.Weapon = d.Item
			gl.Append("%s equips %s", creature.Name, d.Item.Name)
			creature.Inventory.RemoveItem(d.Item)
		}
	case PlayerQuaffPotion:
		if d, ok := message.Data.(PlayerQuaffPotionMessage); ok {
			creature.Quaff(d.Potion)
		}
	case PlayerActivateItem:
		if d, ok := message.Data.(PlayerActivateItemMessage); ok {
			creature.ActivateItem(d.Item)
		}
	case PlayerThrowItem:
		if d, ok := message.Data.(PlayerThrowItemMessage); ok {
			creature.ThrowItem(d)
		}
	case PlayerDropItem:
		if d, ok := message.Data.(PlayerDropItemMessage); ok {
			creature.DropItem(d.Item, d.World)
		}

	case PlayerWarp:
		if d, ok := message.Data.(PlayerWarpMessage); ok {
			creature.TryWarp(d.World, d.TargetX, d.TargetY)
		}
	}
}

// SetColor updates the render color of the player
func (player *Creature) SetColor(color sdl.Color) {
	player.RenderColor = color
}

func (creature *Creature) Render(world *World) {
	world.RenderRuneAt(creature.X, creature.Y, creature.RenderGlyph, creature.RenderColor, gterm.NoColor)
}

func (monster *Creature) Pursue(turn uint64, world *World) bool {
	if world.CurrentLevel().VisionMap.VisibilityAt(monster.X, monster.Y) == Visible {
		monster.State = Pursuing
	}

	if monster.State != Pursuing {
		return true
	}

	var candidates []TrackCandidate
	confused := monster.HasStatus(Confused)

	if confused {
		for _, pos := range world.CurrentLevel().ClampedXY(monster.X, monster.Y, 1) {
			candidates = append(candidates, TrackCandidate{Position: pos, Scent: 0})
		}
		shuffle(world.rng, len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	} else {
		var closeFood []TrackCandidate
		for _, pos := range world.CurrentLevel().ClampedXY(monster.X, monster.Y, 5) {
			if world.CurrentLevel().GetTile(pos.X, pos.Y).Item.Kind == items.Food {
				closeFood = append(closeFood, TrackCandidate{Position: pos, Scent: euclideanDistance(monster.X, monster.Y, pos.X, pos.Y)})
			}
		}

		var visibleFood []TrackCandidate
		for _, tc := range closeFood {
			path := PlotLine(monster.X, monster.Y, tc.X, tc.Y)
			for _, pos := range path {
				if world.CurrentLevel().GetTile(pos.X, pos.Y).IsWall() {
					break
				}
			}
			visibleFood = append(visibleFood, tc)
		}

		sort.Slice(visibleFood, func(i, j int) bool { return visibleFood[i].Scent < visibleFood[j].Scent })

		scent := world.CurrentLevel().ScentMap
		candidates = scent.track(turn, monster.X, monster.Y)
		if len(visibleFood) > 0 {
			closestFood := visibleFood[0]
			pathToFood := PlotLine(monster.X, monster.Y, closestFood.X, closestFood.Y)
			if len(pathToFood) > 1 {
				closestTile := pathToFood[1]
				foodCandidate := TrackCandidate{Position: closestTile, Scent: math.MaxFloat64}
				candidates = append(candidates, foodCandidate)
				sort.Slice(candidates, func(i, j int) bool { return candidates[i].Scent > candidates[j].Scent })
			}
		}
	}

	if len(candidates) > 0 {
		for _, choice := range candidates {
			result, data := monster.TryMove(choice.X, choice.Y, world)
			switch result {
			case MoveIsInvalid:
				continue
			case MoveIsFood:
				tile := world.CurrentLevel().GetTile(choice.X, choice.Y)
				if dice.Roll(tile.Item.Power) == 1 {
					tile.Item.Count--
					gl.Append("%s eats a whole %s!", monster.Name, tile.Item.Name)
				} else {

					gl.Append("%s takes a bite out of a %s!", monster.Name, tile.Item.Name)
				}
				if tile.Item.Count == 0 {
					tile.Item = Item{}
				}
			case MoveIsSuccess:
				oldX := monster.X
				oldY := monster.Y
				monster.X = choice.X
				monster.Y = choice.Y
				m.Broadcast(m.M{ID: MoveEntity, Data: MoveEntityMessage{
					ID:   monster.ID,
					OldX: oldX,
					OldY: oldY,
					NewX: choice.X,
					NewY: choice.Y,
				}})
			case MoveIsEnemy:
				if data, ok := data.(MoveEnemy); ok {
					m.Broadcast(m.M{ID: AttackEntity, Data: AttackEntityMesasge{
						Attacker: data.Attacker,
						Defender: data.Defender,
					}})
				}
			}
			return true
		}
	} else {
		return true
	}
	return false
}

type MoveResult int

const (
	MoveIsInvalid MoveResult = iota
	MoveIsSuccess
	MoveIsFood
	MoveIsEnemy
	MoveIsVictory
)

type MoveEnemy struct {
	Attacker Entity
	Defender Entity
}

type StatusEffect int

const (
	Confused StatusEffect = iota
)

type MonsterBehavior int

const (
	Idle MonsterBehavior = iota
	Pursuing
)
