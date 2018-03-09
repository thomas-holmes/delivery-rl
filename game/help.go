package main

import (
	"fmt"
	"strings"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/gterm"
)

type HelpPop struct {
	PopMenu
}

func NewHelpPop(x, y, w, h int) *HelpPop {
	return &HelpPop{
		PopMenu: PopMenu{
			X: x,
			Y: y,
			W: w,
			H: h,
		},
	}
}

func (pop *HelpPop) Update(input controls.InputEvent) {
	pop.CheckCancel(input)
}

func (pop HelpPop) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	pop.DrawBox(window, White)

	x := pop.X + 2
	y := pop.Y + 2

	window.PutString(x, y, "DeliveryRL Controls", White)

	y += 2

	for _, m := range controls.AllMappings {
		if m.HideHelp {
			continue
		}
		var padded []string
		for _, key := range m.Keys {
			padded = append(padded, fmt.Sprintf("%-8s", key))
		}
		window.PutString(x, y, fmt.Sprintf("%-20s %-30s", m.Name, strings.Join(padded, " ")), White)
		y++
	}
}
