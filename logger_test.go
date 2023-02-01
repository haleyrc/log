package log_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/haleyrc/log"
	"github.com/haleyrc/tag"
)

func TestDebugMessagesAreSuppressedWhenDebuggingIsntEnabled(t *testing.T) {
	var h log.Memory
	logger := log.NewLogger(&h, nil)

	logger.Debug(context.Background(), "this is a test", log.F{})

	if len(h.Messages) != 0 {
		t.Errorf("Expected to have 0 messages, but had %d.", len(h.Messages))
	}

	logger.EnableDebug = true
	logger.Debug(context.Background(), "this is a test", log.F{})

	if len(h.Messages) != 1 {
		t.Errorf("Expected to have 1 message, but had %d.", len(h.Messages))
	}
}

func TestDebugShouldSetTheDebugLogLevel(t *testing.T) {
	var h log.Memory
	logger := log.NewLogger(&h, &log.Options{Debug: true})

	logger.Debug(context.Background(), "this is a test", log.F{})

	got := h.Messages[0]
	if got.Level != log.LevelDebug {
		t.Errorf("Expected log level to be %s, but got %s.", log.LevelDebug, got.Level)
	}
}

func TestInfoShouldSetTheInfoLogLevel(t *testing.T) {
	var h log.Memory
	logger := log.NewLogger(&h, nil)

	logger.Info(context.Background(), "this is a test", log.F{})

	got := h.Messages[0]
	if got.Level != log.LevelInfo {
		t.Errorf("Expected log level to be %s, but got %s.", log.LevelInfo, got.Level)
	}
}

func TestErrorShouldSetTheErrorLogLevel(t *testing.T) {
	var h log.Memory
	logger := log.NewLogger(&h, nil)

	logger.Error(context.Background(), "this is a test", log.F{})

	got := h.Messages[0]
	if got.Level != log.LevelError {
		t.Errorf("Expected log level to be %s, but got %s.", log.LevelError, got.Level)
	}
}

func TestMessageShouldIncludeATimestamp(t *testing.T) {
	var h log.Memory
	logger := log.NewLogger(&h, nil)

	logger.Info(context.Background(), "this is a test", log.F{})

	got := h.Messages[0]
	if got.Timestamp.IsZero() {
		t.Errorf("Expected timestamp to be set, but it was blank.")
	}
}

func TestMessageShouldIncludeTheMessage(t *testing.T) {
	var h log.Memory
	logger := log.NewLogger(&h, nil)

	want := "this is a test"
	logger.Error(context.Background(), want, log.F{})

	got := h.Messages[0]
	if got.Message != want {
		t.Errorf("Expected message to be %q, but got %q.", want, got.Message)
	}
}

func TestMessageShouldIncludeTheFields(t *testing.T) {
	var h log.Memory
	logger := log.NewLogger(&h, nil)

	want := log.F{
		"key1": "value",
		"key2": []string{"value21", "value22"},
	}
	logger.Error(context.Background(), "this is a test", want)

	got := h.Messages[0]
	expectFieldsEqual(t, want, got.Fields)
}

func TestMessageShouldIncludeTagsFromContext(t *testing.T) {
	var h log.Memory
	logger := log.NewLogger(&h, nil)

	ctx := context.Background()
	want := tag.NewGroup(tag.Dict{
		"env":     "prod",
		"service": "database",
	})
	ctx = tag.WithGroup(ctx, want)
	logger.Error(ctx, "this is a test", log.F{})

	got := h.Messages[0]
	expectTagsEqual(t, want.Map(), got.Tags)
}

func expectFieldsEqual(t *testing.T, got, want log.F) {
	if len(got) != len(want) {
		t.Errorf(
			"Fields are of different lengths. Wanted %d, but got %d.",
			len(want), len(got),
		)
		return
	}
	for k := range got {
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Expected to find value %v for key %s, but found %v.", want[k], k, got[k])
		}
	}
}

func expectTagsEqual(t *testing.T, got, want map[string]string) {
	if len(got) != len(want) {
		t.Errorf(
			"Tags are of different lengths. Wanted %d, but got %d.",
			len(want), len(got),
		)
		return
	}
	for k := range got {
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Expected to find value %v for key %s, but found %v.", want[k], k, got[k])
		}
	}
}
