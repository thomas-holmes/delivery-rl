package main

import (
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"
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
	if index >= len(pop.Player.Inventory.Filter(items.Weapon)) {
		return
	}

	item := pop.Player.Inventory.Filter(items.Weapon)[index]
	if item.Kind != items.Weapon {
		gl.Append("You can't wield a %s", item.Name)
		return
	}

	m.Broadcast(m.M{ID: EquipItem, Data: EquipItemMessage{item}})
	pop.done = true
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
		Inventory: pop.Player.Inventory.Filter(items.Weapon),
		PopMenu: PopMenu{
			X: pop.X,
			Y: pop.Y,
			W: pop.W,
			H: pop.H,
		},
	}

	inventoryPop.Render(window)
}
