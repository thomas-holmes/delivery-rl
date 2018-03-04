package items

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type WeaponDefinition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int
	Power       int
	Equippable  bool
}

func LoadWeaponDefinitions(path string) []WeaponDefinition {
	var defs []WeaponDefinition
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
