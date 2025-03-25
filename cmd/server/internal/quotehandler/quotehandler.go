package quotehandler

import (
	"fmt"
	"net"

	"faraway/cmd/server/internal/tcpserver"
)

//go:generate go tool github.com/matryer/moq -pkg quotehandler_test -out mocks_test.go . Connection QuoteStorage

// Connection is an alias for net.Conn to generate a mock.
type Connection net.Conn

// QuoteStorage is an interface for getting a quote.
type QuoteStorage interface {
	GetQuote() (string, error)
}

var _ tcpserver.Handler = (*QuoteHandler)(nil)

// QuoteHandler is a handler for quotes.
type QuoteHandler struct {
	quoteStorage QuoteStorage
}

// New creates a new QuoteHandler.
func New(
	quoteStorage QuoteStorage,
) *QuoteHandler {
	return &QuoteHandler{
		quoteStorage: quoteStorage,
	}
}

// HandleConnection retrieves a quote from the storage and writes it to the connection.
func (s *QuoteHandler) HandleConnection(conn net.Conn) error {
	quote, err := s.quoteStorage.GetQuote()
	if err != nil {
		return fmt.Errorf("get quote: %w", err)
	}

	if _, err := conn.Write([]byte(quote)); err != nil {
		return fmt.Errorf("write quote: %w", err)
	}

	return nil
}
