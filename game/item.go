package main

import (
	"fmt"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/dice"
	m "github.com/thomas-holmes/delivery-rl/game/messages"

	"github.com/thomas-holmes/delivery-rl/game/items"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type Item struct {
	Name        string
	Description string
	Symbol      rune
	Color       sdl.Color
	Power       dice.Notation
	Stacks      bool
	Count       int

	Kind items.Kind
}

func (i Item) CanQuaff() bool {
	return i.Kind == items.Potion
}

func (i Item) CanActivate() bool {
	return i.Kind == items.Warmer
}

func produceItem(itemDef items.Definition) Item {
	if itemDef.Name == "" {
		return Item{}
	}

	r, g, b := uint8(itemDef.Color[0]), uint8(itemDef.Color[1]), uint8(itemDef.Color[2])
	return Item{
		Color:       sdl.Color{R: r, G: g, B: b, A: 255},
		Symbol:      []rune(itemDef.Glyph)[0],
		Name:        itemDef.Name,
		Description: itemDef.Description,
		Power:       itemDef.Power,
		Kind:        itemDef.Kind,
		Stacks:      itemDef.Stacks,
		Count:       1,
	}
}

type ItemDetails struct {
	Item

	PopMenu
}

func (pop *ItemDetails) Update(input controls.InputEvent) {
	switch input.Action() {
	case controls.Cancel:
		m.Broadcast(m.M{ID: ItemDetailClosed, Data: ItemDetailClosedMessage{CloseInventory: false}})
		pop.done = true
	case controls.Quaff:
		if pop.Item.CanQuaff() {
			m.Broadcast(m.M{ID: PlayerQuaffPotion, Data: PlayerQuaffPotionMessage{Potion: pop.Item}})
			m.Broadcast(m.M{ID: ItemDetailClosed, Data: ItemDetailClosedMessage{CloseInventory: true}})
			pop.done = true
		}
	case controls.Activate:
		if pop.Item.CanActivate() {
			m.Broadcast(m.M{ID: PlayerActivateItem, Data: PlayerActivateItemMessage{Item: pop.Item}})
			m.Broadcast(m.M{ID: ItemDetailClosed, Data: ItemDetailClosedMessage{CloseInventory: true}})
			pop.done = true
		}
	}
}

func (pop *ItemDetails) renderShortDescription(row int, window *gterm.Window) int {
	offsetX := pop.X + 1
	offsetY := row
	window.PutRune(offsetX, offsetY, pop.Item.Symbol, pop.Item.Color, gterm.NoColor)
	nameStr := fmt.Sprintf(" - %v", pop.Item.Name)
	offsetX += 2

	window.PutString(offsetX, offsetY, nameStr, White)
	return offsetY + 1
}

func (pop *ItemDetails) renderLongDescription(row int, window *gterm.Window) int {
	offsetX := pop.X + 2
	offsetY := row + 1

	description := pop.Item.Description
	offsetY += putWrappedText(window, description, offsetX, offsetY, 4, 0, pop.W-offsetX+pop.X-1, White)

	return offsetY
}

func (pop *ItemDetails) renderPower(row int, window *gterm.Window) int {
	offsetX := pop.X + 1
	offsetY := row + 1

	powerString := "Power: "
	window.PutString(offsetX, offsetY, powerString, White)

	offsetX += len(powerString)
	window.PutString(offsetX, offsetY, pop.Item.Power.String(), pop.Item.Color)

	return offsetY + 1
}

func (pop *ItemDetails) renderUsage(window *gterm.Window) {
	var usageStr string
	switch pop.Item.Kind {
	case items.Potion:
		usageStr = "Actions: (Q)uaff"
	case items.Warmer:
		usageStr = "Actions: (A)ctivate"
	}
	if usageStr != "" {
		window.PutString(pop.X+1, pop.Y+pop.H-2, usageStr, White)
	}
}

func (pop *ItemDetails) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	pop.DrawBox(window, White)
	row := pop.Y + 1
	row = pop.renderShortDescription(row, window)
	row = pop.renderLongDescription(row, window)
	_ = pop.renderPower(row, window)
	pop.renderUsage(window)
}
