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

	// Memory leak
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

	damage := dice.Roll(attacker.Equipment.Weapon.Power)
	reduction := dice.Roll(defender.Equipment.Armour.Power)
	actual := max(0, damage-reduction)
	defender.Damage(actual)
	if defender.Identity() == attacker.Identity() {
		gl.Append("%s flails around and hurts itself!!", attacker.Name)
	}
	gl.Append("%v hits %v for %v damage but is reduced by (%v)!", attacker.Name, defender.Name, actual, reduction)

	if defender.HP.Current == 0 {
		m.Broadcast(m.M{ID: KillEntity, Data: KillEntityMessage{Attacker: a, Defender: d}})
	}
}

func (combat CombatSystem) Notify(message m.M) {
	switch message.ID {
	case AttackEntity:
		if d, ok := message.Data.(AttackEntityMesasge); ok {
			combat.fight(d.Attacker, d.Defender)
		}
	}
}
