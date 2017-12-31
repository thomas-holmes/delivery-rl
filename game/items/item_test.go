package items

import (
	"testing"

	"github.com/BurntSushi/toml"
)

func TestItemParsing(t *testing.T) {
	appleTOML := `
Name = "Apple"
Description = "A delicious green apple"
Glyph = "a"
Color = [0, 255, 0]
Equippable = false
`

	item := Definition{}

	_, err := toml.Decode(appleTOML, &item)
	if err != nil {
		t.Errorf("Failed to decode apple %v", err)
	}

	if item.Name != "Apple" {
		t.Errorf("Expected name (Apple) but got (%v)", item.Name)
	}
}

func TestParseMultipleItems(t *testing.T) {
	itemsTOML := `
[[Item]]
Glyph = ")"
Color = [255, 0, 0]
Name = "Dagger"
Description = "A plain dagger."
Power = 2
Equippable = true

[[Item]]
Name = "Rapier"
Description = "A slender weapon, perfect for thursting."
Glyph = ')'
Color = [0, 255, 0]
Power = 4
Equippable = true
`

	defs := definitions{}
	_, err := toml.Decode(itemsTOML, &defs)
	if err != nil {
		t.Errorf("Failed to decode items %v", err)
	}

	items := defs.Items
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
