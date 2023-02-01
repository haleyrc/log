package log_test

import (
	"os"
	"time"

	"github.com/haleyrc/log"
)

func ExampleJSON() {
	json := log.JSON{Writer: os.Stdout}

	t, _ := time.Parse(time.RFC3339, "2023-01-31T20:28:16.990084-05:00")
	json.Handle(log.Message{
		Level:     log.LevelInfo,
		Timestamp: t,
		Message:   "test",
		Tags:      map[string]string{"env": "prod", "service": "database"},
		Fields: log.F{
			"user_id": 123,
			"role":    "admin",
		},
	})

	// Output: {"level":"INFO","timestamp":"2023-01-31T20:28:16.990084-05:00","message":"test","tags":{"env":"prod","service":"database"},"fields":{"role":"admin","user_id":123}}
}
