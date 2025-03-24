package quotestorage

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
)

type QuoteStorage struct {
	quotes []string
}

func New() *QuoteStorage {
	return &QuoteStorage{}
}

func (s *QuoteStorage) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("close file: %v", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			s.quotes = append(s.quotes, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	return nil
}

func (s *QuoteStorage) GetQuote() (string, error) {
	n := len(s.quotes)
	if n == 0 {
		return "", fmt.Errorf("no quotes available")
	}

	return s.quotes[rand.Intn(n)], nil
}

func (s *QuoteStorage) AddQuote(quote string) {
	s.quotes = append(s.quotes, quote)
}
