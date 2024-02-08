package custom_logger

import (
	"context"
	"fmt"
	"log/slog"
)

type WrapSlogWriter struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *WrapSlogWriter {
	return &WrapSlogWriter{logger: logger}
}

func (w *WrapSlogWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}

func (w *WrapSlogWriter) Printf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	w.logger.Log(context.Background(), slog.LevelError, msg)
}
