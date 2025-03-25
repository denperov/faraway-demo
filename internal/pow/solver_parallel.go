package pow

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"

	"faraway/internal/duration"
)

type SolverParallel struct {
	parallelism int
}

func NewSolverParallel(parallelism int) *SolverParallel {
	return &SolverParallel{
		parallelism: parallelism,
	}
}

type pseudoErrorSolutionFound struct {
	Nonce Nonce
}

func (e pseudoErrorSolutionFound) Error() string { return "solution found" }

func isResultFound(err error) (Nonce, bool) {
	var e pseudoErrorSolutionFound
	if errors.As(err, &e) {
		return e.Nonce, true
	}

	return 0, false
}

func (c *SolverParallel) SolveChallenge(ctx context.Context, challenge Challenge) (Nonce, error) {
	defer duration.Log("solve challenge", "difficulty", challenge.Difficulty)()

	portion := Nonce((int(NonceMax) + 1) / c.parallelism)

	g, ctx := errgroup.WithContext(ctx)

	for i, start := 0, NonceMin; i < c.parallelism; i, start = i+1, start+portion {
		g.Go(func() error {
			for nonce := start; nonce < start+portion; nonce++ {
				if nonce%1024 == 0 { // check for cancellation every 1024 iterations to avoid overhead
					if err := ctx.Err(); err != nil {
						return err
					}
				}

				if CalculateSolution(challenge.RandomData, nonce).Difficulty() >= challenge.Difficulty {
					return pseudoErrorSolutionFound{Nonce: nonce}
				}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if nonce, ok := isResultFound(err); ok {
			return nonce, nil
		}

		return 0, err
	}

	return 0, ErrorSolutionNotFound{Difficulty: challenge.Difficulty}
}
