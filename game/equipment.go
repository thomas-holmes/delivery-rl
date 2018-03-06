package main

import (
	"log"

	"github.com/thomas-holmes/delivery-rl/game/items"
	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
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

type EquipmentPop struct {
	Player *Creature

	PopMenu
}

func (pop *EquipmentPop) equipItem(index int) {
	// Should probably filter on equippable and change the list, whatever
	// Consider doing this w/o message broadcast. We do have a ref to the player, after all
	if index < len(pop.Player.Inventory.Items) {
		item := pop.Player.Inventory.Items[index]
		log.Printf("Equipping item %+v", item)
		m.Broadcast(m.M{ID: EquipItem, Data: EquipItemMessage{item}})
		pop.done = true
	}
}

func (pop *EquipmentPop) Update(input InputEvent) {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		k := e.Keysym.Sym
		switch {
		case k == sdl.K_ESCAPE:
			pop.done = true
		case k >= sdl.K_a && k <= sdl.K_z:
			pop.equipItem(int(k - sdl.K_a))
		}

	}
}

func (pop *EquipmentPop) Render(window *gterm.Window) {
	// TODO: Don't do this
	inventoryPop := InventoryPop{
		Inventory: pop.Player.Inventory,
		PopMenu: PopMenu{
			X: pop.X,
			Y: pop.Y,
			W: pop.W,
			H: pop.H,
		},
	}

	inventoryPop.Render(window)
}
