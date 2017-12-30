package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type MonsterDefinitions struct {
	Monsters []MonsterDefinition `toml:"Monster"`
}

type MonsterDefinition struct {
	Name        string
	Description string
	Glyph       string
	Color       []int

	Level int
	HP    int
}

func loadMonsterDefinitions(path string) MonsterDefinitions {
	defs := MonsterDefinitions{}
	_, err := toml.DecodeFile(path, &defs)
	if err != nil {
		log.Panicf("Could not load monster definitions from %v. Error %+v", path, err)
	}

	return defs
}
