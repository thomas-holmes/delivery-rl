package game

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func (player *Player) UpdatePosition(xPos int, yPos int, world *World) {
	if xPos >= 0 && xPos < world.Columns &&
		yPos >= 0 && yPos < world.Rows {
		if world.IsTileMonster(xPos, yPos) {
			monster := world.GetMonsterAtTile(xPos, yPos)
			player.Broadcast(PlayerAttack, PlayerAttackMessage{
				Player:  player,
				Monster: monster,
			})
		} else if world.CanStandOnTile(xPos, yPos) {
			oldX := player.xPos
			oldY := player.yPos
			player.xPos = xPos
			player.yPos = yPos
			player.Broadcast(MoveEntity, MoveEntityMessage{ID: player.ID(), OldX: oldX, OldY: oldY, NewX: xPos, NewY: yPos})
			player.Broadcast(PlayerUpdate, nil)
		}
	}
}

func (player *Player) Render(world *World) {
	world.Window.AddToCell(player.xPos, player.yPos, player.RenderGlyph, player.RenderColor)
}

type Health struct {
	Current int
	Max     int
}

// Player pepresents the player
type Player struct {
	id          int
	HP          Health
	Level       int
	Experience  int
	Name        string
	xPos        int
	yPos        int
	RenderGlyph string
	RenderColor sdl.Color

	Messaging
}

func (player *Player) GainExp(exp int) {
	player.Experience += exp
	log.Println("Got some exp", exp)
	if player.Experience >= player.Level {
		player.Experience -= player.Level
		player.Level++
		player.Broadcast(PlayerUpdate, nil)
	}
}

func (player *Player) XPos() int {
	return player.xPos
}

func (player *Player) YPos() int {
	return player.yPos
}

func (player Player) ID() int {
	return player.id
}

func (player *Player) SetID(id int) {
	player.id = id
}

func NewPlayer(xPos int, yPos int) Player {
	player := Player{
		HP:          Health{Current: 5, Max: 5},
		Level:       1,
		RenderGlyph: "@",
		RenderColor: sdl.Color{R: 255, G: 0, B: 0, A: 0},
		xPos:        xPos,
		yPos:        yPos,
	}

	log.Printf("Made a player, %#v", player)
	return player
}

func (player Player) HealthPercentage() float32 {
	current := float32(player.HP.Current)
	max := float32(player.HP.Max)
	return current / max
}

func (player *Player) Damage(amount int) {
	amount = max(amount, 0)

	newHp := max(player.HP.Current-amount, 0)
	player.HP.Current = newHp

	player.Broadcast(PlayerUpdate, nil)
	if newHp == 0 {
		player.Broadcast(PlayerDead, nil)
	}

}

func (player *Player) Heal(amount int) {
	amount = max(amount, 0)

	newHp := min(player.HP.Current+amount, player.HP.Max)
	player.HP.Current = newHp

	player.Broadcast(PlayerUpdate, nil)
}

// HandleInput updates player position based on user input
func (player *Player) HandleInput(event sdl.Event, world *World) {
	switch e := event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_h:
			player.UpdatePosition(player.xPos-1, player.yPos, world)
		case sdl.K_j:
			player.UpdatePosition(player.xPos, player.yPos+1, world)
		case sdl.K_k:
			player.UpdatePosition(player.xPos, player.yPos-1, world)
		case sdl.K_l:
			player.UpdatePosition(player.xPos+1, player.yPos, world)
		case sdl.K_b:
			player.UpdatePosition(player.xPos-1, player.yPos+1, world)
		case sdl.K_n:
			player.UpdatePosition(player.xPos+1, player.yPos+1, world)
		case sdl.K_y:
			player.UpdatePosition(player.xPos-1, player.yPos-1, world)
		case sdl.K_u:
			player.UpdatePosition(player.xPos+1, player.yPos-1, world)
		case sdl.K_1:
			player.Damage(1)
		case sdl.K_2:
			player.Heal(1)
		}
	}
}

// SetColor updates the render color of the player
func (player *Player) SetColor(color sdl.Color) {
	player.RenderColor = color
}