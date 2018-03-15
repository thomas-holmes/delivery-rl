package main

import (
	"flag"
	"log"
	"path"
	"time"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/controls/scene"
	"github.com/thomas-holmes/delivery-rl/game/dice"
	gl "github.com/thomas-holmes/delivery-rl/game/gamelog"

	"github.com/MichaelTJones/pcg"

	"github.com/thomas-holmes/delivery-rl/game/items"
	"github.com/thomas-holmes/delivery-rl/game/monsters"
	"github.com/thomas-holmes/gterm"
	"github.com/veandco/go-sdl2/sdl"

	_ "net/http/pprof"
)

var quit = false

var i = 0

var testMessageCount = 1

func handleInput(input controls.InputEvent) {
	switch input.Event.(type) {
	case *sdl.QuitEvent:
		quit = true
	}
}

// seedDice seeds the default dice roller with four random values from the world RNG
func seedDice(seed uint64) {
	rng := pcg.NewPCG64()
	rng.Seed(rng.Random(), rng.Random(), rng.Random(), rng.Random())
	rng.Seed(seed, DiceSeq, seed*seed, DiceSeq+1)
	dice.SetDefaultRandomness(rng)
}

func configureItemsRepository() {
	if err := items.Configure(path.Join("assets", "definitions")); err != nil {
		log.Fatalln("Could not configure items repository", err)
	}

	if err := items.EnsureLoaded("consumeables", "weapons", "armour", "natural_weapons"); err != nil {
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

func BuildGameFromSeed() {

}

type KeyDownFilter struct{}

func (f KeyDownFilter) FilterEvent(e sdl.Event, userdata interface{}) bool {

	switch e.(type) {
	case *sdl.KeyDownEvent:
		return true
	case *sdl.QuitEvent:
		return true
	default:
		return false
	}
}

func main() {
	window := gterm.NewWindow(100, 60, Font, FontW, FontH, !NoVSync)

	if err := window.Init(); err != nil {
		log.Fatalln("Failed to Init() window", err)
	}

	window.SetTitle("DeliveryRL")
	window.SetBackgroundColor(gterm.NoColor)

	configureItemsRepository()
	configureMonstersRepository()

	sdl.SetEventFilter(KeyDownFilter{}, nil)

	intro := IntroScreen{}

	scene.AddScene(&IntroScene{intro: intro, quitGame: func() { quit = true }, window: window})
	scene.AddScene(NewGameScene(window))

	scene.SetActiveScene(IntroSceneName)

	lastTicks := sdl.GetTicks()
	gl.Append("Press ? for help!")
	for !quit {
		var input controls.InputEvent
		for {
			event, mod := sdl.PollEvent(), sdl.GetModState()
			if event == nil {
				break
			}
			input = controls.InputEvent{Event: event, Keymod: mod}
		}
		nowTicks := sdl.GetTicks()
		delta := nowTicks - lastTicks
		lastTicks = nowTicks
		handleInput(input)

		scene.UpdateActiveScene(input, delta)

		window.ClearWindow()

		scene.RenderActiveScene(window, delta)

		window.Refresh()
	}
}

var NoVSync = true
var Seed int64
var Font string
var FontW int
var FontH int

func GameSeed() int64 {
	if Seed == -1 {
		return time.Now().UnixNano()
	} else {
		return Seed
	}
}

func init() {
	flag.BoolVar(&NoVSync, "no-vsync", false, "disable vsync")
	flag.Int64Var(&Seed, "seed", -1, "Provide a seed for launching the game")
	flag.StringVar(&Font, "font-path", "assets/font/cp437_16x16.png", "Set font relative file path")
	flag.IntVar(&FontW, "font-width", 16, "pixel width per character")
	flag.IntVar(&FontH, "font-height", 16, "pixel height per character")
	flag.Parse()
}
