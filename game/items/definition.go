package items

import (
	"io/ioutil"
	"log"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type ItemKind int

const (
	Unknown ItemKind = iota
	Consumeable
	Weapon
)

func (i *ItemKind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kindStr string
	if err := unmarshal(&kindStr); err != nil {
		log.Println("Failed to unmarshal ItemKind", err)
	}

	k, ok := parseKind(kindStr)
	if !ok {
		log.Println("Failed to parse item kind")
	}
	*i = k

	return nil
}

func parseKind(kind string) (ItemKind, bool) {
	switch strings.ToLower(kind) {
	case "consumeable":
		return Consumeable, true
	case "weapon":
		return Weapon, true
	}

	return Unknown, false
}

type ItemDefinition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int
	Power       float64
	Kind        ItemKind
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
