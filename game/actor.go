package main

import "github.com/thomas-holmes/delivery-rl/game/controls"

type Actor interface {
	CanAct() bool

	// Update Should return false if we have not completed our turn and need
	// another pass through the game loop
	Update(turn uint64, input controls.InputEvent, world *World) bool

	StartTurn()
	EndTurn()
}
