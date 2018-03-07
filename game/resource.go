package main

import (
	"math"
)

// Resource represents some character resource that hax a max and [re/de]generates over time
type Resource struct {
	RegenRate float64

	Current int
	Max     int

	regenPartial float64
}

// Tick applies the regen rate to the resource
func (resource *Resource) Tick() {
	resource.regenPartial += resource.RegenRate

	if math.Abs(resource.regenPartial) >= 1 {
		adjustment := math.Floor(resource.regenPartial)
		newValue := resource.Current + int(math.Floor(resource.regenPartial))
		newValue = min(newValue, resource.Max)
		newValue = max(newValue, 0)
		resource.Current = newValue
		resource.regenPartial -= adjustment
	}
}

func (resource Resource) Percentage() float64 {
	current := float64(resource.Current)
	max := float64(resource.Max)
	return current / max
}
