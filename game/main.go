package main

import (
	"flag"
	"log"
	"path"
	"time"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/dice"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"

	"github.com/MichaelTJones/pcg"

	"github.com/thomas-holmes/delivery-rl/game/items"
	"github.com/thomas-holmes/delivery-rl/game/monsters"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"

	"net/http"
	_ "net/http/pprof"
)

var quit = false

var i = 0

func filterActionableEvents(input controls.InputEvent) controls.InputEvent {
	switch e := input.Event.(type) {
	case *sdl.KeyDownEvent:
		input.Event = *e
		return input
	case *sdl.QuitEvent:
		input.Event = *e
		return input
	}
	return controls.InputEvent{Event: nil}
}

var showFPS = true

var testMessageCount = 1

func handleInput(input controls.InputEvent, world *World) {
	switch e := input.Event.(type) {
	case sdl.KeyDownEvent:
		switch e.Keysym.Sym {
		case sdl.K_BACKSLASH:
			world.ToggleScentOverlay()
		case sdl.K_F12:
			showFPS = !showFPS
			world.Window.ShouldRenderFps(showFPS)
		case sdl.K_F10:
			gl.Append("This is a test message %d", testMessageCount)
			testMessageCount++
		case sdl.K_F11:
			gl.Append("This is a test very long message, its number is %d so you can figure it out", testMessageCount)
			testMessageCount++
		}
		world.Input = input
	case sdl.QuitEvent:
		quit = true
	}
}

const (
	DefaultSeq uint64 = iota * 1000
)

func MakeNewWorld(window *gterm.Window, rng *pcg.PCG64) *World {
	world := NewWorld(window, true, rng)

	// TODO: Roll this up into some kind of registering a system function on the world
	NewCombatSystem(world)

	player := NewPlayer()

	player.Name = "Euclid"

	world.AddEntityToCurrentLevel(player)

	return world
}

// seedDice seeds the default dice roller with four random values from the world RNG
func seedDice(pcgRng *pcg.PCG64) {
	rng := pcg.NewPCG64()
	rng.Seed(rng.Random(), rng.Random(), rng.Random(), rng.Random())
	dice.SetDefaultRandomness(rng)
}

func configureItemsRepository() {
	if err := items.Configure(path.Join("assets", "definitions")); err != nil {
		log.Fatalln("Could not configure items repository", err)
	}

	if err := items.EnsureLoaded("consumeables", "weapons", "shoes", "natural_weapons"); err != nil {
		log.Fatalln("Failed to load all item repositories", err)
	}
}

func configureMonstersRepository() {
	if err := monsters.Configure(path.Join("assets", "definitions")); err != nil {
		log.Fatalln("Could not configure monsters repository", err)
	}

	if err := monsters.EnsureLoaded("monsters"); err != nil {
		log.Fatalln("Could not load monster repositories", err)
	}
}

func main() {
	// Disable FPS limit, generally, so I can monitor performance.
	window := gterm.NewWindow(100, 30, path.Join("assets", "font", "MorePerfectDOSVGA.ttf"), 26, !NoVSync)

	pcgRng := pcg.NewPCG64()
	seed := uint64(Seed)
	pcgRng.Seed(seed, DefaultSeq, seed*seed, DefaultSeq+1)

	if err := window.Init(); err != nil {
		log.Fatalln("Failed to Init() window", err)
	}

	window.SetTitle("DeliveryRL")
	window.SetBackgroundColor(gterm.NoColor)

	configureItemsRepository()
	configureMonstersRepository()

	seedDice(pcgRng)
	window.ShouldRenderFps(showFPS)
	world := MakeNewWorld(window, pcgRng)

	hud := NewHud(world.Player, world, 60, 0)

	intro := IntroScreen{}
	for !quit && !world.QuitGame {
		for {
			event, mod := sdl.PollEvent(), sdl.GetModState()
			if event == nil {
				break
			}
			inputEvent := controls.InputEvent{Event: event, Keymod: mod}
			inputEvent = filterActionableEvents(inputEvent)
			handleInput(inputEvent, world)
			if !intro.Done() {
				intro.Update(inputEvent)
			}
		}
		world.Update()
		world.Input = controls.InputEvent{}

		window.ClearWindow()

		if world.Animating() {
			world.UpdateAnimations()
		}

		if !intro.Done() {
			intro.Render(window)
		} else {
			world.Render()
			hud.Render(world)
		}
		window.Refresh()
	}
}

var NoVSync = true
var Seed int64

func init() {
	go http.ListenAndServe("localhost:6060", nil)
	flag.BoolVar(&NoVSync, "no-vsync", false, "disable vsync")
	flag.Int64Var(&Seed, "seed", time.Now().UnixNano(), "Provide a seed for launching the game")
	flag.Parse()
	log.Println("Starting game with seed", Seed)
}
