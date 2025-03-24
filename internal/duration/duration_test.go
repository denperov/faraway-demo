package duration_test

import (
	"log/slog"
	"testing"
	"time"

	"faraway/internal/duration"
)

func TestLog(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	defer duration.Log("msg_1")()

	time.Sleep(100 * time.Millisecond)
}
