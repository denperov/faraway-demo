package main

import (
	"context"
	"fmt"
	"log/slog"

	"faraway/cmd/server/ddosprotection"
	"faraway/cmd/server/quotehandler"
	"faraway/cmd/server/quotestorage"
	"faraway/cmd/server/tcpserver"
	"faraway/internal/app"
	"faraway/internal/pow"
)

func main() { app.Run(service) }

type Config struct {
	ListenAddress        string         `envconfig:"LISTEN_ADDRESS" default:"127.0.0.1:1337"`
	QuotesFile           string         `envconfig:"QUOTES_FILE" default:"quotes.txt"`
	ProtectionDifficulty pow.Difficulty `envconfig:"PROTECTION_DIFFICULTY" default:"20"`
}

func service(ctx context.Context, cfg Config) error {
	slog.Debug("service started", "listen_address", cfg.ListenAddress, "protection_difficulty", cfg.ProtectionDifficulty)

	powGenerator := pow.NewGenerator(cfg.ProtectionDifficulty)
	powVerifier := pow.NewVerifier()

	quoteStorage := quotestorage.New()
	if err := quoteStorage.Load(cfg.QuotesFile); err != nil {
		return fmt.Errorf("load quotes: %w", err)
	}

	quoteHandler := quotehandler.New(
		quoteStorage,
	)

	ddosProtection := ddosprotection.New(
		powGenerator,
		powVerifier,
		quoteHandler,
	)

	tcpServer := tcpserver.New(
		cfg.ListenAddress,
		ddosProtection,
	)

	return tcpServer.Run(ctx)
}
