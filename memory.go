package log

// NewMemoryLogger returns a logger that caches Messages, unchanged, in memory.
func NewMemoryLogger(opts *Options) *Logger {
	return NewLogger(&Memory{}, opts)
}

// Memory stores Messages in an in-memory cache. This is useful for writing
// tests or doing debugging where you need to access the individual fields
// directly without a second conversion step.
type Memory struct {
	Messages []Message
}

func (h *Memory) Handle(m Message) {
	if h.Messages == nil {
		h.Messages = []Message{m}
		return
	}
	h.Messages = append(h.Messages, m)
}
