package main

import (
	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/controls/scene"

	"github.com/thomas-holmes/gterm"
)

const (
	IntroSceneName = "INTRO_SCENE"
	GameSceneName  = "GAME_SCENE"
)

type IntroScene struct {
	intro *IntroScreen

	quitGame func()

	window *gterm.Window
}

func (i *IntroScene) Name() string {
	return IntroSceneName
}

func (i *IntroScene) OnActivate(previous string) {
	i.intro.Reset()
}

func (i *IntroScene) Update(input controls.InputEvent, deltaT uint32) {
	i.intro.Update(input.Action())

	switch {
	case i.intro.QuitGame():
		i.quitGame()
	case i.intro.StartGame():
		scene.SetActiveScene(GameSceneName)
	}
}

func (i *IntroScene) Render(window *gterm.Window, deltaT uint32) {
	i.intro.Render(i.window)
}
