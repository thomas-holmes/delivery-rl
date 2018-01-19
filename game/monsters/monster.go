package monsters

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/veandco/go-sdl2/sdl"

	yaml "gopkg.in/yaml.v2"
)

type color struct {
	sdl.Color
}

type Definition struct {
	Name        string
	Description string
	Glyph       string
	Color       color

	Level int
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

func LoadDefinitions(path string) []Definition {
	var defs []Definition

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("Could not load monster definitions from %v. Error %+v", path, err)
	}
	err = yaml.Unmarshal(bytes, &defs)
	if err != nil {
		log.Panicf("Could not load monster definitions from %v. Error %+v", path, err)
	}

	return defs
}
