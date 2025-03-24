package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/kelseyhightower/envconfig"
)

type ExitCode int

const (
	ExitCodeSuccess ExitCode = 0
	ExitCodeError   ExitCode = 1
)

func Run[T any](appMain func(ctx context.Context, cfg T) error) {
	// Do not use defer here!!!
	if exitCode := run(appMain); exitCode != 0 {
		os.Exit(int(exitCode))
	}
}

func run[T any](appMain func(ctx context.Context, cfg T) error) (exitCode ExitCode) {
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt)
	defer done()

	defer func() {
		if r := recover(); r != nil {
			slog.Error("app panic", "error", fmt.Errorf("%v", r))
			exitCode = ExitCodeError
		}
	}()

	var cfg T

	if err := readConfigEnv(&cfg); err != nil {
		slog.Error("read env failed", "error", err)
		return ExitCodeError
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)

	if err := appMain(ctx, cfg); err != nil {
		slog.Error("app failed", "error", err)
		return ExitCodeError
	}

	return ExitCodeSuccess
}

func readConfigEnv[T any](cfg *T) error {
	if err := envconfig.Process("", cfg); err != nil {
		return fmt.Errorf("process env: %w", err)
	}

	return nil
}
