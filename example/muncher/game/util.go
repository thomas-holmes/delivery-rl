package game

import "github.com/veandco/go-sdl2/sdl"

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

var Green = sdl.Color{R: 0, G: 255, B: 0, A: 255}
var Yellow = sdl.Color{R: 255, G: 255, B: 0, A: 255}
var Orange = sdl.Color{R: 255, G: 192, B: 0, A: 255}
var Red = sdl.Color{R: 255, G: 0, B: 0, A: 255}