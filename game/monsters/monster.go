package monsters

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/veandco/go-sdl2/sdl"

	yaml "gopkg.in/yaml.v2"
)

var Dragon Definition = Definition{
	Name:        "Hungry Dragon",
	Description: "Your best customer, a very hungry dragon. He can't wait to get his claws on some pizza!",
	Glyph:       "D",
	Color:       color{sdl.Color{R: 225, G: 25, B: 25, A: 255}},
	Level:       99,
	HP:          5000,
}

type color struct {
	sdl.Color
}

type Definition struct {
	Name        string
	Description string
	Glyph       string
	Color       color

	Level int
	Power int
	HP    int
}

// Feels a bit error prone, but I think it'll be ok.
func (c *color) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var colors []uint8
	if err := unmarshal(&colors); err != nil {
		log.Panicln("Big surprise. Unmarshal is broken?", err)
	}

	if len(colors) != 3 {
		return errors.New("Expected colors array of format [R, G, B]")
	}

	c.R = colors[0]
	c.G = colors[1]
	c.B = colors[2]
	c.A = 255

	return nil
}

func LoadDefinitions(path string) ([]Definition, error) {
	var defs []Definition

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(bytes, &defs)
	if err != nil {
		return nil, err
	}

	return defs, nil
}
