package main

import "testing"

func TestResourceTickWholeNumbers(t *testing.T) {
	r := Resource{Current: 1, Max: 5, RateTimes100: 100}
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
	r := Resource{Current: 4, Max: 5, RateTimes100: 100}
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
	r := Resource{Current: 1, Max: 5, RateTimes100: 50}
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
	r := Resource{Current: 1, Max: 5, RateTimes100: 25}
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
	r := Resource{Current: 5, Max: 5, RateTimes100: -50}
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

func TestResourceMinisculeRateTimes100(t *testing.T) {
	r := Resource{Current: 10, Max: 20, RateTimes100: 5}

	for i := 0; i < 20; i++ {
		r.Tick()
	}

	expected := 11
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}

}

func TestResourceDegenHeatBug(t *testing.T) {
	r := Resource{Current: 125, Max: 125, RateTimes100: -10}

	for i := 0; i < 10; i++ {
		r.Tick()
	}

	expected := 124
	if r.Current != expected {
		t.Errorf("Got current value %v, but expected %v. %+v", r.Current, expected, r)
	}

}
