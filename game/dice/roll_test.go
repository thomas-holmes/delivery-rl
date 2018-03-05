package dice

import (
	"fmt"
	"testing"
)

func testParse(t *testing.T, num int, sides int, bonus int) {
	var notationStr string
	if bonus == 0 {
		notationStr = fmt.Sprintf("%dd%d", num, sides)
	} else {
		notationStr = fmt.Sprintf("%dd%d+%d", num, sides, bonus)
	}

	notation, err := ParseNotation(notationStr)
	if err != nil {
		t.Error("Failed to parse", err)
	}
	if notation.Num != num {
		t.Errorf("Expected to get %d dice, instead got %d", num, notation.Num)
	}
	if notation.Sides != sides {
		t.Errorf("Expected to get %d sides, instead got %d", sides, notation.Sides)
	}
	if notation.Bonus != bonus {
		t.Errorf("Expected to get bonus of %d, instead got %d", bonus, notation.Bonus)
	}

	if bonus == 0 {
		notationStr = fmt.Sprintf("%dd%d", num, sides)
	} else {
		notationStr = fmt.Sprintf("%dd%d + %d", num, sides, bonus)
	}

	notation, err = ParseNotation(notationStr)
	if err != nil {
		t.Error("Failed to parse", err)
	}
	if notation.Num != num {
		t.Errorf("Expected to get %d dice, instead got %d", num, notation.Num)
	}
	if notation.Sides != sides {
		t.Errorf("Expected to get %d sides, instead got %d", sides, notation.Sides)
	}
	if notation.Bonus != bonus {
		t.Errorf("Expected to get bonus of %d, instead got %d", bonus, notation.Bonus)
	}
}

func TestParseNotation(t *testing.T) {
	for num := 1; num < 20; num++ {
		for sides := 2; sides <= 100; sides += 2 {
			for bonus := 0; bonus < 1; bonus++ {
				testParse(t, num, sides, bonus)
			}
		}
	}
}

func TestParseNotationRegex(t *testing.T) {
	notation, err := ParseNotation("2d4+6")
	if err != nil {
		t.Error("failed to parse", err)
	}

	if notation.Num != 2 {
		t.Error("Expected 2 dice, got", notation.Num)
	}
	if notation.Sides != 4 {
		t.Error("Expected 4 sides, got", notation.Sides)
	}
	if notation.Bonus != 6 {
		t.Error("Expected 6 bonus, got", notation.Bonus)
	}
}
