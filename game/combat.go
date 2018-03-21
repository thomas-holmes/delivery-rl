package main

import (
	"github.com/thomas-holmes/delivery-rl/game/dice"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
)

type CombatSystem struct {
	World *World
}

func NewCombatSystem(world *World) *CombatSystem {
	c := &CombatSystem{World: world}

	// Memory leak
	m.Subscribe(c.Notify)
	return c
}

func (combat CombatSystem) fight(attacker *Creature, defender *Creature) {
	damage := dice.Roll(attacker.Equipment.Weapon.Power)
	reduction := dice.Roll(defender.Equipment.Armour.Power)
	actual := max(0, damage-reduction)
	defender.Damage(actual)
	if defender == attacker {
		gl.Append("%s flails around and hurts itself!!", attacker.Name)
	}
	gl.Append("%v hits %v for %v damage but is reduced by (%v)!", attacker.Name, defender.Name, actual, reduction)

	if defender.HP.Current == 0 {
		m.Broadcast(m.M{ID: KillCreature, Data: KillCreatureMessage{Attacker: attacker, Defender: defender}})
	}
}

func (combat CombatSystem) Notify(message m.M) {
	switch message.ID {
	case AttackCreature:
		if d, ok := message.Data.(AttackCreatureMessage); ok {
			combat.fight(d.Attacker, d.Defender)
		}
	}
}
