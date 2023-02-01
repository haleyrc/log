package log_test

import (
	"fmt"
	"time"

	"github.com/haleyrc/log"
)

func ExampleMemory() {
	var mem log.Memory

	t, _ := time.ParseInLocation(time.RFC3339, "2023-01-31T20:28:16.990084-05:00", time.UTC)
	tags := map[string]string{"env": "prod", "service": "database"}

	mem.Handle(log.Message{
		Level:     log.LevelDebug,
		Timestamp: t,
		Message:   "test debug",
		Tags:      tags,
		Fields: log.F{
			"username": "joeuser",
		},
	})

	mem.Handle(log.Message{
		Level:     log.LevelInfo,
		Timestamp: t,
		Message:   "test info",
		Tags:      tags,
		Fields: log.F{
			"email": "joe.user@example.com",
			"address": log.F{
				"street": "123 Nowhere Ave.",
				"city":   "Anywhere",
				"state":  "NY",
			},
		},
	})

	mem.Handle(log.Message{
		Level:     log.LevelError,
		Timestamp: t,
		Message:   "test error",
		Tags:      tags,
		Fields: log.F{
			"code": 404,
		},
	})

	for idx, msg := range mem.Messages {
		fmt.Printf("%d: %s\n", idx, msg)
	}

	// Output: 0: <Message level:DEBUG timestamp:"2023-01-31 20:28:16.990084 -0500 -0500" msg:"test debug" tags:map[env:prod service:database] fields:map[username:joeuser]>
	// 1: <Message level:INFO timestamp:"2023-01-31 20:28:16.990084 -0500 -0500" msg:"test info" tags:map[env:prod service:database] fields:map[address:map[city:Anywhere state:NY street:123 Nowhere Ave.] email:joe.user@example.com]>
	// 2: <Message level:ERROR timestamp:"2023-01-31 20:28:16.990084 -0500 -0500" msg:"test error" tags:map[env:prod service:database] fields:map[code:404]>
}
