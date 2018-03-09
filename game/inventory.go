package main

import (
	"fmt"
	"log"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/items"

	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
)

type Inventory []Item

func (inventory *Inventory) Filter(filter items.Kind) []Item {
	var filtered []Item

	for _, item := range *inventory {
		if item.Kind&filter > 0 {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func (inventory *Inventory) Add(item Item) {
	if !item.Stacks {
		*inventory = append(*inventory, item)
		return
	}
	for i, it := range *inventory {
		if it.Name == item.Name {
			(*inventory)[i].Count += item.Count
			return
		}
	}
	*inventory = append(*inventory, item)
}

func (inventory *Inventory) RemoveItem(item Item) {
	for i, it := range *inventory {
		if it.Name == item.Name {
			if item.Count > 1 {
				(*inventory)[i].Count--
				return
			} else {
				*inventory = append((*inventory)[:i], (*inventory)[i+1:]...)
				return
			}
		}
	}
}

type InventoryPop struct {
	Inventory

	selectedIndex int
	itemDetailPop ItemDetails

	PopMenu
}

func (pop *InventoryPop) adjustSelection(delta int) {
	pop.selectedIndex += delta
	pop.selectedIndex = max(0, pop.selectedIndex)
	pop.selectedIndex = min(len(pop.Inventory)-1, pop.selectedIndex)
}

func (pop *InventoryPop) selectedItem() (Item, bool) {
	if len(pop.Inventory) <= 0 {
		return Item{}, false
	}

	return pop.Inventory[pop.selectedIndex], true
}

func (pop *InventoryPop) Update(input controls.InputEvent) {
	item, itemSelected := pop.selectedItem()
	pop.CheckCancel(input)
	switch input.Action() {
	case controls.Up:
		pop.adjustSelection(-1)
	case controls.Down:
		pop.adjustSelection(1)
	case controls.SkipUp:
		pop.adjustSelection(-5)
	case controls.SkipDown:
		pop.adjustSelection(5)
	case controls.Top:
		pop.adjustSelection(-len(pop.Inventory))
	case controls.Bottom:
		pop.adjustSelection(len(pop.Inventory))
	case controls.Quaff:
		if itemSelected && item.CanQuaff() {
			m.Broadcast(m.M{ID: PlayerQuaffPotion, Data: PlayerQuaffPotionMessage{Potion: item}})
			pop.done = true
		}
	case controls.Activate:
		if itemSelected && item.CanActivate() {
			m.Broadcast(m.M{ID: PlayerActivateItem, Data: PlayerActivateItemMessage{Item: item}})
			pop.done = true
		}
	case controls.Equip:
		if itemSelected && item.CanEquip() {
			m.Broadcast(m.M{ID: EquipItem, Data: EquipItemMessage{item}})
			pop.done = true
		}
	}
}

func (pop *InventoryPop) renderItem(index int, row int, window *gterm.Window) int {
	offsetY := row
	offsetX := pop.X + 1

	item := pop.Inventory[index]

	prefix := ""
	if item.Count > 1 {
		prefix = fmt.Sprintf("[%d] ", item.Count)
	}

	color := Grey
	if pop.selectedIndex == index {
		color = White
		window.PutRune(offsetX+1, offsetY, rightArrow, White, gterm.NoColor)
	}

	name := item.Name

	offsetY += putWrappedText(window, prefix+name, offsetX, offsetY, 4, 2, pop.W-offsetX+pop.X-1, color)
	return offsetY
}

func (pop *InventoryPop) Render(window *gterm.Window) {
	if err := window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Printf("(%v,%v) (%v,%v)", pop.X, pop.Y, pop.W, pop.H)
		log.Println("Failed to render inventory", err)
	}

	nextRow := pop.Y + 1
	for i := 0; i < len(pop.Inventory); i++ {
		nextRow = pop.renderItem(i, nextRow, window)
	}

	if item, ok := pop.selectedItem(); ok {
		menu := ItemDetails{PopMenu: PopMenu{X: pop.X + pop.W, Y: pop.Y, W: 30, H: 26}, Item: item}
		menu.Render(window)
	}
	pop.DrawBox(window, White)
}
