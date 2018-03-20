package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/thomas-holmes/delivery-rl/game/controls"
	"github.com/thomas-holmes/delivery-rl/game/controls/scene"
	"github.com/thomas-holmes/delivery-rl/game/dice"

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
	if err := items.Configure(path.Join(AssetRoot, "assets", "definitions")); err != nil {
		log.Fatalln("Could not configure items repository", err)
	}

	if err := items.EnsureLoaded("consumeables", "weapons", "armour", "natural_weapons"); err != nil {
		log.Fatalln("Failed to load all item repositories", err)
	}
}

func configureMonstersRepository() {
	if err := monsters.Configure(path.Join(AssetRoot, "assets", "definitions")); err != nil {
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
	case *sdl.WindowEvent:
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
	window.SetBackgroundColor(DeepPurple)

	configureItemsRepository()
	configureMonstersRepository()

	sdl.SetEventFilter(KeyDownFilter{}, nil)

	scene.AddScene(&IntroScene{intro: NewIntroScreen(window), quitGame: func() { quit = true }, window: window})
	scene.AddScene(NewGameScene(window))

	scene.SetActiveScene(IntroSceneName)

	lastTicks := sdl.GetTicks()
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

var AssetRoot string

var DefaultFontPath string

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
	var err error
	if AssetRoot, err = os.Executable(); err != nil {
		log.Println("Failed to determine executable path, loading assets relative to working directory")
		AssetRoot = ""
	} else {
		AssetRoot = filepath.Dir(AssetRoot)
	}
	DefaultFontPath = path.Join(AssetRoot, "assets", "font", "cp437_12x12.png")

	flag.BoolVar(&NoVSync, "no-vsync", false, "disable vsync")
	flag.Int64Var(&Seed, "seed", -1, "Provide a seed for launching the game")
	flag.StringVar(&Font, "font-path", DefaultFontPath, "Set font relative file path")
	flag.IntVar(&FontW, "font-width", 12, "pixel width per character")
	flag.IntVar(&FontH, "font-height", 12, "pixel height per character")
	log.Println("DefaultFontPath", DefaultFontPath)
	flag.Parse()
}
