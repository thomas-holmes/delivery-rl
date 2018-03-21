package main

import (
	"log"

	"github.com/MichaelTJones/pcg"
	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/controls/scene"
	"github.com/thomas-holmes/delivery-rl/game/gamelog"
	"github.com/thomas-holmes/delivery-rl/game/messages"
	"github.com/thomas-holmes/gterm"
)

type GameScene struct {
	world  *World
	window *gterm.Window
	hud    *HUD
	rng    *pcg.PCG64
}

const (
	DefaultSeq uint64 = iota * 1000
	DiceSeq
)

func clearSystems() {
	gamelog.Clear()
	messages.UnSubAll()
}

func (scene *GameScene) OnActivate(previous string) {
	clearSystems()
	gamelog.Append("Press ? for controls and more info!")
	pcgRng := pcg.NewPCG64()
	seed := uint64(GameSeed())
	log.Printf("Starting game with seed %d", seed)
	pcgRng.Seed(seed, DefaultSeq, seed*seed, DefaultSeq+1)

	seedDice(seed)

	scene.world = MakeNewWorld(scene.window, pcgRng)
	NewCombatSystem(scene.world)

	scene.hud = NewHud(scene.world.Player, scene.world, 65, 1)
}

func MakeNewWorld(window *gterm.Window, rng *pcg.PCG64) *World {
	world := NewWorld(window, true, rng)

	player := NewPlayer()

	player.Name = "Euclid"

	world.AddCreatureToCurrentLevel(player)

	return world
}

func NewGameScene(window *gterm.Window) *GameScene {
	scene := &GameScene{window: window}
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
	// TODO: What is up with this comment?
	//g.world.UpdateAnimations()
	g.world.Render()
	g.hud.Render(g.world)
}
