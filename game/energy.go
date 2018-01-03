package main

// Very similar to resource, but doesn't work on the same ticking interval
type Energy struct {
	Current int
	Max     int
}

func (energy *Energy) AddEnergy(e int) {
	newEnergy := min(energy.Max, energy.Current+e)
	energy.Current = newEnergy
}
