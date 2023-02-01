package log_test

import (
	"context"
	"os"

	"github.com/haleyrc/log"
	"github.com/haleyrc/tag"
)

func Example() {
	ctx := context.Background()
	ctx = tag.WithGroup(ctx, tag.NewGroup(tag.Dict{
		"env":     "prod",
		"service": "database",
	}))

	log.SetDebug(true)
	log.Debug(ctx, "this is a test", log.F{
		"user_id": 123,
		"role":    "admin",
	})
	log.Info(ctx, "this is a test", log.F{
		"user_id": 123,
		"role":    "admin",
	})
	log.Error(ctx, "this is a test", log.F{
		"user_id": 123,
		"role":    "admin",
	})
}

func ExampleLogger() {
	ctx := context.Background()
	ctx = tag.WithGroup(ctx, tag.NewGroup(tag.Dict{
		"env":     "prod",
		"service": "database",
	}))

	logger := log.NewJSONLogger(os.Stdout, &log.Options{
		Debug: true,
	})
	logger.Debug(ctx, "this is a test", log.F{
		"user_id": 123,
		"role":    "admin",
	})
	logger.Info(ctx, "this is a test", log.F{
		"user_id": 123,
		"role":    "admin",
	})
	logger.Error(ctx, "this is a test", log.F{
		"user_id": 123,
		"role":    "admin",
	})
}
