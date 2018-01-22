package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type ResumeScreen struct {
	world *World

	PopMenu
}

func (resume *ResumeScreen) Update(input InputEvent) {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_y:
			s := SaveV0{}
			dat, err := ioutil.ReadFile("/tmp/save.dat")
			if err != nil {
				log.Fatalln("Couldn't load save", err)
			}

			if err := s.Decode(bytes.NewReader(dat)); err != nil {
				log.Println("Failed to properly decode", err)
			}

			s.Restore(resume.world)
			resume.done = true
		case sdl.K_n:
			resume.done = true
		}
	}
}

func (resume ResumeScreen) Render(window *gterm.Window) {
	window.ClearRegion(resume.X, resume.Y, resume.W, resume.Y)
	window.PutString(resume.X, resume.Y, "Resume game? Y/N", Red)
}
