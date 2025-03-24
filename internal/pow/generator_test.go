package pow_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"faraway/internal/pow"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_GenerateChallenge(t *testing.T) {
	difficulty := pow.Difficulty(10)
	generator := pow.NewGenerator(difficulty)

	challenge, err := generator.GenerateChallenge()

	require.NoError(t, err)
	assert.Equal(t, difficulty, challenge.Difficulty)
	assert.NotEqual(t, pow.RandomData{}, challenge.RandomData)
}
