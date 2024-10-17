package log

type config struct {
	debug      bool
	freezeTime bool
}

// An Option is a function that modifies the default handler configuration.
type Option func(*config)

// Debug forces a handler to output debug lines.
func Debug() Option {
	return func(cfg *config) {
		cfg.debug = true
	}
}

// FreezeTime forces a handler to use a constant timestamp when emitting log
// lines. This is only useful for testing where deterministic output is
// required.
func FreezeTime() Option {
	return func(cfg *config) {
		cfg.freezeTime = true
	}
}
