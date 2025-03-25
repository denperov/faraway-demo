package quotestorage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"faraway/cmd/server/internal/quotestorage"
)

func TestQuoteStorage_GetQuote(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockQuote := "quote_1"

		quoteStorage := quotestorage.New()
		quoteStorage.AddQuote(mockQuote)

		quote, err := quoteStorage.GetQuote()

		require.NoError(t, err)
		assert.Equal(t, mockQuote, quote)
	})

	t.Run("empty storage", func(t *testing.T) {
		quoteStorage := quotestorage.New()

		_, err := quoteStorage.GetQuote()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no quotes available")
	})
}
