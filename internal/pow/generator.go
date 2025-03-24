package pow

import (
	"crypto/rand"
	"fmt"
)

type Generator struct {
	difficulty Difficulty
}

func NewGenerator(difficulty Difficulty) *Generator {
	return &Generator{
		difficulty: difficulty,
	}
}

func (s *Generator) GenerateChallenge() (Challenge, error) {
	randomData, err := generateRandomData()
	if err != nil {
		return Challenge{}, err
	}

	return Challenge{
		Difficulty: s.difficulty,
		RandomData: randomData,
	}, nil
}

func generateRandomData() (RandomData, error) {
	var randomData RandomData

	if _, err := rand.Read(randomData[:]); err != nil {
		return RandomData{}, fmt.Errorf("generate random data: %w", err)
	}

	return randomData, nil
}
