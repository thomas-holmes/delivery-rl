package items

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Definition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int
	Power       int
	Equippable  bool
}

func LoadDefinitions(path string) []Definition {
	var defs []Definition
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
