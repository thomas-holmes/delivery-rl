package pcg

// PCG Random Number Generation
// Developed by Melissa O'Neill <oneill@pcg-random.org>
// Paper and details at http://www.pcg-random.org
// Ported to Go by Michael Jones <michael.jones@gmail.com>

const (
	pcg32State      = 0x853c49e6748fea9b //  9600629759793949339
	pcg32Increment  = 0xda3e39cb94b95bdb // 15726070495360670683
	pcg32Multiplier = 0x5851f42d4c957f2d //  6364136223846793005
)

type PCG32 struct {
	State     uint64
	Increment uint64
}

func NewPCG32() PCG32 {
	return PCG32{pcg32State, pcg32Increment}
}

func (p *PCG32) Seed(State, sequence uint64) *PCG32 {
	p.Increment = (sequence << 1) | 1
	p.State = (State+p.Increment)*pcg32Multiplier + p.Increment
	return p
}

func (p *PCG32) Random() uint32 {
	// Advance 64-bit linear congruential generator to new State
	oldState := p.State
	p.State = oldState*pcg32Multiplier + p.Increment

	// Confuse and permute 32-bit output from old State
	xorShifted := uint32(((oldState >> 18) ^ oldState) >> 27)
	rot := uint32(oldState >> 59)
	return (xorShifted >> rot) | (xorShifted << ((-rot) & 31))
}

func (p *PCG32) Bounded(bound uint32) uint32 {
	if bound == 0 {
		return 0
	}
	threshold := -bound % bound
	for {
		r := p.Random()
		if r >= threshold {
			return r % bound
		}
	}
}

func (p *PCG32) Advance(delta uint64) *PCG32 {
	p.State = p.advanceLCG64(p.State, delta, pcg32Multiplier, p.Increment)
	return p
}

func (p *PCG32) Retreat(delta uint64) *PCG32 {
	return p.Advance(-delta)
}

func (p *PCG32) advanceLCG64(State, delta, curMult, curPlus uint64) uint64 {
	accMult := uint64(1)
	accPlus := uint64(0)
	for delta > 0 {
		if delta&1 != 0 {
			accMult *= curMult
			accPlus = accPlus*curMult + curPlus
		}
		curPlus = (curMult + 1) * curPlus
		curMult *= curMult
		delta /= 2
	}
	return accMult*State + accPlus
}
