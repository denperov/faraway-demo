package ddosprotection

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"

	"faraway/cmd/server/internal/tcpserver"
	"faraway/internal/pow"
)

//go:generate go tool github.com/matryer/moq -pkg ddosprotection_test -out mocks_test.go . Connection TCPServerHandler POWGenerator POWVerifier

// Connection is an alias for net.Conn to generate a mock.
type Connection net.Conn

// TCPServerHandler is an alias for tcpserver.Handler to generate a mock.
type TCPServerHandler tcpserver.Handler

var _ tcpserver.Handler = (*DDOSProtection)(nil)

type POWGenerator interface {
	GenerateChallenge() (pow.Challenge, error)
}

type POWVerifier interface {
	VerifySolution(pow.Challenge, pow.Nonce) bool
}

type DDOSProtection struct {
	powGenerator POWGenerator
	powVerifier  POWVerifier
	handler      tcpserver.Handler
}

func New(
	powGenerator POWGenerator,
	powVerifier POWVerifier,
	handler tcpserver.Handler,
) *DDOSProtection {
	return &DDOSProtection{
		powGenerator: powGenerator,
		powVerifier:  powVerifier,
		handler:      handler,
	}
}

func (s *DDOSProtection) HandleConnection(conn net.Conn) error {
	if err := s.protection(conn); err != nil {
		return fmt.Errorf("protection: %w", err)
	}

	return s.handler.HandleConnection(conn)
}

func (s *DDOSProtection) protection(conn net.Conn) error {
	slog.Debug("start protection")

	challenge, err := s.powGenerator.GenerateChallenge()
	if err != nil {
		return fmt.Errorf("generate challenge: %w", err)
	}

	if err := binary.Write(conn, binary.BigEndian, challenge); err != nil {
		return fmt.Errorf("write challenge: %w", err)
	}

	var nonce pow.Nonce

	if err := binary.Read(conn, binary.BigEndian, &nonce); err != nil {
		return fmt.Errorf("read nonce: %w", err)
	}

	if !s.powVerifier.VerifySolution(challenge, nonce) {
		return fmt.Errorf("invalid nonce")
	}

	slog.Debug("protection passed", "difficulty", challenge.Difficulty)

	return nil
}
