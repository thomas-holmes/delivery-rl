package main

import (
	"github.com/thomas-holmes/delivery-rl/game/items"
)

type Equipment struct {
	Weapon Item
	Shoes  Item
}

func NewEquipment() Equipment {
	def, ok := items.GetCollection("natural_weapons").GetByName("Bare Hands")
	var weapon Item
	if ok {
		weapon = produceItem(def)
	}

	return Equipment{
		Weapon: weapon,
	}
}
