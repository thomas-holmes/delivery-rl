package main

import (
	"bytes"
	"log"
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

	p.Items = append(p.Items, Item{Name: "TestItem", Power: 9})
	p.Items = append(p.Items, Item{Name: "Mallet", Power: 3})
	p.Equipment.Weapon = p.Items[1]
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

func TestEncodeLevel(t *testing.T) {
	l := Level{}
	l.Columns = 9
	l.Rows = 9
	l.Depth = 1
	l.VisionMap = &VisionMap{}
	l.ScentMap = &ScentMap{}

	buffer := new(bytes.Buffer)

	el := exportLevel(l)

	if err := el.Encode(buffer); err != nil {
		t.Error(err)
	}
	log.Printf("%+v", el)

	decoded := &ExportedLevelV0{}

	serialized := buffer.Bytes()
	log.Printf("%+v", serialized)

	if err := decoded.Decode(bytes.NewReader(serialized)); err != nil {
		t.Fatal(err)
	}

	if decoded.Columns != l.Columns {
		t.Errorf("Exported columns (%v) but got (%v)", l.Columns, decoded.Columns)
	}

	if decoded.Rows != l.Rows {
		t.Errorf("Exported rows (%v) but got (%v)", l.Rows, decoded.Rows)
	}

	if decoded.Depth != l.Depth {
		t.Errorf("Exported depth (%v) but got (%v)", l.Depth, decoded.Depth)
	}
}
