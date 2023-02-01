package log

import (
	"context"
	"time"

	"github.com/haleyrc/tag"
)

// Options can be passed to the constructors provided in this package to
// customize the behavior of the resulting Logger.
type Options struct {
	// Debug enables debug level messages.
	Debug bool
}

// NewLogger returns a Logger that sends messages to the provided MessageHandler
// subject to the provided Options.
func NewLogger(h MessageHandler, opts *Options) *Logger {
	l := Logger{Handler: h}

	if opts != nil {
		l.EnableDebug = opts.Debug
	}

	return &l
}

// Logger provides methods for creating Messages and passing them to an
// implementation of the MessageHandler interface for formatting.
type Logger struct {
	// If true, debug messages will be passed through to the Handler, o.w. they
	// are dropped.
	EnableDebug bool

	// Handler is responsible for formatting/handling Messages passed to it by
	// the logging methods (Debug, Info, Error).
	Handler MessageHandler
}

// Debug logs a message at the DEBUG log level. If debug output is not disabled,
// this is a no-op.
func (l *Logger) Debug(ctx context.Context, msg string, fields F) {
	if !l.EnableDebug {
		return
	}
	l.log(ctx, LevelDebug, msg, fields)
}

// Info logs a message at the INFO log level.
func (l *Logger) Info(ctx context.Context, msg string, fields F) {
	l.log(ctx, LevelInfo, msg, fields)
}

// Error logs a message at the ERROR log level.
func (l *Logger) Error(ctx context.Context, msg string, fields F) {
	l.log(ctx, LevelError, msg, fields)
}

func (l *Logger) log(ctx context.Context, level Level, msg string, fields F) {
	tags := tag.FromContext(ctx)
	l.Handler.Handle(Message{
		Level:     level,
		Timestamp: time.Now(),
		Message:   msg,
		Tags:      tags.Map(),
		Fields:    fields,
	})
}
