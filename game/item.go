package main

import (
	"fmt"
	"strconv"

	"github.com/thomas-holmes/delivery-rl/game/items"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

// This is used for empty hands, maybe?
var BareHands = Item{Name: "Bare Hands", Power: 1}

type Item struct {
	Name        string
	Description string
	Symbol      rune
	Color       sdl.Color
	Power       float64

	Kind items.ItemKind
}

func produceItem(itemDef items.Definition) Item {
	r, g, b := uint8(itemDef.Color[0]), uint8(itemDef.Color[1]), uint8(itemDef.Color[2])
	return Item{
		Color:       sdl.Color{R: r, G: g, B: b, A: 255},
		Symbol:      []rune(itemDef.Glyph)[0],
		Name:        itemDef.Name,
		Description: itemDef.Description,
		Power:       itemDef.Power,
		Kind:        itemDef.Kind,
	}
}

type ItemDetails struct {
	Item

	PopMenu
}

func (pop *ItemDetails) Update(input InputEvent) {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_ESCAPE:
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
	window.PutString(offsetX, offsetY, strconv.Itoa(int(pop.Item.Power)), pop.Item.Color)

	return offsetY + 1
}

func (pop *ItemDetails) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	row := pop.Y + 1
	row = pop.renderShortDescription(row, window)
	row = pop.renderLongDescription(row, window)
	_ = pop.renderPower(row, window)
}
