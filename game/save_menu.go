package main

import "github.com/veandco/go-sdl2/sdl"

type SaveMenu struct {
	World *World

	PopMenu
}

func (s SaveMenu) saveGame() {
	s.World.SaveGame()
}

func (pop *SaveMenu) Update(input InputEvent) {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		k := e.Keysym.Sym
		switch {
		case k == sdl.K_y:
			pop.saveGame()
			pop.done = true
		case k == sdl.K_n:
			pop.done = true
		}

	}
}
