package main

import (
	"log"
	"math/rand"
	"strconv"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Team int

const (
	PlayerTeam Team = iota
	MonsterTeam
)

type Health struct {
	Current int
	Max     int
}

type Creature struct {
	Identifiable

	IsPlayer bool

	Experience int

	RenderGlyph rune
	RenderColor sdl.Color

	Team

	State MonsterBehavior

	X int
	Y int

	Energy

	Inventory

	Equipment

	HP Health

	Level int

	Name string

	Messaging
}

func (c Creature) CanAct() bool {
	return c.currentEnergy >= 100
}

func (c Creature) XPos() int {
	return c.X
}

func (c Creature) YPos() int {
	return c.Y
}

func (c *Creature) Damage(damage int) {
	log.Printf("%v is Taking damage of %v", *c, damage)
	c.HP.Current = max(0, c.HP.Current-damage)

	if c.IsPlayer {
		c.Broadcast(PlayerUpdate, nil)
		if c.HP.Current == 0 {
			c.Broadcast(PlayerDead, nil)
		}
	}
}

// TODO: Delete once monster is a creature
func (c *Creature) Combatant() *Creature {
	return c
}

func (c *Creature) TryMove(newX int, newY int, world *World) (MoveResult, interface{}) {

	if world.CanStandOnTile(newX, newY) {
		return MoveIsSuccess, nil
	}

	if defender, ok := world.GetCreatureAtTile(newX, newY); ok {
		log.Printf("Got a creature in TryMove, %+v", *defender)
		if c.Team != defender.Team {
			a, aOk := world.GetEntity(c.ID)
			d, dOk := world.GetEntity(defender.ID)
			if aOk && dOk {
				return MoveIsEnemy, MoveEnemy{Attacker: a, Defender: d}
			}
		}
	}

	return MoveIsInvalid, nil
}

func NewPlayer() Creature {
	player := NewCreature(1, 100, 5)
	player.RenderGlyph = '@'
	player.RenderColor = Red
	player.IsPlayer = true
	player.Team = PlayerTeam

	log.Printf("Made a player, %#v", player)
	return player
}

func NewCreature(level int, maxEnergy int, maxHP int) Creature {
	return Creature{
		Level: level,
		Team:  MonsterTeam,
		Energy: Energy{
			currentEnergy: maxEnergy,
			maxEnergy:     maxEnergy,
		},
		HP:        Health{Current: 5, Max: 5},
		Equipment: NewEquipment(),
	}
}

func (player *Creature) LevelUp() {
	player.Experience -= max(0, player.Experience-player.Level)
	player.Level++
	player.HP.Max = int(float32(player.HP.Max) * 1.5)
	player.HP.Current = player.HP.Max
}

func (player *Creature) GainExp(exp int) {
	player.Experience += exp
	log.Println("Got some exp", exp)
	if player.Experience >= player.Level {
		player.LevelUp()
		player.Broadcast(PlayerUpdate, nil)
	}
}

func (player Creature) HealthPercentage() float32 {
	current := float32(player.HP.Current)
	max := float32(player.HP.Max)
	return current / max
}

func (player *Creature) Heal(amount int) {
	amount = max(amount, 0)

	newHp := min(player.HP.Current+amount, player.HP.Max)
	player.HP.Current = newHp

	player.Broadcast(PlayerUpdate, nil)
}

func (player *Creature) PickupItem(world *World) bool {
	tile := world.GetTile(player.X, player.Y)
	if tile.Item == nil {
		return false
	}

	player.Items = append(player.Items, tile.Item)
	tile.Item = nil
	return true
}

func (creature *Creature) Update(turn uint64, event sdl.Event, world *World) bool {
	// TODO: There should be a better way!
	if creature.IsPlayer {
		if creature.HandleInput(event, world) {
			creature.currentEnergy -= 100
			return true
		}
	} else {
		if creature.Pursue(turn, world) {
			creature.currentEnergy -= 100
			return true
		}
	}
	return false
}

// HandleInput updates player position based on user input
func (player *Creature) HandleInput(event sdl.Event, world *World) bool {
	newX := player.X
	newY := player.Y

	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_COMMA:
			tile := world.GetTile(player.X, player.Y)
			if tile.TileKind == UpStair {
				if stair, ok := world.CurrentLevel.getStair(player.X, player.Y); ok {
					player.Broadcast(PlayerFloorChange, PlayerFloorChangeMessage{
						Stair: stair,
					})
				} else {
					return false
				}
			}
			return true
		case sdl.K_PERIOD:
			tile := world.GetTile(player.X, player.Y)
			if tile.TileKind == DownStair {
				if stair, ok := world.CurrentLevel.getStair(player.X, player.Y); ok {
					player.Broadcast(PlayerFloorChange, PlayerFloorChangeMessage{
						Stair: stair,
					})
				} else {
					return false
				}
			}
			return true
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
		case sdl.K_p:
			return player.PickupItem(world)
		case sdl.K_i:
			menu := &InventoryPop{X: 10, Y: 2, W: 30, H: world.Window.Rows - 4, Inventory: player.Inventory}
			player.Broadcast(ShowMenu, ShowMenuMessage{Menu: menu})
			return false
		case sdl.K_e:
			menu := &EquipmentPop{X: 10, Y: 2, W: 30, H: world.Window.Rows - 4, Player: player}
			player.Broadcast(ShowMenu, ShowMenuMessage{Menu: menu})
			return false
		case sdl.K_ESCAPE:
			world.GameOver = true
			world.QuitGame = true
			return true
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
				player.Broadcast(MoveEntity, MoveEntityMessage{ID: player.ID, OldX: oldX, OldY: oldY, NewX: newX, NewY: newY})
			case MoveIsEnemy:
				if data, ok := data.(MoveEnemy); ok {
					player.Broadcast(AttackEntity, AttackEntityMesasge{
						Attacker: data.Attacker,
						Defender: data.Defender,
					})
				}
			}
		}
		return true
	}
	return false
}

func (creature *Creature) Notify(message Message, data interface{}) {
	if !creature.IsPlayer {
		return
	}
	switch message {
	case KillEntity:
		if d, ok := data.(KillEntityMessage); ok {
			aCombatant, ok := d.Attacker.(Combatant)
			if !ok {
				return
			}
			dCombatant, ok := d.Defender.(Combatant)
			if !ok {
				return
			}
			attacker, defender := aCombatant.Combatant(), dCombatant.Combatant()

			if defender.ID == creature.ID {
				creature.Broadcast(PlayerDead, nil)
				return
			}
			if attacker.ID != creature.ID {
				return
			}
			if creature.Level > defender.Level {
				creature.GainExp((defender.Level + 1) / 4)
			} else {
				creature.GainExp((defender.Level + 1) / 2)
			}
		}
	case EquipItem:
		if d, ok := data.(EquipItemMessage); ok {
			creature.Equipment.Weapon = d.Item // This is super low effort, but should work?
		}
	}
}

func (c *Creature) NeedsInput() bool {
	return c.IsPlayer
}

// SetColor updates the render color of the player
func (player *Creature) SetColor(color sdl.Color) {
	player.RenderColor = color
}

func (creature *Creature) Render(world *World) {
	world.RenderRuneAt(creature.X, creature.Y, creature.RenderGlyph, creature.RenderColor, gterm.NoColor)
}

// Monster Stuff
func NewMonster(xPos int, yPos int, level int, color sdl.Color, hp int) Creature {
	monster := NewCreature(level, 100, hp)

	monster.X = xPos
	monster.Y = yPos
	monster.Team = MonsterTeam
	monster.RenderColor = color
	monster.RenderGlyph = []rune(strconv.Itoa(monster.Level))[0]

	return monster
}

// TODO: If a monster is blocking the ideal path our monster should go around
func (monster *Creature) Pursue(turn uint64, world *World) bool {
	if world.CurrentLevel.VisionMap.VisibilityAt(monster.X, monster.Y) == Visible {
		monster.State = Pursuing
	}

	if monster.State != Pursuing {
		return true
	}

	scent := world.CurrentLevel.ScentMap

	// TODO: Maybe short circuit tracking here and just attack the player instead
	// if in ranger?
	candidates := scent.track(turn, monster.X, monster.Y)

	// TODO: Sometimes the monster takes a suboptimal path
	if len(candidates) > 0 {
		randomIndex := rand.Intn(len(candidates))
		choice := candidates[randomIndex]
		if len(candidates) > 1 {
			// TODO: Not actually sure if this is invalid but for now I want to know if it happens.
			log.Printf("More than one candidate, %+v", candidates)
		}

		result, data := monster.TryMove(choice.XPos, choice.YPos, world)
		log.Printf("Tried to move %#v, got result: %v, data %#v", monster, result, data)
		switch result {
		case MoveIsInvalid:
			log.Panicf("Monsters aren't allowed to yield their turn")
			return false
		case MoveIsSuccess:
			oldX := monster.X
			oldY := monster.Y
			monster.X = choice.XPos
			monster.Y = choice.YPos
			monster.Broadcast(MoveEntity, MoveEntityMessage{
				ID:   monster.ID,
				OldX: oldX,
				OldY: oldY,
				NewX: choice.XPos,
				NewY: choice.YPos,
			})
		case MoveIsEnemy:
			if data, ok := data.(MoveEnemy); ok {
				monster.Broadcast(AttackEntity, AttackEntityMesasge{
					Attacker: data.Attacker,
					Defender: data.Defender,
				})
			}
		}
		return true
	}
	return false
}

type MoveResult int

const (
	MoveIsInvalid MoveResult = iota
	MoveIsSuccess
	MoveIsEnemy
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
