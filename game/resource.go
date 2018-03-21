package main

// Resource represents some character resource that hax a max and [re/de]generates over time
type Resource struct {
	RateTimes100 int

	Current int
	Max     int

	regenPartial int
}

// Tick applies the regen rate to the resource
func (resource *Resource) Tick() {
	resource.regenPartial += resource.RateTimes100

	sign := 1

	partial := resource.regenPartial
	if partial < 0 {
		partial = -partial
		sign = -1
	}

	actualRegen := partial / 100
	if actualRegen > 0 {
		adjustment := actualRegen * sign
		resource.Current += adjustment
		resource.regenPartial -= (adjustment * 100)
	}

	resource.Current = max(0, resource.Current)
	resource.Current = min(resource.Max, resource.Current)

}

func (resource Resource) Percentage() float64 {
	current := float64(resource.Current)
	max := float64(resource.Max)
	return current / max
}
