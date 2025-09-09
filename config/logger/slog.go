package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

type ColorTextHandler struct {
	slog.Handler
}

func (h *ColorTextHandler) Handle(ctx context.Context, r slog.Record) error {
	var levelColor *color.Color
	switch r.Level {
	case slog.LevelInfo:
		levelColor = color.New(color.FgGreen)
	case slog.LevelWarn:
		levelColor = color.New(color.FgYellow)
	case slog.LevelError:
		levelColor = color.New(color.FgRed)
	default:
		levelColor = color.New(color.FgWhite)
	}
	levelStr := r.Level.String()
	msg := fmt.Sprintf("%s %s", levelStr, r.Message)
	levelColor.Fprintf(os.Stdout, "%s\n", msg)
	return nil
}

func NewSlog() *slog.Logger {
	baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	colorHandler := &ColorTextHandler{Handler: baseHandler}
	return slog.New(colorHandler)
}
