package main

import (
	"github.com/thomas-holmes/delivery-rl/game/items"
	"github.com/veandco/go-sdl2/sdl"
)

var TShirt = Item{
	Name:        "T-Shirt",
	Description: "A T-Shirt with the logo of your pizza shop. They make you wear it...",
	Symbol:      rune('Ω'),
	Color:       sdl.Color{R: 255, G: 0, B: 0, A: 255},
	Count:       1,
	Kind:        items.Armour,
}

type Equipment struct {
	Weapon Item
	Armour Item
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
