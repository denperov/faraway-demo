package pow

import "fmt"

type Solver interface {
	SolveChallenge(challenge Challenge) (Nonce, error)
}

// NewSolver is factory function for creating a new Solver.
func NewSolver(parallelism int) Solver {
	// TODO: Select the appropriate solver based on CPU characteristics.
	switch parallelism {
	case 1:
		return NewSolverSingle()
	default:
		return NewSolverParallel(parallelism)
	}
}

type ErrorSolutionNotFound struct {
	Difficulty Difficulty
}

func (e ErrorSolutionNotFound) Error() string {
	return fmt.Sprintf("solution not found for difficulty %d", e.Difficulty)
}
