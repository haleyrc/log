package log

import (
	"context"
	"fmt"
	"os"
	"time"
)

var dl = NewJSONLogger(os.Stdout, &Options{
	Debug: true,
})

// Debug logs a message at the DEBUG log level using the default logger. If
// debug output is not disabled, this is a no-op.
func Debug(ctx context.Context, msg string, fields F) { dl.Debug(ctx, msg, fields) }

// Info logs a message at the INFO log level using the default logger.
func Info(ctx context.Context, msg string, fields F) { dl.Info(ctx, msg, fields) }

// Error logs a message at the ERROR log level using the default logger.
func Error(ctx context.Context, msg string, fields F) { dl.Error(ctx, msg, fields) }

// SetDebug enables debug messages from the default logger.
func SetDebug(value bool) { dl.EnableDebug = value }

// SetHandler sets the message handler for the default logger.
func SetHandler(h MessageHandler) { dl.Handler = h }

// F is a shortcut for passing fields to the logging methods.
type F map[string]interface{}

// MessageHandler is used by a Logger to format incoming messages.
type MessageHandler interface {
	Handle(Message)
}

// Level loosely corresponds to a standard syslog level.
type Level string

// The following Levels are supported. Note that not all of the standard syslog
// levels are available. The README goes into greater detail on the motivation
// for this restriction.
const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelError Level = "ERROR"
)

// Message is the basic unit that is passed to all implementations of
// MessageHandler. It represents a single "log line" resulting from a call to
// one of the logging methods (Debug, Info, Error).
type Message struct {
	// The log level for the message.
	Level Level `json:"level"`

	// The server time that the message was created.
	Timestamp time.Time `json:"timestamp"`

	// A human-readable description of the logged event. This should usually be
	// a basic string (read: no interpolation) in order to ensure that logs are
	// maximally greppable. If you have data points to include, use the Tags or
	// Fields fields.
	Message string `json:"message"`

	// Tags are key/value pairs that are included with the message. These are
	// usually set per-request e.g. by middleware.
	Tags map[string]string `json:"tags"`

	// Fields represents a generic blob of data. This is the primary mechanism
	// for including specific data points with log messages.
	Fields F `json:"fields"`
}

func (msg Message) String() string {
	return fmt.Sprintf(
		"<Message level:%s timestamp:%q msg:%q tags:%v fields:%v>",
		msg.Level,
		msg.Timestamp,
		msg.Message,
		msg.Tags,
		msg.Fields,
	)
}
