package pow

import (
	"crypto/sha256"
	"encoding"
	"fmt"
	"math"
	"math/bits"
	"strconv"
)

type Nonce uint32

const NonceMin = Nonce(0)
const NonceMax = Nonce(math.MaxUint32)
const NonceSize = 4

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

const RandomDataSize = 32

type RandomData [RandomDataSize]byte

type Challenge struct {
	Difficulty Difficulty
	RandomData RandomData
}

const SolutionSize = 32

type Solution [SolutionSize]byte

func CalculateSolution(randomData RandomData, nonce Nonce) Solution {
	buffer := [NonceSize + RandomDataSize]byte{
		byte(nonce >> 24),
		byte(nonce >> 16),
		byte(nonce >> 8),
		byte(nonce >> 0),
	} // copy nonce to buffer in big-endian order
	copy(buffer[NonceSize:NonceSize+RandomDataSize], randomData[:])

	hash := sha256.Sum256(buffer[:])
	hash = sha256.Sum256(hash[:])

	return hash
}

const bitsPerByte = 8

func (s Solution) Difficulty() Difficulty {
	for i, b := range s {
		if b != 0 {
			return Difficulty(i*bitsPerByte + bits.LeadingZeros8(b))
		}
	}

	return Difficulty(len(s) * bitsPerByte)
}
