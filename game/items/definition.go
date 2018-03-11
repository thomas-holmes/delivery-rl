package items

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/MichaelTJones/pcg"
	"github.com/thomas-holmes/delivery-rl/game/dice"

	yaml "gopkg.in/yaml.v2"
)

type Kind uint32

const (
	Unknown Kind = 0
	Potion  Kind = 1 << iota
	Warmer
	Weapon
	Armour
	Missile
	Food
)

func GetLevelBoundedItem(rng *pcg.PCG64, collection Collection, levelBound int) Definition {
	for {
		def := collection.Sample(rng)
		if def.Level <= levelBound {
			return def
		}
	}
}

func (i *Kind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var kindStr string
	if err := unmarshal(&kindStr); err != nil {
		log.Println("Failed to unmarshal Kind", err)
	}

	k, ok := parseKind(kindStr)
	if !ok {
		log.Println("Failed to parse item kind")
	}
	*i = k

	return nil
}

func parseKind(kind string) (Kind, bool) {
	switch strings.ToLower(kind) {
	case "potion":
		return Potion, true
	case "warmer":
		return Warmer, true
	case "weapon":
		return Weapon, true
	case "armour":
		return Armour, true
	case "missile":
		return Missile, true
	case "food":
		return Food, true
	}

	return Unknown, false
}

type Definition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int
	Stacks      bool
	Level       int
	Power       dice.Notation
	Kind        Kind
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
