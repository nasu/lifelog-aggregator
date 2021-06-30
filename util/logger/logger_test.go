package logger

import (
	"bytes"
	"log"
	"testing"
)

type testObject struct {
	title string
	level Level
	want  string
}

func TestTrace(t *testing.T) {
	tests := []testObject{
		{"trace", TRACE, "test 1 foo\n"},
		{"debug", DEBUG, ""},
		{"info", INFO, ""},
		{"warn", WARN, ""},
		{"error", ERROR, ""},
	}

	for _, tt := range tests {
		l := NewLogger(tt.level)
		got := spyOutput(l, tt, func() { l.Trace("test %d %s", 1, "foo") })
		if got != tt.want {
			t.Errorf("%s got=%s wnat=%s", tt.title, got, tt.want)
		}

	}
}

func TestDebug(t *testing.T) {
	tests := []testObject{
		{"trace", TRACE, "test 1 foo\n"},
		{"debug", DEBUG, "test 1 foo\n"},
		{"info", INFO, ""},
		{"warn", WARN, ""},
		{"error", ERROR, ""},
	}

	for _, tt := range tests {
		l := NewLogger(tt.level)
		got := spyOutput(l, tt, func() { l.Debug("test %d %s", 1, "foo") })
		if got != tt.want {
			t.Errorf("%s got=%s wnat=%s", tt.title, got, tt.want)
		}
	}
}

func TestInfo(t *testing.T) {
	tests := []testObject{
		{"trace", TRACE, "test 1 foo\n"},
		{"debug", DEBUG, "test 1 foo\n"},
		{"info", INFO, "test 1 foo\n"},
		{"warn", WARN, ""},
		{"error", ERROR, ""},
	}

	for _, tt := range tests {
		l := NewLogger(tt.level)
		got := spyOutput(l, tt, func() { l.Info("test %d %s", 1, "foo") })
		if got != tt.want {
			t.Errorf("%s got=%s wnat=%s", tt.title, got, tt.want)
		}

	}
}
func TestWarn(t *testing.T) {
	tests := []testObject{
		{"trace", TRACE, "test 1 foo\n"},
		{"debug", DEBUG, "test 1 foo\n"},
		{"info", INFO, "test 1 foo\n"},
		{"warn", WARN, "test 1 foo\n"},
		{"error", ERROR, ""},
	}

	for _, tt := range tests {
		l := NewLogger(tt.level)
		got := spyOutput(l, tt, func() { l.Warn("test %d %s", 1, "foo") })
		if got != tt.want {
			t.Errorf("%s got=%s wnat=%s", tt.title, got, tt.want)
		}

	}
}
func TestError(t *testing.T) {
	tests := []testObject{
		{"trace", TRACE, "test 1 foo\n"},
		{"debug", DEBUG, "test 1 foo\n"},
		{"info", INFO, "test 1 foo\n"},
		{"warn", WARN, "test 1 foo\n"},
		{"error", ERROR, "test 1 foo\n"},
	}

	for _, tt := range tests {
		l := NewLogger(tt.level)
		got := spyOutput(l, tt, func() { l.Error("test %d %s", 1, "foo") })
		if got != tt.want {
			t.Errorf("%s got=%s wnat=%s", tt.title, got, tt.want)
		}

	}
}
func spyOutput(l *Logger, tt testObject, f func()) string {
	buf := &bytes.Buffer{}
	buf.Reset()
	l.logger = log.New(buf, "", 0)

	f()
	return buf.String()
}
