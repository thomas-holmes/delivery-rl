package monsters

import (
	"errors"
	"log"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/BurntSushi/toml"
)

type definitions struct {
	Monsters []Definition `toml:"Monster"`
}

type color struct {
	sdl.Color
}

type Definition struct {
	Name        string
	Description string
	Glyph       string
	Color       color `toml:"Color"`

	Level int
	HP    int
}

func (c *color) UnmarshalTOML(data interface{}) error {
	colors, ok := data.([]interface{})
	if !ok {
		return errors.New("Colors not an array")
	}

	if len(colors) != 3 {
		return errors.New("Expected colors array of format [R, G, B]")
	}
	r := colors[0].(int64)
	g := colors[1].(int64)
	b := colors[2].(int64)

	c.R = uint8(r)
	c.G = uint8(g)
	c.B = uint8(b)
	c.A = 255

	return nil
}

func LoadDefinitions(path string) []Definition {
	defs := definitions{}
	_, err := toml.DecodeFile(path, &defs)
	if err != nil {
		log.Panicf("Could not load monster definitions from %v. Error %+v", path, err)
	}

	return defs.Monsters
}
