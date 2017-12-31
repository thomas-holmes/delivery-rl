package main

type Actor interface {
	CanAct() bool

	// Update Should return false if we have not completed our turn and need
	// another pass through the game loop
	Update(turn uint64, input InputEvent, world *World) bool
}
