package items

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/thomas-holmes/delivery-rl/game/dice"

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

type Definition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int
	Power       dice.Notation
	Kind        ItemKind
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
