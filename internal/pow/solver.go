package pow

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"faraway/internal/duration"
)

type ErrorSolutionNotFound struct {
	Difficulty Difficulty
}

func (e ErrorSolutionNotFound) Error() string {
	return fmt.Sprintf("solution not found for difficulty %d", e.Difficulty)
}

type Solver struct{}

func NewSolver() *Solver {
	return &Solver{}
}

func (c *Solver) SolveChallenge(challenge Challenge) (Nonce, error) {
	defer duration.Log("solve challenge", "difficulty", challenge.Difficulty)()

	var buffer [32 + 4]byte
	copy(buffer[:32], challenge.RandomData[:])

	for nonce := MinNonce; nonce < MaxNonce; nonce++ {
		binary.BigEndian.PutUint32(buffer[32:], uint32(nonce))

		firstHash := sha256.Sum256(buffer[:])
		solution := Solution(sha256.Sum256(firstHash[:]))

		if solution.Difficulty() >= challenge.Difficulty {
			return nonce, nil
		}
	}

	return 0, ErrorSolutionNotFound{Difficulty: challenge.Difficulty}
}
