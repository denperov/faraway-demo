package pow

import (
	"encoding"
	"fmt"
	"math/bits"
	"strconv"
)

type Nonce uint32

type Difficulty uint16

var _ encoding.TextUnmarshaler = (*Difficulty)(nil)

func (d *Difficulty) UnmarshalText(text []byte) error {
	value, err := strconv.ParseUint(string(text), 10, 16)
	if err != nil {
		return fmt.Errorf("parse difficulty: %w", err)
	}
	*d = Difficulty(value)
	return nil
}

type RandomData [32]byte

type Challenge struct {
	Difficulty Difficulty
	RandomData RandomData
}

type Solution [32]byte

func (s Solution) Difficulty() Difficulty {
	for i, b := range s {
		if b != 0 {
			return Difficulty(i*8 + bits.LeadingZeros8(b))
		}
	}

	return Difficulty(len(s) * 8)
}
