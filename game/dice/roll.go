package dice

import (
	"errors"
	"strconv"
	"strings"

	"github.com/MichaelTJones/pcg"
)

/*
 1 - 1d4  (2.5)
 2 - 1d6  (3.5)
 3 - 1d8  (4.5)
 4 - 2d4  (5)
 5 - 1d10 (5.5)
 6 - 1d12 (6.5)
 7 - 3d4  (7.5)
 8 - 2d6  (7)
 9 - 2d8  (9)
10 - 3d6  (10.5)
11 - 2d10 (11)
12 - 2d12 (13)
13 - 3d8  (13.5)
14 - 3d10 (16.5)
15 - 3d12 (19.5)
*/

var defaultRoller Roller

type Roller struct {
	rng *pcg.PCG64
}

// Roll rolls n dice with y sides. To simulate a roll of 4d8 you
// would call Roll(4, 8)
func (r Roller) Roll(num int, sides int) int {
	if r.rng == nil {
		panic("Can't roll dice without randomness")
	}

	total := 0
	for rolls := 0; rolls < num; rolls++ {
		total += (int(r.rng.Bounded(uint64(sides))) + 1)
	}

	return total
}

func (r Roller) RollDice(diceNotation string) (int, error) {
	parts := strings.Split(strings.TrimSpace(diceNotation), "d")
	if len(parts) != 2 {
		return 0, errors.New("Could not parse dice notation. Pass in form NdY")
	}
	numStr, sidesStr := parts[0], parts[1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}
	sides, err := strconv.Atoi(sidesStr)
	if err != nil {
		return 0, err
	}

	return r.Roll(num, sides), nil
}

// NewRoller constructs a roller with the provided PCG64 rng
func NewRoller(rng *pcg.PCG64) Roller {
	return Roller{rng: rng}
}

// SetDefaultRandomness set the rng for the default Roller
func SetDefaultRandomness(rng *pcg.PCG64) {
	defaultRoller = NewRoller(rng)
}

// Roll roles dice using the default roller
func Roll(num int, sides int) int {
	return defaultRoller.Roll(num, sides)
}

func RollDice(diceNotation string) (int, error) {
	return defaultRoller.RollDice(diceNotation)
}
