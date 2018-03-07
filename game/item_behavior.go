package main

import (
	"github.com/thomas-holmes/delivery-rl/game/dice"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
)

func QuaffPotion(creature *Creature, potion Item) {
	power := dice.Roll(potion.Power)
	switch potion.Name {
	case "Health Potion":
		gl.Append("%s quaffs a %s and regains %d health!", creature.Name, potion.Name, power)
		creature.Heal(power)
	case "Power Potion":
		gl.Append("%s quaffs a %s and gains %d max health!", creature.Name, potion.Name, power)
		creature.BoostMaxHP(power)
	}
}
