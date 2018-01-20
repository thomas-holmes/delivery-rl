package main

import (
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"
)

type ResumeScreen struct {
	PopMenu
}

func (resume *ResumeScreen) Update(input InputEvent) {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_y:
			// Reload it
			resume.done = true
		case sdl.K_n:
			// Trash it
			resume.done = true
		}
	}
}

func (resume ResumeScreen) Render(window *gterm.Window) {
	window.ClearRegion(resume.X, resume.Y, resume.W, resume.Y)
	window.PutString(resume.X, resume.Y, "Resume game? Y/N", Red)
}
