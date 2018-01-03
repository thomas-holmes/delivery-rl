package monsters

import (
	"errors"
	"fmt"
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

// Feels a bit error prone, but I think it'll be ok.
func (c *color) UnmarshalTOML(data interface{}) error {
	colors, ok := data.([]interface{})
	if !ok {
		return errors.New("Colors not an array")
	}

	if len(colors) != 3 {
		return errors.New("Expected colors array of format [R, G, B]")
	}
	r, ok := colors[0].(int64)
	if !ok {
		return fmt.Errorf("R %v was not an int64 but %T instead, failed", colors[0], colors[0])
	}
	g, ok := colors[1].(int64)
	if !ok {
		return fmt.Errorf("G %v was not an int64 but %T instead, failed", colors[1], colors[1])
	}
	b, ok := colors[2].(int64)
	if !ok {
		return fmt.Errorf("B %v was not an int64 but %T instead, failed", colors[2], colors[2])
	}

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
