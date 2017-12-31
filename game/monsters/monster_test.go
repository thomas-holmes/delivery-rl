package monsters

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestMonsterParsing(t *testing.T) {
	goblinTOML := `
Name = "Goblin"
Description = "A hungry little goblin."
Glyph = "g"
Color = [0, 255, 0] 
Level = 1
HP = 5
`
	monster := Definition{}
	_, err := toml.Decode(goblinTOML, &monster)
	if err != nil {
		t.Errorf("Failed to parse goblin TOML, %+v", err)
	}

	if monster.Name != "Goblin" {
		t.Errorf("Expected name (Goblin) got (%v)", monster.Name)
	}

	if monster.Description != "A hungry little goblin." {
		t.Errorf("Expected description (A hungry littel goblin.) got (%v)", monster.Description)
	}
}

func TestParseMultipleMonsters(t *testing.T) {
	goblinTOML := `
[[Monster]]
Name = "Goblin"
Description = "A hungry little goblin."
Glyph = "g"
Color = [0, 255, 0] 
Level = 1
HP = 5

[[Monster]]
Name = "Orc"
Description = "A furious orc."
Glyph = "o"
Color = [25, 225, 10] 
Level = 2
HP = 10
`
	monsters := definitions{}
	_, err := toml.Decode(goblinTOML, &monsters)
	if err != nil {
		t.Errorf("Failed to parse goblin TOML, %+v", err)
	}

	if len(monsters.Monsters) != 2 {
		t.Errorf("Expected (2) monsters got (%v)", len(monsters.Monsters))
	}

	if monsters.Monsters[0].Name != "Goblin" {
		t.Errorf("Expected name (Goblin) got (%v)", monsters.Monsters[0].Name)
	}
	if monsters.Monsters[0].Description != "A hungry little goblin." {
		t.Errorf("Expected description (A hungry littel goblin.) got (%v)", monsters.Monsters[0].Description)
	}

	if monsters.Monsters[1].Name != "Orc" {
		t.Errorf("Expected name (Goblin) got (%v)", monsters.Monsters[1].Name)
	}
	if monsters.Monsters[1].Description != "A furious orc." {
		t.Errorf("Expected description (A curious orc.) got (%v)", monsters.Monsters[1].Description)
	}

}
