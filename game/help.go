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

func (pop *HelpPop) Update(action controls.Action) {
	pop.CheckCancel(action)
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

	y += 2
	objective1 := "Welcome to DeliveryRL! You are a typical delivery person for a most unusual pizza shop. " +
		"Sometimes your store gets orders from mythical creatures. Today an Ancient Dragon has ordered a pizza and expects " +
		"it to be delivered promptly, and warm! Why does a dragon need delivery? It's not your job to ask those questions, " +
		"but you drew the short straw this time."

	y += putWrappedText(window, objective1, x, y, 2, 0, pop.W-3, White)
	y++

	objective2 := "Race to the depths of the caverns and deliver the pizza to the Dragon who is waiting. " +
		"Survive by avoiding monsters, fighting for your life, and distracting them with some extra food " +
		"you brought along. You also have a special trick up your sleeve, the ability to warp a short distance. " +
		"It tires you out but don't let that slow you down; time is of the essence, after all! Maybe you can " +
		"scrounge up something useful from past adventurers, but remember: even though you can teleport you're no warrior!"

	y += putWrappedText(window, objective2, x, y, 2, 0, pop.W-3, White)

	y++
	window.PutString(x, y, "Good Luck!", White)

	y += 2
	window.PutString(x, y, "You will find the dragon on the 10th floor.", Red)

	content := "2018 7DRL by"
	name := "keipra"
	x, y = pop.X+pop.W-2-len(content)-len(name), pop.Y+pop.H-2
	window.PutString(x, y, content, Grey)

	window.PutString(x+len(content)+1, y, name, KeipraPurple)
}
