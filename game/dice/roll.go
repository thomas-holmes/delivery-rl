package dice

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

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

type Notation struct {
	Num   int
	Sides int
	Bonus int
}

func (n Notation) String() string {
	if n.Bonus > 0 {
		return fmt.Sprintf("%dd%d+%d", n.Num, n.Sides, n.Bonus)
	} else {
		return fmt.Sprintf("%dd%d", n.Num, n.Sides)
	}
}

// Roll rolls n dice with y sides. To simulate a roll of 4d8 you
// would call Roll(4, 8)
func (r Roller) Roll(notation Notation) int {
	if r.rng == nil {
		panic("Can't roll dice without randomness")
	}

	num, sides := notation.Num, notation.Sides

	total := 0
	for rolls := 0; rolls < num; rolls++ {
		total += (int(r.rng.Bounded(uint64(sides))) + 1)
	}

	return total
}

var diceRegex = regexp.MustCompile(`(?P<num>\d+)d(?P<sides>\d+)(\s*\+\s*(?P<bonus>\d+))?`)

func ParseNotation(notationStr string) (Notation, error) {
	var notation Notation
	if !diceRegex.MatchString(notationStr) {
		return Notation{}, errors.New("Not a valid dice notation")
	}
	matches := diceRegex.FindStringSubmatch(notationStr)
	if matches == nil {
		return notation, errors.New("Didn't match properly")
	}
	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return notation, err
	}

	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return notation, err
	}

	var bonus int
	if matches[4] != "" {
		bonus, err = strconv.Atoi(matches[4])
		if err != nil {
			return notation, err
		}
	}
	notation.Num = num
	notation.Sides = sides
	notation.Bonus = bonus

	return notation, nil
}

// RollDice roll dice using dice notation
func (r Roller) RollDice(diceNotation string) (int, error) {
	notation, err := ParseNotation(diceNotation)
	if err != nil {
		return 0, err
	}

	return r.Roll(notation), nil
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
func Roll(notation Notation) int {
	return defaultRoller.Roll(notation)
}

// RollDice rolls dice using the default roller
func RollDice(diceNotation string) (int, error) {
	return defaultRoller.RollDice(diceNotation)
}

func (i *Notation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var notationStr string
	if err := unmarshal(&notationStr); err != nil {
		return err
	}

	notation, err := ParseNotation(notationStr)
	if err != nil {
		return err
	}

	*i = notation

	return nil
}
