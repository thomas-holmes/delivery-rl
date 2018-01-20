package main

import (
	"bytes"
	"testing"
)

func makeTestPlayer() *Creature {
	p := NewPlayer()
	p.Name = "TestMan"
	p.LevelUp()
	p.LevelUp()
	p.SetIdentity(9)
	p.X = 4
	p.Y = 19

	p.Items = append(p.Items, &Item{Name: "TestItem", Power: 9})
	p.Items = append(p.Items, &Item{Name: "Mallet", Power: 3})
	p.Equipment.Weapon = *p.Items[1]
	return &p
}

func TestExportCreature(t *testing.T) {
	p := makeTestPlayer()
	e := ExportCreature(p)

	if p.Name != e.Name {
		t.Errorf("Expected exported name (%v) but got (%v) instead", p.Name, e.Name)
	}
}

func TestEncodeDecodeCreature(t *testing.T) {
	p := makeTestPlayer()

	e := ExportCreature(p)

	buffer := new(bytes.Buffer)

	if err := e.Encode(buffer); err != nil {
		t.Fatal(err)
	}

	decoded := &ExportedCreatureV0{}

	if err := decoded.Decode(buffer); err != nil {
		t.Fatal(err)
	}

	if decoded.Name != e.Name {
		t.Errorf("Expected decoded name (%v) but instead got (%v)", e.Name, decoded.Name)
	}

	if len(decoded.Inventory) != len(e.Inventory) {
		t.Errorf("Expected decoded items len (%v) but instead got (%v)", len(e.Inventory), len(decoded.Inventory))
	}
	if decoded.Equipment.Weapon != e.Equipment.Weapon {
		t.Errorf("Expected equipped weapon (%+v) but instead got (%+v)", e.Equipment.Weapon, decoded.Equipment.Weapon)
	}
}
