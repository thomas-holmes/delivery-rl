package main

import (
	"github.com/thomas-holmes/delivery-rl/game/dice"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
	"github.com/thomas-holmes/delivery-rl/game/items"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
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

// ThrowItem chucks a single item across the map. It's a bit weird because of the
// way removing an item works. So it makes a copy, sets the quantity to 1 so
// we don't accidentally create a dup bug, but we need to provide the full count
// so the inventory removeitem function properly treats it as a stack.
// I could fix this a better way but not right now :()
func ThrowItem(creature *Creature, world *World, item Item, targetX, targetY int) bool {
	singleItem := item
	singleItem.Count = 1
	switch singleItem.Name {
	case "Garlic Butter":
		gl.Append("Threw %s", singleItem.Name)
		a := NewLinearSpellAnimation(creature.X, creature.Y, targetX, targetY, 20, 0, singleItem.Symbol, singleItem.Color)
		world.AddAnimation(&a)
		if singleItem.Name == "Garlic Butter" {
			m.Broadcast(m.M{ID: SplashGrease, Data: SplashGreaseMessage{Item: singleItem, X: targetX, Y: targetY}})
		}
		return true
	default:
		if world.PlaceItemAround(singleItem, targetX, targetY) {
			gl.Append("Threw %s", singleItem.Name)
			creature.Inventory.RemoveItem(item) // This one is "item" so it just reduces by one. Ugh
			animation := NewLinearSpellAnimation(creature.X, creature.Y, targetX, targetY, 20, 0, singleItem.Symbol, singleItem.Color)
			world.AddAnimation(&animation)
			return true
		} else {
			gl.Append("Could not throw %s, there is nowhere for it to land", singleItem.Name)
		}
	}
	return false
}
