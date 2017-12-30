package main

import "math"

type Resource struct {
	RegenRate float64

	Current int
	Max     int

	regenPartial float64
}

func (resource *Resource) Tick() {
	resource.regenPartial += resource.RegenRate

	if math.Floor(resource.regenPartial) >= 1 {
		resource.Current = min(int(math.Floor(resource.regenPartial))+resource.Current, resource.Max)
		resource.regenPartial -= math.Floor(resource.regenPartial)
	}
}
