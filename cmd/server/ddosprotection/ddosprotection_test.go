package ddosprotection_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"faraway/cmd/server/ddosprotection"
	"faraway/internal/pow"
)

func TestDDOSProtection_HandleConnection(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockChallenge := pow.Challenge{Difficulty: 1, RandomData: pow.RandomData{0x02}}
		mockChallengeBytes := []byte{
			0x00, 0x01,
			0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		} // Big-endian encoding of mockChallenge
		mockNonce := pow.Nonce(3)
		mockNonceBytes := []byte{0x00, 0x00, 0x00, 0x03} // Big-endian encoding of mockNonce

		mockPOWGenerator := &POWGeneratorMock{
			GenerateChallengeFunc: func() (pow.Challenge, error) {
				return mockChallenge, nil
			},
		}

		mockPOWVerifier := &POWVerifierMock{
			VerifySolutionFunc: func(challenge pow.Challenge, nonce pow.Nonce) bool {
				assert.Equal(t, mockChallenge, challenge)
				assert.Equal(t, mockNonce, nonce)
				return true
			},
		}

		mockConnection := &ConnectionMock{
			WriteFunc: func(b []byte) (int, error) {
				assert.Equal(t, mockChallengeBytes, b)
				return len(b), nil
			},
			ReadFunc: func(p []byte) (n int, err error) {
				copy(p, mockNonceBytes)
				return len(mockNonceBytes), nil
			},
		}

		mockTCPHandler := &TCPServerHandlerMock{
			HandleConnectionFunc: func(conn net.Conn) error {
				return nil
			},
		}

		ddosProtection := ddosprotection.New(mockPOWGenerator, mockPOWVerifier, mockTCPHandler)

		err := ddosProtection.HandleConnection(mockConnection)

		require.NoError(t, err)
	})

	t.Run("generate challenge error", func(t *testing.T) {
		mockPOWGenerator := &POWGeneratorMock{
			GenerateChallengeFunc: func() (pow.Challenge, error) {
				return pow.Challenge{}, assert.AnError
			},
		}

		mockConnection := &ConnectionMock{}

		mockTCPHandler := &TCPServerHandlerMock{}

		ddosProtection := ddosprotection.New(mockPOWGenerator, nil, mockTCPHandler)

		err := ddosProtection.HandleConnection(mockConnection)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "protection")
	})

	t.Run("write challenge error", func(t *testing.T) {
		mockChallenge := pow.Challenge{Difficulty: 1, RandomData: pow.RandomData{0x02}}

		mockPOWGenerator := &POWGeneratorMock{
			GenerateChallengeFunc: func() (pow.Challenge, error) {
				return mockChallenge, nil
			},
		}

		mockConnection := &ConnectionMock{
			WriteFunc: func(b []byte) (int, error) {
				return 0, assert.AnError
			},
		}

		mockTCPHandler := &TCPServerHandlerMock{}

		ddosProtection := ddosprotection.New(mockPOWGenerator, nil, mockTCPHandler)

		err := ddosProtection.HandleConnection(mockConnection)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "write challenge")
	})

	t.Run("read nonce error", func(t *testing.T) {
		mockChallenge := pow.Challenge{Difficulty: 1, RandomData: pow.RandomData{0x02}}
		mockNonceBytes := []byte{0x00, 0x00, 0x00, 0x03} // Big-endian encoding of mockNonce

		mockPOWGenerator := &POWGeneratorMock{
			GenerateChallengeFunc: func() (pow.Challenge, error) {
				return mockChallenge, nil
			},
		}

		mockConnection := &ConnectionMock{
			WriteFunc: func(b []byte) (int, error) {
				return len(b), nil
			},
			ReadFunc: func(p []byte) (n int, err error) {
				copy(p, mockNonceBytes)
				return 0, assert.AnError
			},
		}

		mockTCPHandler := &TCPServerHandlerMock{}

		ddosProtection := ddosprotection.New(mockPOWGenerator, nil, mockTCPHandler)

		err := ddosProtection.HandleConnection(mockConnection)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "read nonce")
	})

	t.Run("invalid nonce", func(t *testing.T) {
		mockChallenge := pow.Challenge{Difficulty: 1, RandomData: pow.RandomData{0x02}}

		mockPOWGenerator := &POWGeneratorMock{
			GenerateChallengeFunc: func() (pow.Challenge, error) {
				return mockChallenge, nil
			},
		}

		mockConnection := &ConnectionMock{
			WriteFunc: func(b []byte) (int, error) {
				return len(b), nil
			},
			ReadFunc: func(p []byte) (n int, err error) {
				return len(p), nil
			},
		}

		mockPOWVerifier := &POWVerifierMock{
			VerifySolutionFunc: func(challenge pow.Challenge, nonce pow.Nonce) bool {
				return false
			},
		}

		mockTCPHandler := &TCPServerHandlerMock{}

		ddosProtection := ddosprotection.New(mockPOWGenerator, mockPOWVerifier, mockTCPHandler)

		err := ddosProtection.HandleConnection(mockConnection)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid nonce")
	})
}
