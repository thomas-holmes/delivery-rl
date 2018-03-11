package main

import (
	"log"

	"github.com/thomas-holmes/delivery-rl/game/dice"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
)

type CombatSystem struct {
	World *World
}

func NewCombatSystem(world *World) *CombatSystem {
	c := &CombatSystem{World: world}
	m.Subscribe(c.Notify)
	return c
}

func (combat CombatSystem) fight(a Entity, d Entity) {
	attacker, ok := a.(*Creature)
	if !ok {
		log.Panicf("Got a non-creature %+v", a)
		return
	}
	defender, ok := d.(*Creature)
	if !ok {
		log.Panicf("Got a non-creature %+v", d)
		return
	}

	log.Printf("Fighting with %+v", attacker.Equipment)
	damage := dice.Roll(attacker.Equipment.Weapon.Power) + (attacker.Level / 3)
	reduction := max(0, dice.Roll(dice.Notation{Num: 1, Sides: defender.Level})-1)
	actual := max(0, damage-reduction)
	defender.Damage(actual)
	gl.Append("%v hits %v for %v damage but is reduced by (%v)!", attacker.Name, defender.Name, actual, reduction)

	// This should be done by the entity instead of here, I think?
	// I think this used to attribute the experience gain on death. Maybe
	// need a more sophisticated combat/exp tracking system instead based
	// on damage dealt & proximity?
	if defender.HP.Current == 0 {
		m.Broadcast(m.M{ID: KillEntity, Data: KillEntityMessage{Attacker: a, Defender: d}})
	}
}

func (combat CombatSystem) zap(a Entity, d Entity, s Spell) {
	attacker, ok := a.(*Creature)
	if !ok {
		log.Panicf("Got a non-creature %+v", a)
		return
	}
	defender, ok := d.(*Creature)
	if !ok {
		log.Panicf("Got a non-creature %+v", d)
		return
	}

	log.Printf("Spell attacking with %+v", s)
	for i := 0; i < s.Hits; i++ {
		defender.Damage(s.Power)
		gl.Append("%v hits %v with %v for %v damage!", attacker.Name, defender.Name, s.Name, s.Power)
	}

	// This should be done by the entity instead of here, I think?
	// I think this used to attribute the experience gain on death. Maybe
	// need a more sophisticated combat/exp tracking system instead based
	// on damage dealt & proximity?
	if defender.HP.Current == 0 {
		m.Broadcast(m.M{ID: KillEntity, Data: KillEntityMessage{Attacker: a, Defender: d}})
	}

}

func (combat CombatSystem) zapCone(launch SpellLaunchMessage) {
	spell := launch.Spell
	player := combat.World.Player
	for _, pos := range conePositions(player.X, player.Y, launch.X, launch.Y, spell.Size) {
		if pos.X < 0 || pos.Y < 0 || pos.X >= combat.World.CurrentLevel().Columns || pos.Y >= combat.World.CurrentLevel().Rows {
			continue
		}
		if c, ok := combat.World.CurrentLevel().GetCreatureAtTile(pos.X, pos.Y); ok {
			combat.zap(launch.Caster, c, spell)
		}
	}
}

func (combat CombatSystem) zapSquare(launch SpellLaunchMessage) {
	spell := launch.Spell
	minX := max(launch.X-spell.Size, 0)
	maxX := min(launch.X+spell.Size, combat.World.CurrentLevel().Columns)

	minY := max(launch.Y-spell.Size, 0)
	maxY := min(launch.Y+spell.Size, combat.World.CurrentLevel().Rows)

	for y := minY; y < maxY+1; y++ {
		for x := minX; x < maxX+1; x++ {
			if c, ok := combat.World.CurrentLevel().GetCreatureAtTile(x, y); ok {
				combat.zap(launch.Caster, c, spell)
			}
		}
	}
}

func (combat CombatSystem) MovementSpell(launch SpellLaunchMessage) {
	if c, ok := launch.Caster.(*Creature); ok {
		c.TryTeleport(launch.X, launch.Y, combat.World)
	}
}

func (combat CombatSystem) resolveSpell(launch SpellLaunchMessage) {
	if launch.Spell.Name == "Fire Bolt" {
		if c, ok := launch.Caster.(*Creature); ok {
			a := NewLinearSpellAnimation(c.X, c.Y, launch.X, launch.Y, 20, 0, '*', Red)
			combat.World.AddAnimation(&a)
		}
	} else if launch.Spell.Name == "Magic Missile" {
		if c, ok := launch.Caster.(*Creature); ok {
			for i := 0; i < launch.Spell.Hits; i++ {
				a := NewLinearSpellAnimation(c.X, c.Y, launch.X, launch.Y, 20, uint32(i*40), '*', LightBlue)
				combat.World.AddAnimation(&a)
			}
		}

	}
	if launch.Spell.Hits == 0 {
		combat.MovementSpell(launch)
	} else {
		switch launch.Spell.Shape {
		case Square:
			combat.zapSquare(launch)
		case Line:
			// Nothing yet
		case Cone:
			combat.zapCone(launch)
		}
	}
}

func (combat CombatSystem) Notify(message m.M) {
	switch message.ID {
	case AttackEntity:
		if d, ok := message.Data.(AttackEntityMesasge); ok {
			combat.fight(d.Attacker, d.Defender)
		}
	case SpellLaunch:
		if d, ok := message.Data.(SpellLaunchMessage); ok {
			combat.resolveSpell(d)
		}
	}
}
