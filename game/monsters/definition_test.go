package monsters

import (
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestMonsterParsing(t *testing.T) {
	goblinYAML := `
name: "Goblin"
description: "A hungry little goblin."
glyph: "g"
color: [0, 255, 0]
level: 1
power: 4
hp: 5
`
	monster := Definition{}
	err := yaml.Unmarshal([]byte(goblinYAML), &monster)
	if err != nil {
		t.Errorf("Failed to parse goblin YAML, %+v", err)
	}

	if monster.Name != "Goblin" {
		t.Errorf("Expected name (Goblin) got (%v)", monster.Name)
	}

	if monster.Description != "A hungry little goblin." {
		t.Errorf("Expected description (A hungry littel goblin.) got (%v)", monster.Description)
	}

	if monster.Color.G != 255 {
		t.Errorf("Expected monster color R to be 255, got (%v)", monster.Color.G)
	}

	if monster.Power != 4 {
		t.Errorf("Expected Power (4) but got (%v)", monster.Power)
	}

	if monster.HP != 5 {
		t.Errorf("Expected HP (5) got (%v)", monster.HP)
	}
}

func TestParseMultipleMonsters(t *testing.T) {
	monstersYAML := `---
  - name: "Goblin"
    description: "A hungry little goblin."
    glyph: "g"
    color: [0, 255, 0]
    level: 1
    hp: 5

  - name: "Orc"
    description: "A furious orc."
    glyph: "o"
    color: [25, 225, 10]
    level: 3
    hp:  14
`
	var monsters []Definition
	err := yaml.Unmarshal([]byte(monstersYAML), &monsters)
	if err != nil {
		t.Errorf("Failed to parse goblin YAML, %+v", err)
	}

	if len(monsters) != 2 {
		t.Errorf("Expected (2) monsters got (%v)", len(monsters))
	}

	if monsters[0].Name != "Goblin" {
		t.Errorf("Expected name (Goblin) got (%v)", monsters[0].Name)
	}
	if monsters[0].Description != "A hungry little goblin." {
		t.Errorf("Expected description (A hungry littel goblin.) got (%v)", monsters[0].Description)
	}

	if monsters[1].Name != "Orc" {
		t.Errorf("Expected name (Goblin) got (%v)", monsters[1].Name)
	}
	if monsters[1].Description != "A furious orc." {
		t.Errorf("Expected description (A curious orc.) got (%v)", monsters[1].Description)
	}
	if monsters[1].HP != 14 {
		t.Errorf("Expected HP (14) got (%v)", monsters[1].HP)
	}

}
