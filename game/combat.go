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
	if defender.Identity() == attacker.Identity() {
		gl.Append("%s flails around and hurts itself!!", attacker.Name)
	}
	gl.Append("%v hits %v for %v damage but is reduced by (%v)!", attacker.Name, defender.Name, actual, reduction)

	// This should be done by the entity instead of here, I think?
	// I think this used to attribute the experience gain on death. Maybe
	// need a more sophisticated combat/exp tracking system instead based
	// on damage dealt & proximity?
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
