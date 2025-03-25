package pow

import (
	"faraway/internal/duration"
)

type SolverSingle struct{}

func NewSolverSingle() *SolverSingle {
	return &SolverSingle{}
}

func (c *SolverSingle) SolveChallenge(challenge Challenge) (Nonce, error) {
	defer duration.Log("solve challenge", "difficulty", challenge.Difficulty)()

	for nonce := NonceMin; nonce < NonceMax; nonce++ {
		if CalculateSolution(challenge.RandomData, nonce).Difficulty() >= challenge.Difficulty {
			return nonce, nil
		}
	}

	return 0, ErrorSolutionNotFound{Difficulty: challenge.Difficulty}
}
