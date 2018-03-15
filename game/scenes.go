package main

import (
	"github.com/MichaelTJones/pcg"
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

func (i *IntroScene) Render(window *gterm.Window, deltaT uint32) {
	i.intro.Render(i.window)
}

type GameScene struct {
	world *World
	hud   *HUD
	rng   *pcg.PCG64
}

func MakeNewWorld(window *gterm.Window, rng *pcg.PCG64) *World {
	world := NewWorld(window, true, rng)

	// TODO: Roll this up into some kind of registering a system function on the world
	NewCombatSystem(world)

	player := NewPlayer()

	player.Name = "Euclid"

	world.AddEntityToCurrentLevel(player)

	return world
}

func NewGameScene(window *gterm.Window) *GameScene {
	scene := &GameScene{}

	pcgRng := pcg.NewPCG64()
	seed := uint64(Seed)
	pcgRng.Seed(seed, DefaultSeq, seed*seed, DefaultSeq+1)

	seedDice(pcgRng)

	scene.world = MakeNewWorld(window, pcgRng)

	scene.hud = NewHud(scene.world.Player, scene.world, 65, 2)

	return scene
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

func (g *GameScene) Render(window *gterm.Window, deltaT uint32) {
	g.world.UpdateAnimations()
	g.world.Render()
	g.hud.Render(g.world)
}
