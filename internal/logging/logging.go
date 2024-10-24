package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/upassed/upassed-account-service/internal/config"
)

const (
	timeFormat = "[Mon Jan 2 2006 15:04:05]"
)

const (
	reset = "\033[0m"

	cyan        = 36
	lightGray   = 37
	darkGray    = 90
	lightRed    = 91
	lightYellow = 93
	white       = 97
)

func New(envType config.EnvType) *slog.Logger {
	var log *slog.Logger

	switch envType {
	case config.EnvLocal:
		log = slog.New(newCustomFormattedJsonHandler(&slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case config.EnvTesting:
		log = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}

	return log
}

type customFormattedJsonHandler struct {
	handler slog.Handler
	buffer  *bytes.Buffer
	mutex   *sync.Mutex
}

func newCustomFormattedJsonHandler(opts *slog.HandlerOptions) *customFormattedJsonHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	b := &bytes.Buffer{}
	return &customFormattedJsonHandler{
		buffer: b,
		handler: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		mutex: &sync.Mutex{},
	}
}

func (handler *customFormattedJsonHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.handler.Enabled(ctx, level)
}

func (handler *customFormattedJsonHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &customFormattedJsonHandler{
		handler: handler.handler.WithAttrs(attrs),
		buffer:  handler.buffer,
		mutex:   handler.mutex,
	}
}

func (handler *customFormattedJsonHandler) WithGroup(name string) slog.Handler {
	return &customFormattedJsonHandler{
		handler: handler.handler.WithGroup(name),
		buffer:  handler.buffer,
		mutex:   handler.mutex,
	}
}

func (handler *customFormattedJsonHandler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String()

	switch record.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	attrs, err := handler.computeAttrs(ctx, record)
	if err != nil {
		return err
	}

	indent, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	fmt.Printf("%s %s: %s\n%s\n",
		colorize(lightGray, record.Time.Format(timeFormat)),
		level,
		colorize(white, record.Message),
		colorize(darkGray, string(indent)),
	)

	return nil
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func (handler *customFormattedJsonHandler) computeAttrs(ctx context.Context, record slog.Record) (map[string]any, error) {
	handler.mutex.Lock()
	defer func() {
		handler.buffer.Reset()
		handler.mutex.Unlock()
	}()

	if err := handler.handler.Handle(ctx, record); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(handler.buffer.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}

func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey || a.Key == slog.LevelKey || a.Key == slog.MessageKey {
			return slog.Attr{}
		}

		if next == nil {
			return a
		}

		return next(groups, a)
	}
}

func Error(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
