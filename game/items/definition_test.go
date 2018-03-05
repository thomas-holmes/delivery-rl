package items

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestItemParsing(t *testing.T) {
	healthPotionYAML := `
name: "Health Potion"
description: "Chug this to heal some HP!"
glyph: "!"
color: [255, 0, 0]
power: 2d6+4
kind: potion
`

	item := Definition{}

	err := yaml.Unmarshal([]byte(healthPotionYAML), &item)
	if err != nil {
		t.Errorf("Failed to decode Health Potion %v", err)
	}

	if item.Name != "Health Potion" {
		t.Errorf("Expected name (Health Potion) but got (%v)", item.Name)
	}

	if item.Kind != Potion {
		t.Errorf("Expected item to be Consumeable(%v) but instead got %v", Potion, item.Kind)
	}

	if item.Power.String() != "2d6+4" {
		t.Error("Expected power to be 2d6+4 instead got", item.Power.String())
	}

}

func TestParseMultipleWeapons(t *testing.T) {
	itemsYAML := `---
  - name: "Dagger"
    description: "A plain dagger."
    glyph: ")"
    color: [255, 0, 0]
    power: 1d4
    equippable: true
  - name: "Rapier"
    description: "A slender weapon, perfect for thrusting."
    glyph: ')'
    color: [0, 255, 0]
    power: 1d8
    equippable: true
`

	var items []Definition
	err := yaml.Unmarshal([]byte(itemsYAML), &items)
	if err != nil {
		t.Errorf("Failed to decode items %v", err)
	}

	if len(items) != 2 {
		t.Errorf("Expected item length (2) but got (%v)", len(items))
	}

	if items[0].Name != "Dagger" {
		t.Errorf("Expected name (Dagger) but got (%v)", items[0].Name)
	}

	if items[1].Name != "Rapier" {
		t.Errorf("Expected name (Rapier) but got (%v)", items[1].Name)
	}
}
