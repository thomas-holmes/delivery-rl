package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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

	Inventory

	Equipment

	Speed int

	HP Resource // health
	ST Resource // stamina

	HT Resource // heat

	Spells []Spell

	Level int

	Name string

	m.Unsubscribe
}

func (c *Creature) StartTurn() {
	if !c.GainedEnergy {
		c.Energy.AddEnergy(100)
		c.GainedEnergy = true
	}

	if c.Energy.Current >= 100 {
		c.CurrentlyActing = true
	}
}

func (c *Creature) EndTurn() {
	c.CurrentlyActing = false
	c.GainedEnergy = false
	c.Regen()
}

func (c Creature) IsDead() bool {
	return c.HP.Current <= 0
}

func (c Creature) CanAct() bool {
	return !c.IsDragon && !c.IsDead() && (c.Energy.Current >= 100 || (c.CurrentlyActing && c.Energy.Current >= c.Speed))
}

func (c Creature) XPos() int {
	return c.X
}

func (c Creature) YPos() int {
	return c.Y
}

func (c *Creature) Damage(damage int) {
	c.HP.Current = max(0, c.HP.Current-damage)
}

func (c *Creature) TryMove(newX int, newY int, world *World) (MoveResult, interface{}) {
	if world.CurrentLevel().CanStandOnTile(newX, newY) {
		return MoveIsSuccess, nil
	}

	if defender, ok := world.CurrentLevel().GetCreatureAtTile(newX, newY); ok {
		if c.IsPlayer && defender.IsDragon {
			return MoveIsVictory, nil
		}
		if (c.Team != NeutralTeam) && (c.Team != defender.Team) {
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

func (player *Creature) TryTeleport(newX int, newY int, world *World) bool {
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
		Speed:     100,
		Equipment: NewEquipment(),
	}
}

func NewPlayer() *Creature {
	player := NewCreature(1, 20)
	player.Team = PlayerTeam
	player.RenderGlyph = '@'
	player.RenderColor = Red
	player.IsPlayer = true
	player.Spells = DefaultSpells
	player.VisionDistance = 12
	player.HT = Resource{Current: 125, Max: 125, RegenRate: -0.2}

	player.Unsubscribe = m.Subscribe(player.Notify)

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

func (player *Creature) GainExp(exp int) {
	player.Experience += exp
	if player.Experience >= player.NextLevelCost() {
		player.LevelUp()
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
		creature.Energy.Current -= creature.Speed
		return true
	}

	return false
}

func (creature *Creature) TargetSpell(spell Spell, world *World) {
	menu := &SpellTargeting{PopMenu: PopMenu{X: 0, Y: 0, W: 0, H: 0}, TargetX: creature.X, TargetY: creature.Y, World: world, Spell: spell}
	m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
}

func (creature *Creature) CanCast(spell Spell) bool {
	if spell.Cost <= creature.ST.Current {
		return true
	}
	return false
}

func (creature *Creature) CastSpell(spell Spell, world *World, targetX int, targetY int) {
	fmt.Printf("Firing at (%v,%v) with %+v", targetX, targetY, spell)
	creature.CompletedExternalAction = true
	creature.ST.Current -= spell.Cost
	// Can attack self. Do we care?
	m.Broadcast(m.M{ID: SpellLaunch, Data: SpellLaunchMessage{Caster: creature, Spell: spell, X: targetX, Y: targetY}})
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

// HandleInput updates player position based on user input
func (player *Creature) HandleInput(input controls.InputEvent, world *World) bool {
	newX := player.X
	newY := player.Y

	if player.CompletedExternalAction {
		player.CompletedExternalAction = false
		return true
	}

	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_COMMA:
			if input.Keymod&sdl.KMOD_SHIFT > 0 {
				tile := world.CurrentLevel().GetTile(player.X, player.Y)
				if tile.TileKind == UpStair {
					if stair, ok := world.CurrentLevel().getStair(player.X, player.Y); ok {
						m.Broadcast(m.M{ID: PlayerFloorChange, Data: PlayerFloorChangeMessage{
							Stair: stair,
						}})
					} else {
						return false
					}
				}
			}
			return false
		case sdl.K_PERIOD:
			if input.Keymod&sdl.KMOD_SHIFT > 0 {
				tile := world.CurrentLevel().GetTile(player.X, player.Y)
				if tile.TileKind == DownStair {
					if stair, ok := world.CurrentLevel().getStair(player.X, player.Y); ok {
						m.Broadcast(m.M{ID: PlayerFloorChange, Data: PlayerFloorChangeMessage{
							Stair: stair,
						}})
					} else {
						return false
					}
				}
			}
			// Period returns true because it means "wait"
			return true
		case sdl.K_SLASH:
			if input.Keymod&sdl.KMOD_SHIFT > 0 {
				m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: NewHelpPop(2, 1, 50, world.Window.Rows-2)}})
				return false
			}
		case sdl.K_h:
			newX = player.X - 1
		case sdl.K_j:
			newY = player.Y + 1
		case sdl.K_k:
			newY = player.Y - 1
		case sdl.K_l:
			newX = player.X + 1
		case sdl.K_b:
			newX, newY = player.X-1, player.Y+1
		case sdl.K_n:
			newX, newY = player.X+1, player.Y+1
		case sdl.K_y:
			newX, newY = player.X-1, player.Y-1
		case sdl.K_u:
			newX, newY = player.X+1, player.Y-1
		case sdl.K_1:
			player.Damage(1)
			return false
		case sdl.K_2:
			player.Heal(1)
			return false
		case sdl.K_3:
			player.HT.Current++
			return false
		case sdl.K_g:
			return player.PickupItem(world)
		case sdl.K_i:
			menu := &InventoryPop{PopMenu: PopMenu{X: 6, Y: 2, W: 40, H: world.Window.Rows - 4}, Inventory: player.Inventory}
			m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
			return false
		case sdl.K_e:
			menu := &EquipmentPop{PopMenu: PopMenu{X: 6, Y: 2, W: 40, H: world.Window.Rows - 4}, Player: player}
			m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
			return false
		case sdl.K_x:
			menu := &InspectionPop{PopMenu: PopMenu{X: 60, Y: 20, W: 30, H: 5}, World: world, InspectX: player.X, InspectY: player.Y}
			m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
			return false
		case sdl.K_z:
			menu := &SpellPop{PopMenu: PopMenu{X: 10, Y: 2, W: 30, H: world.Window.Rows - 4}, World: world}
			m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: menu}})
			return false
		case sdl.K_m:
			m.Broadcast(m.M{ID: ShowFullGameLog})
			return false
		case sdl.K_q:
			if input.Keymod&sdl.KMOD_CTRL > 0 {
				world.QuitGame = true
				world.GameOver = true
			}
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
			}
		}
		return true
	}
	return false
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
			creature.Inventory.RemoveItem(d.Item)
		}
	case SpellTarget:
		if d, ok := message.Data.(SpellTargetMessage); ok {
			creature.TargetSpell(d.Spell, d.World)
		}
	case PlayerQuaffPotion:
		if d, ok := message.Data.(PlayerQuaffPotionMessage); ok {
			creature.Quaff(d.Potion)
		}
	case PlayerActivateItem:
		if d, ok := message.Data.(PlayerActivateItemMessage); ok {
			creature.ActivateItem(d.Item)
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

	scent := world.CurrentLevel().ScentMap

	// TODO: Maybe short circuit tracking here and just attack the player instead
	// if in range?
	candidates := scent.track(turn, monster.X, monster.Y)

	if len(candidates) > 0 {
		for _, choice := range candidates {
			result, data := monster.TryMove(choice.X, choice.Y, world)
			switch result {
			case MoveIsInvalid:
				continue
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
	MoveIsEnemy
	MoveIsVictory
)

type MoveEnemy struct {
	Attacker Entity
	Defender Entity
}

type MonsterBehavior int

const (
	Idle MonsterBehavior = iota
	Pursuing
)
