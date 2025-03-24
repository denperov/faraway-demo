package duration

import (
	"log/slog"
	"time"
)

func Log(msg string, args ...any) func() {
	start := time.Now()
	return func() {
		slog.Debug(msg, append(args, "duration", time.Since(start))...)
	}
}
