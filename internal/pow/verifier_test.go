package pow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"faraway/internal/pow"
)

func TestVerifier_VerifySolution(t *testing.T) {
	randomData := [32]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}

	challenge := pow.Challenge{
		Difficulty: 10,
		RandomData: randomData,
	}

	verifier := pow.NewVerifier()

	t.Run("valid solution", func(t *testing.T) {
		nonce := pow.Nonce(2147946924)

		result := verifier.VerifySolution(challenge, nonce)

		assert.True(t, result)
	})

	t.Run("invalid solution", func(t *testing.T) {
		nonce := pow.Nonce(1234)

		result := verifier.VerifySolution(challenge, nonce)

		assert.False(t, result)
	})
}
