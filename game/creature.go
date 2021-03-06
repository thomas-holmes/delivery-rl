package main

import (
	"flag"
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
	// Probably add ID? Do I even need it?

	CompletedExternalAction bool

	IsPlayer       bool
	IsDragon       bool
	VisionDistance int

	RenderGlyph rune
	RenderColor sdl.Color

	Depth int

	Team

	State MonsterBehavior

	Resting bool

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
		c.RemoveStatusEffect(Slow)
	case Greasy:
		c.ApplyStatusEffect(Slow, 10, false)
	}
}

func (c *Creature) StartTurn(world *World) {
	c.CheckTerrain(world)
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
	if c.HasStatus(Slow) {
		c.SpeedModifier = 4
	} else {
		c.SpeedModifier = 1
	}
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
			return MoveIsEnemy, MoveEnemy{Attacker: c, Defender: defender}
		}
	}

	return MoveIsInvalid, nil
}

func (player *Creature) TryWarp(world *World, newX, newY, cost int) bool {
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
			player.ST.Current -= cost
			m.Broadcast(m.M{ID: MoveCreature, Data: MoveCreatureMessage{Creature: player, OldX: oldX, OldY: oldY, NewX: newX, NewY: newY}})
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
		HP:            Resource{Current: maxHP, Max: maxHP, RateTimes100: 5},
		ST:            Resource{Current: 2, Max: 2, RateTimes100: 15},
		Effects:       make(map[StatusEffect]int),
		BaseSpeed:     100,
		SpeedModifier: 1,
		Equipment:     NewEquipment(),
	}
}

func (p *Creature) AddStartingItems() {
	consumeables := items.GetCollection("consumeables")

	if def, ok := consumeables.GetByName("Hand Warmer"); ok && InitialWarmers > 0 {
		it := produceItem(def)
		it.Count = InitialWarmers
		p.Inventory.Add(it)
	}

	if def, ok := consumeables.GetByName("Chicken Wing"); ok && InitialWings > 0 {
		it := produceItem(def)
		it.Count = InitialWings
		p.Inventory.Add(it)
	}

	if def, ok := consumeables.GetByName("Garlic Butter"); ok && InitialButter > 0 {
		it := produceItem(def)
		it.Count = InitialButter
		p.Inventory.Add(it)
	}

	if def, ok := consumeables.GetByName("Red Pepper Flakes"); ok && InitialPepper > 0 {
		it := produceItem(def)
		it.Count = InitialPepper
		p.Inventory.Add(it)
	}

	if def, ok := consumeables.GetByName("Breadstick"); ok && InitialBread > 0 {
		it := produceItem(def)
		it.Count = InitialBread
		p.Inventory.Add(it)
	}

	p.Equipment.Armour = TShirt
}

func NewPlayer() *Creature {
	player := NewCreature(1, 30)
	player.Team = PlayerTeam
	player.RenderGlyph = '@'
	player.RenderColor = Red
	player.IsPlayer = true
	player.VisionDistance = 12
	player.HP.RateTimes100 = int(HPRegen * 100)
	player.ST = Resource{Current: 4, Max: 4, RateTimes100: int(STRegen * 100)}
	player.HT = Resource{Current: 125, Max: 125, RateTimes100: int(-HeatDecay * 100)}

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

	if player.Inventory.Add(tile.Item) {
		tile.Item = Item{}
		return true
	}
	return false
}

func (creature *Creature) IsFoodRuined() bool {
	return creature.HT.Current <= 0
}

func (creature *Creature) EndGame() {
	m.Broadcast(m.M{ID: FoodSpoiled})
}

func (player *Creature) Safe(world *World) bool {
	return !world.CurrentLevel().CreatureInSight()
}

// Update returns true if an action that would constitute advancing the turn took place
func (creature *Creature) Update(turn uint64, action controls.Action, world *World) bool {
	success := false

	if creature.IsPlayer {
		if creature.IsFoodRuined() {
			creature.EndGame()
			return true
		} else if creature.Resting {
			if creature.HP.Current >= creature.HP.Max {
				creature.Resting = false
			} else if !creature.Safe(world) {
				creature.Resting = false
			} else {
				success = true
			}
		} else {
			success = creature.HandleInput(action, world)
		}
	} else {
		success = creature.Pursue(turn, world)
	}

	creature.CheckTerrain(world)

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

func (creature *Creature) RemoveStatusEffect(effect StatusEffect) {
	delete(creature.Effects, effect)
}

func (creature *Creature) ApplyStatusEffect(effect StatusEffect, count int, stacks bool) {
	if remaining, ok := creature.Effects[effect]; ok {
		if stacks {
			creature.Effects[effect] = remaining + count
			log.Printf("Stacking Effect, now has %d", remaining+count)
		}
	} else {
		creature.Effects[effect] = count
	}
}

// HandleInput updates player position based on user input
func (player *Creature) HandleInput(action controls.Action, world *World) bool {
	newX := player.X
	newY := player.Y

	if player.CompletedExternalAction {
		player.CompletedExternalAction = false
		return true
	}

	switch action {
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
		menu := &InventoryPop{World: world, PopMenu: PopMenu{X: 2, Y: 2, W: 30, H: 30}, Inventory: player.Inventory}
		m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
		return false
	case controls.Warp:
		if player.ST.Current >= 1 {
			menu := NewWarpPop(world)
			m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
		} else {
			gl.Append("Costs at least 1 ST to Warp.")
		}
		return false
	case controls.Rest:
		if player.Safe(world) {
			player.Resting = true
		}
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
			m.Broadcast(m.M{ID: MoveCreature, Data: MoveCreatureMessage{Creature: player, OldX: oldX, OldY: oldY, NewX: newX, NewY: newY}})
		case MoveIsEnemy:
			if data, ok := data.(MoveEnemy); ok {
				m.Broadcast(m.M{ID: AttackCreature, Data: AttackCreatureMessage{
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

func (player *Creature) Notify(message m.M) {
	switch message.ID {
	case KillCreature:
		if d, ok := message.Data.(KillCreatureMessage); ok {
			attacker, defender := d.Attacker, d.Defender

			if defender == player {
				m.Broadcast(m.M{ID: PlayerDead})
				return
			}
			if attacker != player {
				return
			}
		}
	case EquipItem:
		if d, ok := message.Data.(EquipItemMessage); ok {
			player.CompletedExternalAction = true

			switch d.Item.Kind {
			case items.Weapon:
				// Put it back in my inventory
				player.Inventory.RemoveItem(d.Item)

				if player.Equipment.Weapon.Name != "Bare Hands" {
					if ok := player.Inventory.Add(player.Equipment.Weapon); !ok {
						// Shouldn't happen
						log.Printf("Failed to add a weapon to inventory when equipping. There should have been space.")
					}
				}

				player.Equipment.Weapon = d.Item
				gl.Append("%s equips %s", player.Name, d.Item.Name)
			case items.Armour:
				player.Inventory.RemoveItem(d.Item)
				if ok := player.Inventory.Add(player.Equipment.Armour); !ok {
					// Shouldn't happen
					log.Printf("Failed to add an armour to inventory when equipping. There should have been space.")
				}
				player.Equipment.Armour = d.Item
				gl.Append("%s equips %s", player.Name, d.Item.Name)
			}
		}
	case PlayerQuaffPotion:
		if d, ok := message.Data.(PlayerQuaffPotionMessage); ok {
			player.Quaff(d.Potion)
		}
	case PlayerActivateItem:
		if d, ok := message.Data.(PlayerActivateItemMessage); ok {
			player.ActivateItem(d.Item)
		}
	case PlayerThrowItem:
		if d, ok := message.Data.(PlayerThrowItemMessage); ok {
			player.ThrowItem(d)
		}
	case PlayerDropItem:
		if d, ok := message.Data.(PlayerDropItemMessage); ok {
			player.DropItem(d.Item, d.World)
		}
	case PlayerWarp:
		if d, ok := message.Data.(PlayerWarpMessage); ok {
			player.TryWarp(d.World, d.TargetX, d.TargetY, d.Cost)
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

func (monster *Creature) FindFood(world *World) []TrackCandidate {

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

	return visibleFood
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
		visibleFood := monster.FindFood(world)

		sort.Slice(visibleFood, func(i, j int) bool { return visibleFood[i].Scent < visibleFood[j].Scent })

		scent := world.CurrentLevel().ScentMap
		candidates = scent.track(turn, monster.X, monster.Y)
		if len(visibleFood) > 0 {
			closestFood := visibleFood[0]

			pathToFood := PlotLine(monster.X, monster.Y, closestFood.X, closestFood.Y)
			if len(pathToFood) >= 0 {
				var closestTile Position
				if len(pathToFood) > 1 {
					closestTile = pathToFood[1]
				} else {
					closestTile = pathToFood[0]

				}
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
				m.Broadcast(m.M{ID: MoveCreature, Data: MoveCreatureMessage{
					Creature: monster,
					OldX:     oldX,
					OldY:     oldY,
					NewX:     choice.X,
					NewY:     choice.Y,
				}})
			case MoveIsEnemy:
				if data, ok := data.(MoveEnemy); ok {
					m.Broadcast(m.M{ID: AttackCreature, Data: AttackCreatureMessage{
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
	Attacker *Creature
	Defender *Creature
}

type StatusEffect int

const (
	Confused StatusEffect = iota
	Slow
)

type MonsterBehavior int

const (
	Idle MonsterBehavior = iota
	Pursuing
)

var InitialWarmers int
var InitialWings int
var InitialBread int
var InitialPepper int
var InitialButter int

var HeatDecay float64
var HPRegen float64
var STRegen float64

func init() {
	flag.IntVar(&InitialWarmers, "starting-warmers", 0, "Starting amount of Hand Warmers.")
	flag.IntVar(&InitialWings, "starting-wings", 3, "Starting amount of Chicken Wings.")
	flag.IntVar(&InitialBread, "starting-bread", 3, "Starting amount of Breadsticks.")
	flag.IntVar(&InitialPepper, "starting-pepper", 3, "Starting amount of Red Pepper Flakes.")
	flag.IntVar(&InitialButter, "starting-Butter", 3, "Starting amount of Garlic Butter.")

	flag.Float64Var(&HeatDecay, "heat-decay", 0.10, "Amount of heat lost per turn.")
	flag.Float64Var(&HPRegen, "hp-regen", 0.15, "HP regen per turn.")
	flag.Float64Var(&STRegen, "st-regen", 0.10, "ST regen per turn.")
}
