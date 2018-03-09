package main

import (
	"fmt"
	"log"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/items"

	m "github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
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

	PopMenu
}

func (pop *InventoryPop) tryShowItem(index int) {
	if index < len(pop.Inventory) {
		var unsub m.Unsubscribe
		unsub = m.Subscribe(func(message m.M) {
			if message.ID == ItemDetailClosed {
				if d, ok := message.Data.(ItemDetailClosedMessage); ok {
					if unsub != nil {
						unsub()
					}
					if d.CloseInventory {
						pop.done = true
					}
				}
			}
		})

		menu := ItemDetails{PopMenu: PopMenu{X: 2, Y: 2, W: 50, H: 26}, Item: pop.Inventory[index]}
		m.Broadcast(m.M{ID: ShowMenu, Data: ShowMenuMessage{Menu: &menu}})
	}
}

func (pop *InventoryPop) Update(input controls.InputEvent) {
	pop.CheckCancel(input)
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		k := e.Keysym.Sym
		switch {
		case k >= sdl.K_a && k <= sdl.K_z:
			pop.tryShowItem(int(k - sdl.K_a))
		}
	}
}

func (pop *InventoryPop) renderItem(index int, row int, window *gterm.Window) int {
	offsetY := row
	offsetX := pop.X + 1

	item := pop.Inventory[index]

	var selectionStr string
	if item.Count > 1 {
		selectionStr = fmt.Sprintf("%v - [%d] ", string('a'+index), item.Count)
	} else {
		selectionStr = fmt.Sprintf("%v - ", string('a'+index))
	}

	window.PutString(offsetX, offsetY, selectionStr, White)

	name := item.Name

	offsetY += putWrappedText(window, name, offsetX, offsetY, len(selectionStr), 2, pop.W-offsetX+pop.X-1, White)
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

	pop.DrawBox(window, White)
}
