package quotestorage

import (
	"fmt"
	"io/ioutil"
	"log/slog"
	"math/rand"
	"os"

	"gopkg.in/yaml.v2"
)

// Quote represents a single quote entry in the YAML file.
type Quote struct {
	Text string `yaml:"text"`
}

// quotesFile represents the overall YAML structure.
type quotesFile struct {
	Quotes []Quote `yaml:"quotes"`
}

type QuoteStorage struct {
	quotes []string
}

func New() *QuoteStorage {
	return &QuoteStorage{}
}

// Load reads the YAML file, parses it, and stores the quotes.
func (s *QuoteStorage) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("close file", "error", err)
		}
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var qf quotesFile
	if err := yaml.Unmarshal(data, &qf); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	for _, q := range qf.Quotes {
		if q.Text != "" {
			s.quotes = append(s.quotes, q.Text)
		}
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
