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
	intro IntroScreen

	quitGame func()

	window *gterm.Window
}

func (i *IntroScene) Name() string {
	return IntroSceneName
}

func (i *IntroScene) Update(input controls.InputEvent, deltaT uint32) {
	i.intro.Update(input.Action())

	switch {
	case i.intro.QuitGame():
		i.quitGame()
	case i.intro.StartGame():
		i.intro.Reset()
		scene.SetActiveScene(GameSceneName)
	}
}

func (i *IntroScene) Render(deltaT uint32) {
	i.intro.Render(i.window)
}

type GameScene struct {
	world *World
	hud   *HUD
}

func (g *GameScene) Name() string {
	return GameSceneName
}

func (g *GameScene) Update(input controls.InputEvent, deltaT uint32) {
	g.world.Update(input.Action())
	if g.world.QuitGame {
		scene.SetActiveScene(IntroSceneName)
	}
}

func (g *GameScene) Render(deltaT uint32) {
	g.world.UpdateAnimations()
	g.world.Render()
	g.hud.Render(g.world)
}
