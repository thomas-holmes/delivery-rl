package items

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestItemParsing(t *testing.T) {
	appleYAML := `
name: "Apple"
description: "A delicious green apple"
glyph: "a"
color: [0, 255, 0]
equippable: false
`

	item := WeaponDefinition{}

	err := yaml.Unmarshal([]byte(appleYAML), &item)
	if err != nil {
		t.Errorf("Failed to decode apple %v", err)
	}

	if item.Name != "Apple" {
		t.Errorf("Expected name (Apple) but got (%v)", item.Name)
	}
}

func TestParseMultipleWeapons(t *testing.T) {
	itemsYAML := `---
  - name: "Dagger"
    description: "A plain dagger."
    glyph: ")"
    color: [255, 0, 0]
    power: 2
    equippable: true
  - name: "Rapier"
    description: "A slender weapon, perfect for thursting."
    glyph: ')'
    color: [0, 255, 0]
    power: 4
    equippable: true
`

	var items []WeaponDefinition
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
