package main

import (
	"github.com/thomas-holmes/delivery-rl/game/controls"
)

type Scene interface {
	Update(input controls.InputEvent, deltaT uint32)
	Render(deltaT uint32)
}
