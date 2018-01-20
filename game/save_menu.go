package main

import (
	"bytes"
	"log"

	"io/ioutil"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type SaveMenu struct {
	World *World

	PopMenu
}

func (s SaveMenu) saveGame() {
	log.Printf("OMGOMGOMG SAVING")
	save := NewSaveV0()

	save.SaveWorld(s.World)

	buf := new(bytes.Buffer)

	save.Encode(buf)

	// TODO: Do something with this serialized buffer

	ioutil.WriteFile("/tmp/save.dat", buf.Bytes(), 0666)

	log.Printf("[%v] %v", buf.Len(), buf)

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

func (pop SaveMenu) Render(window *gterm.Window) {
	window.ClearRegion(pop.X, pop.Y, pop.W, pop.H)
	window.PutString(pop.X, pop.Y, "Save? Y/N", White)
}
