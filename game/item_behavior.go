package main

import (
	"github.com/thomas-holmes/delivery-rl/game/dice"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	"github.com/thomas-holmes/delivery-rl/game/items"
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

func ActivateItem(creature *Creature, item Item) {
	power := dice.Roll(item.Power)
	switch item.Kind {
	case items.Warmer:
		gl.Append("%s stuffs a %s into the pizza bag, regaining %d heat.", creature.Name, item.Name, power)
		creature.RestoreHeat(power)
	}
}
