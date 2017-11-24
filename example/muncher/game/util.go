package game

import (
	"log"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

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

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -1 * a
}

func timeMe(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%v %s", name, elapsed)

}

func distance(x0 int, y0 int, x1 int, y1 int) int {
	x := x1 - x0
	y := y1 - y0
	return int(math.Sqrt(float64(x*x) + float64(y*y)))
}
