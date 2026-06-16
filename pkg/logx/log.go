package logx

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type LogFormat string
type LogLevel string

const (
	FormatText LogFormat = "text"
	FormatJSON LogFormat = "json"

	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type Options struct {
	Level     LogLevel
	Format    LogFormat
	AddSource bool
	Output    io.Writer
}

type Logger struct {
	logger *slog.Logger
	level  *slog.LevelVar
}

var (
	mu  sync.RWMutex
	std *Logger
)

type ctxKey struct{}

func New(opts Options) *Logger {
	levelVar := new(slog.LevelVar)
	levelVar.Set(parseLevel(opts.Level))

	output := opts.Output
	if output == nil {
		output = os.Stdout
	}

	handlerOpts := &slog.HandlerOptions{
		Level:     levelVar,
		AddSource: opts.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				if a.Value.Kind() == slog.KindTime {
					return slog.String(slog.TimeKey, a.Value.Time().Format("2006-01-02 15:04:05"))
				}
				return a

			case slog.SourceKey:
				src, ok := a.Value.Any().(*slog.Source)
				if !ok || src == nil {
					return a
				}
				return slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", filepath.Base(src.File), src.Line))

			default:
				return a
			}
		},
	}

	var handler slog.Handler
	switch opts.Format {
	case FormatJSON:
		handler = slog.NewJSONHandler(output, handlerOpts)
	default:
		handler = slog.NewTextHandler(output, handlerOpts)
	}

	base := slog.New(handler)
	return &Logger{
		logger: base,
		level:  levelVar,
	}
}

func Init(opts Options) {
	lg := New(opts)

	mu.Lock()
	std = lg
	mu.Unlock()

	slog.SetDefault(lg.logger)
}

func SetDefault(lg *Logger) {
	if lg == nil {
		return
	}

	mu.Lock()
	std = lg
	mu.Unlock()

	slog.SetDefault(lg.logger)
}

func L() *Logger {
	mu.RLock()
	if std != nil {
		lg := std
		mu.RUnlock()
		return lg
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	if std == nil {
		std = New(Options{
			Level:  LevelInfo,
			Format: FormatText,
			Output: os.Stdout,
		})
		slog.SetDefault(std.logger)
	}

	return std
}

func parseLevel(level LogLevel) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (g *Logger) clone(l *slog.Logger) *Logger {
	if g == nil {
		return &Logger{
			logger: l,
			level:  nil,
		}
	}
	return &Logger{
		logger: l,
		level:  g.level,
	}
}

func (g *Logger) base() *Logger {
	if g == nil || g.logger == nil {
		return L()
	}
	return g
}

func (g *Logger) Slog() *slog.Logger {
	return g.base().logger
}

func (g *Logger) With(args ...any) *Logger {
	bg := g.base()
	return bg.clone(bg.logger.With(args...))
}

func (g *Logger) WithGroup(name string) *Logger {
	bg := g.base()
	return bg.clone(bg.logger.WithGroup(name))
}

func (g *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	bg := g.base()
	return bg.logger.Enabled(ctx, level)
}

func (g *Logger) SetLevel(level slog.Level) {
	bg := g.base()
	if bg.level != nil {
		bg.level.Set(level)
	}
}

func (g *Logger) Level() slog.Level {
	bg := g.base()
	if bg.level == nil {
		return slog.LevelInfo
	}
	return bg.level.Level()
}

func (g *Logger) Debug(msg string, args ...any) {
	g.base().logger.Debug(msg, args...)
}

func (g *Logger) Info(msg string, args ...any) {
	g.base().logger.Info(msg, args...)
}

func (g *Logger) Warn(msg string, args ...any) {
	g.base().logger.Warn(msg, args...)
}

func (g *Logger) Error(msg string, args ...any) {
	g.base().logger.Error(msg, args...)
}

func (g *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	g.base().logger.DebugContext(ctx, msg, args...)
}

func (g *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	g.base().logger.InfoContext(ctx, msg, args...)
}

func (g *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	g.base().logger.WarnContext(ctx, msg, args...)
}

func (g *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	g.base().logger.ErrorContext(ctx, msg, args...)
}

func (g *Logger) Debugf(format string, args ...any) {
	g.base().logger.Debug(fmt.Sprintf(format, args...))
}

func (g *Logger) Infof(format string, args ...any) {
	g.base().logger.Info(fmt.Sprintf(format, args...))
}

func (g *Logger) Warnf(format string, args ...any) {
	g.base().logger.Warn(fmt.Sprintf(format, args...))
}

func (g *Logger) Errorf(format string, args ...any) {
	g.base().logger.Error(fmt.Sprintf(format, args...))
}

// 全局快捷方法
func With(args ...any) *Logger      { return L().With(args...) }
func WithGroup(name string) *Logger { return L().WithGroup(name) }
func Slog() *slog.Logger            { return L().Slog() }
func SetLevel(level slog.Level)     { L().SetLevel(level) }
func Level() slog.Level             { return L().Level() }
func Debug(msg string, args ...any) { L().Debug(msg, args...) }
func Info(msg string, args ...any)  { L().Info(msg, args...) }
func Warn(msg string, args ...any)  { L().Warn(msg, args...) }
func Error(msg string, args ...any) { L().Error(msg, args...) }
func DebugContext(ctx context.Context, msg string, args ...any) {
	L().DebugContext(ctx, msg, args...)
}
func InfoContext(ctx context.Context, msg string, args ...any) {
	L().InfoContext(ctx, msg, args...)
}
func WarnContext(ctx context.Context, msg string, args ...any) {
	L().WarnContext(ctx, msg, args...)
}
func ErrorContext(ctx context.Context, msg string, args ...any) {
	L().ErrorContext(ctx, msg, args...)
}
func Debugf(format string, args ...any) { L().Debugf(format, args...) }
func Infof(format string, args ...any)  { L().Infof(format, args...) }
func Warnf(format string, args ...any)  { L().Warnf(format, args...) }
func Errorf(format string, args ...any) { L().Errorf(format, args...) }

func NewContext(ctx context.Context, lg *Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if lg == nil {
		lg = L()
	}
	return context.WithValue(ctx, ctxKey{}, lg)
}

func FromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return L()
	}
	if lg, ok := ctx.Value(ctxKey{}).(*Logger); ok && lg != nil {
		return lg
	}
	return L()
}

type Fields map[string]any

type fieldsKey struct{}

func ContextWithFields(ctx context.Context, fields Fields) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(fields) == 0 {
		return ctx
	}
	return context.WithValue(ctx, fieldsKey{}, fields)
}

func FieldsFromContext(ctx context.Context) Fields {
	if ctx == nil {
		return nil
	}
	if v, ok := ctx.Value(fieldsKey{}).(Fields); ok {
		return v
	}
	return nil
}
