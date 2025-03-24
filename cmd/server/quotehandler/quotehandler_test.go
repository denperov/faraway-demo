package quotehandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"faraway/cmd/server/quotehandler"
)

func TestQuoteServer_HandleConnection(t *testing.T) {
	mockQuote := "quote_1"

	t.Run("success", func(t *testing.T) {
		mockQuoteStorage := &QuoteStorageMock{
			GetQuoteFunc: func() (string, error) {
				return mockQuote, nil
			},
		}

		var writtenData []byte

		mockConnection := &ConnectionMock{
			WriteFunc: func(p []byte) (n int, err error) {
				writtenData = append([]byte(nil), p...)
				return len(p), nil
			},
		}

		quoteHandler := quotehandler.New(mockQuoteStorage)

		err := quoteHandler.HandleConnection(mockConnection)

		require.NoError(t, err)
		assert.Equal(t, mockQuote, string(writtenData))
	})

	t.Run("storage error", func(t *testing.T) {
		mockQuoteStorage := &QuoteStorageMock{
			GetQuoteFunc: func() (string, error) {
				return "", assert.AnError
			},
		}

		mockConnection := &ConnectionMock{}

		quoteHandler := quotehandler.New(mockQuoteStorage)

		err := quoteHandler.HandleConnection(mockConnection)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "get quote")
	})

	t.Run("write error", func(t *testing.T) {
		mockQuoteStorage := &QuoteStorageMock{
			GetQuoteFunc: func() (string, error) {
				return mockQuote, nil
			},
		}

		mockConnection := &ConnectionMock{
			WriteFunc: func(p []byte) (n int, err error) {
				return 0, assert.AnError
			},
		}

		quoteHandler := quotehandler.New(mockQuoteStorage)

		err := quoteHandler.HandleConnection(mockConnection)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "write quote")
	})
}
