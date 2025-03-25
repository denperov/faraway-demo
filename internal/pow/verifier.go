package pow

type Verifier struct {
}

func NewVerifier() *Verifier {
	return &Verifier{}
}

func (s *Verifier) VerifySolution(challenge Challenge, nonce Nonce) bool {
	return CalculateSolution(challenge.RandomData, nonce).Difficulty() >= challenge.Difficulty
}
