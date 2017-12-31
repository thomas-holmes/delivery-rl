package monsters

import (
	"log"

	"github.com/BurntSushi/toml"
)

type definitions struct {
	Monsters []Definition `toml:"Monster"`
}

type Definition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int

	Level int
	HP    int
}

func LoadDefinitions(path string) []Definition {
	defs := definitions{}
	_, err := toml.DecodeFile(path, &defs)
	if err != nil {
		log.Panicf("Could not load monster definitions from %v. Error %+v", path, err)
	}

	return defs.Monsters
}
