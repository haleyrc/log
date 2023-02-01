package log

import (
	"encoding/json"
	"fmt"
	"io"
)

// NewJSONLogger returns a Logger that writes JSON-formatted messages to the
// provided io.Writer, one message per line.
func NewJSONLogger(w io.Writer, opts *Options) *Logger {
	h := &JSON{Writer: w}
	return NewLogger(h, opts)
}

// JSON writes Messages to its Writer in JSON format, one message per line.
type JSON struct {
	Writer io.Writer
}

func (h *JSON) Handle(m Message) {
	b, _ := json.Marshal(m)
	fmt.Fprintln(h.Writer, string(b))
}
