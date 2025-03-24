package pow

import (
	"crypto/sha256"
	"encoding/binary"
)

type Verifier struct {
}

func NewVerifier() *Verifier {
	return &Verifier{}
}

func (s *Verifier) VerifySolution(challenge Challenge, nonce Nonce) bool {
	var data [32 + 4]byte
	copy(data[:32], challenge.RandomData[:])
	binary.BigEndian.PutUint32(data[32:], uint32(nonce))

	firstHash := sha256.Sum256(data[:])
	solution := Solution(sha256.Sum256(firstHash[:]))

	return solution.Difficulty() >= challenge.Difficulty
}
