package main

import "testing"

func TestResourceTickWholeNumbers(t *testing.T) {
	r := Resource{Current: 1, Max: 5, RegenRate: 1}
	r.Tick()
	expected := 2
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 3
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 4
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}

func TestResourceTickCapsAtMax(t *testing.T) {
	r := Resource{Current: 4, Max: 5, RegenRate: 1}
	r.Tick()
	expected := 5
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}

func TestResourceTickOneHalf(t *testing.T) {
	r := Resource{Current: 1, Max: 5, RegenRate: 0.5}
	r.Tick()
	expected := 1
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 2
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 2
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 3
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}

func TestResourceTickOneQuarter(t *testing.T) {
	r := Resource{Current: 1, Max: 5, RegenRate: 0.25}
	r.Tick()
	expected := 1
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 1
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 1
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 2
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}

func TestResourceNegativeTickOneHalf(t *testing.T) {
	r := Resource{Current: 5, Max: 5, RegenRate: -0.5}
	r.Tick()
	expected := 5
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 4
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 4
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
	r.Tick()
	expected = 3
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}
}
