package items

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type ItemKind int

const (
	Weapon ItemKind = iota
)

func parseKind(kind string) (ItemKind, bool) {
	switch kind {
	case "Weapon":
		return Weapon, true
	}

	return -1, false
}

type ItemDefinition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int
	Power       float64
	Equippable  bool
}

func LoadItemDefinitions(path string) []ItemDefinition {
	var defs []ItemDefinition
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicf("Could not load item definitions from %v. Error %+v", path, err)
	}
	err = yaml.Unmarshal(bytes, &defs)
	if err != nil {
		log.Panicf("Could not parse item definitions from %v. Error %+v", path, err)
	}
	return defs
}
