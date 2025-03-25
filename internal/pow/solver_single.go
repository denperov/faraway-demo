package pow

import (
	"context"

	"faraway/internal/duration"
)

type SolverSingle struct{}

func NewSolverSingle() *SolverSingle {
	return &SolverSingle{}
}

func (c *SolverSingle) SolveChallenge(ctx context.Context, challenge Challenge) (Nonce, error) {
	defer duration.Log("solve challenge", "difficulty", challenge.Difficulty)()

	for nonce := NonceMin; nonce < NonceMax; nonce++ {
		if nonce%1024 == 0 { // check for cancellation every 1024 iterations to avoid overhead
			if err := ctx.Err(); err != nil {
				return 0, err
			}
		}

		if CalculateSolution(challenge.RandomData, nonce).Difficulty() >= challenge.Difficulty {
			return nonce, nil
		}
	}

	return 0, ErrorSolutionNotFound{Difficulty: challenge.Difficulty}
}
