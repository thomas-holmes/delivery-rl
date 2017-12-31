package items

import (
	"log"

	"github.com/BurntSushi/toml"
)

type definitions struct {
	Items []Definition `toml:"Item"`
}

type Definition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int
	Power       int
	Equippable  bool
}

func LoadDefinitions(path string) []Definition {
	defs := definitions{}
	_, err := toml.DecodeFile(path, &defs)
	if err != nil {
		log.Panicf("Could not load item definitions from %v. Error %+v", path, err)
	}

	return defs.Items
}
