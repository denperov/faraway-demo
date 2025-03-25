package quoteclient

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"

	"faraway/internal/pow"
)

type POWSolver interface {
	SolveChallenge(context.Context, pow.Challenge) (pow.Nonce, error)
}

type QuoteClient struct {
	powSolver     POWSolver
	serverAddress string
}

func New(
	powSolver POWSolver,
	serverAddress string,
) *QuoteClient {
	return &QuoteClient{
		powSolver:     powSolver,
		serverAddress: serverAddress,
	}
}

func (c *QuoteClient) GetQuote(ctx context.Context) (string, error) {
	slog.Debug("dial", "address", c.serverAddress)

	conn, err := net.Dial("tcp", c.serverAddress)
	if err != nil {
		return "", fmt.Errorf("dial tcp: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("close connection", "error", err)
		}
	}()

	slog.Debug("read challenge")

	var challenge pow.Challenge

	if err := binary.Read(conn, binary.BigEndian, &challenge); err != nil {
		return "", fmt.Errorf("read challenge: %w", err)
	}

	nonce, err := c.powSolver.SolveChallenge(ctx, challenge)
	if err != nil {
		return "", fmt.Errorf("solve challenge: %w", err)
	}

	if err := binary.Write(conn, binary.BigEndian, nonce); err != nil {
		return "", fmt.Errorf("write nonce: %w", err)
	}

	slog.Debug("read quote")

	quote, err := io.ReadAll(conn)
	if err != nil {
		return "", fmt.Errorf("read quote: %w", err)
	}

	return string(quote), nil
}
