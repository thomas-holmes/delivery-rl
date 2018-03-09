package controls

import (
	"fmt"
	"testing"

	"github.com/veandco/go-sdl2/sdl"
)

func TestVerifyNoDuplicateMappings(t *testing.T) {
	controlSeen := make(map[sdl.Keycode]int)
	shiftSeen := make(map[sdl.Keycode]int)
	seen := make(map[sdl.Keycode]int)

	for _, m := range AllMappings {
		switch {
		case m.Shift:
			for _, k := range m.SdlKeys {
				_, ok := shiftSeen[k]
				if !ok {
					shiftSeen[k] = 1
				} else {
					panic(fmt.Sprintf("Found duplicate %v mapping %+v", k, m))
				}
			}
		case m.Control:
			for _, k := range m.SdlKeys {
				_, ok := controlSeen[k]
				if !ok {
					controlSeen[k] = 1
				} else {
					panic(fmt.Sprintf("Found duplicate %v mapping %+v", k, m))
				}
			}
		default:
			for _, k := range m.SdlKeys {
				_, ok := seen[k]
				if !ok {
					seen[k] = 1
				} else {
					panic(fmt.Sprintf("Found duplicate %v mapping %+v", k, m))
				}
			}
		}
	}
}

func TestVerifyMappingActionCounts(t *testing.T) {
	if len(AllMappings) != int(FinalUnused-1) {
		panic(fmt.Sprintf("Expected to find %d Actions but instead have %d", FinalUnused-1, len(AllMappings)))
	}
}
