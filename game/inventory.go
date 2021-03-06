package main

import (
	"fmt"
	"log"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/items"

	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
)

const MaxInventory = 28

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

// Add Returns false if there was no room
func (inventory *Inventory) Add(item Item) bool {
	if len(*inventory) >= MaxInventory {
		return false
	}

	if !item.Stacks {
		*inventory = append(*inventory, item)
		return true
	}
	for i, it := range *inventory {
		if it.Name == item.Name {
			(*inventory)[i].Count += item.Count
			return true
		}
	}
	*inventory = append(*inventory, item)

	return true
}

func (inventory *Inventory) RemoveAllItem(item Item) {
	for i, it := range *inventory {
		if it.Name == item.Name {
			*inventory = append((*inventory)[:i], (*inventory)[i+1:]...)
			return
		}
	}
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

	*World

	selectedIndex int
	itemDetailPop ItemDetails

	PopMenu
}

func (pop *InventoryPop) adjustSelectionWrap(delta int) {
	pop.selectedIndex += delta
	if pop.selectedIndex < 0 {
		pop.selectedIndex = max(0, len(pop.Inventory)-1)
	} else if pop.selectedIndex >= len(pop.Inventory) {
		pop.selectedIndex = 0
	}
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

func (pop *InventoryPop) tryUseSelectedItem(action controls.Action) {
	item, itemSelected := pop.selectedItem()
	if !itemSelected {
		return
	}
	switch {
	case item.CanQuaff():
		if action != controls.Confirm && action != controls.Quaff {
			return
		}
		m.Broadcast(m.M{ID: PlayerQuaffPotion, Data: PlayerQuaffPotionMessage{Potion: item}})
	case item.CanActivate():
		if action != controls.Confirm && action != controls.Activate {
			return
		}
		m.Broadcast(m.M{ID: PlayerActivateItem, Data: PlayerActivateItemMessage{Item: item}})
	case item.CanEquip():
		if action != controls.Confirm && action != controls.Equip {
			return
		}
		m.Broadcast(m.M{ID: EquipItem, Data: EquipItemMessage{item}})
	case item.CanThrow():
		if action != controls.Confirm && action != controls.Throw {
			return
		}
		m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: NewThrowPop(item, pop.World)}})
	default:
		return
	}

	pop.done = true
}

func (pop *InventoryPop) tryDropItem() {
	selectedItem, ok := pop.selectedItem()
	if ok {
		// x, y := pop.Player.X, pop.Player.Y
		//m.Broadcast(m.M{ID: PlaceItem, Data: PlaceItemMessage{Creature: pop.World.Player, Item: selectedItem, TargetX: x, TargetY: y}})
		m.Broadcast(m.M{ID: PlayerDropItem, Data: PlayerDropItemMessage{World: pop.World, Item: selectedItem}})
		pop.done = true
	}
}

func (pop *InventoryPop) Update(action controls.Action) {
	pop.CheckCancel(action)
	switch action {
	case controls.Up:
		pop.adjustSelectionWrap(-1)
	case controls.Down:
		pop.adjustSelectionWrap(1)
	case controls.SkipUp:
		pop.adjustSelection(-5)
	case controls.SkipDown:
		pop.adjustSelection(5)
	case controls.Top:
		pop.adjustSelection(-len(pop.Inventory))
	case controls.Bottom:
		pop.adjustSelection(len(pop.Inventory))
	case controls.Drop:
		pop.tryDropItem()
	case controls.Confirm, controls.Quaff, controls.Activate, controls.Equip, controls.Throw:
		pop.tryUseSelectedItem(action)
	}
}

func (pop *InventoryPop) renderItem(index int, row int, window *gterm.Window) int {
	offsetY := row
	offsetX := pop.X

	item := pop.Inventory[index]

	color := Grey
	if pop.selectedIndex == index {
		color = White
		window.PutRune(offsetX+2, offsetY, rightArrow, White, gterm.NoColor)
	}

	offsetX += 4
	prefix := ""
	window.PutRune(offsetX, offsetY, item.Symbol, item.Color, gterm.NoColor)

	if item.Count > 1 {
		prefix = fmt.Sprintf("[%d] ", item.Count)
	}

	name := item.Name

	offsetY += putWrappedText(window, prefix+name, offsetX, offsetY, 2, 2, pop.W-offsetX+pop.X-1, color)
	return offsetY
}

func (pop *InventoryPop) Render(window *gterm.Window) {
	if err := window.ClearRegion(pop.X, pop.Y, pop.W, pop.H); err != nil {
		log.Println("Failed to render inventory", err)
	}

	nextRow := pop.Y + 1
	for i := 0; i < len(pop.Inventory); i++ {
		nextRow = pop.renderItem(i, nextRow, window)
	}

	if item, ok := pop.selectedItem(); ok {
		menu := ItemDetails{PopMenu: PopMenu{X: pop.X + pop.W, Y: pop.Y, W: 30, H: 30}, Item: item}
		menu.Render(window)
	}
	pop.DrawBox(window, White)
}
