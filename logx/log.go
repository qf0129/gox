package logx

import (
	"context"
	"log/slog"
	"os"
)

var logger = slog.Default()

func GetLogger() *slog.Logger {
	return logger
}

func SetJsonLogger() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
}

func SetTextLogger() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
}

func NewTextFileLogger(filepath string) *slog.Logger {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		logger.Error("CreateLogFileError", "err", err)
	}
	return slog.New(slog.NewTextHandler(f, nil))
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger.Log(ctx, level, msg, args...)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logger.LogAttrs(ctx, level, msg, attrs...)
}
