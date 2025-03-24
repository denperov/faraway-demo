package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"faraway/cmd/client/quoteclient"
	"faraway/internal/app"
	"faraway/internal/delay"
	"faraway/internal/pow"
)

func main() { app.Run(client) }

type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:"127.0.0.1:1337"`
}

func client(ctx context.Context, cfg Config) error {
	powSolver := pow.NewSolver()

	quoteClient := quoteclient.New(
		powSolver,
		cfg.ServerAddress,
	)

	for {
		quote, err := quoteClient.GetQuote()
		if err != nil {
			slog.Error("get quote failed", "error", err)
		}

		fmt.Println(quote)

		if err := delay.For(ctx, 1*time.Second); err != nil {
			return nil // context finished, normal shutdown
		}
	}
}
