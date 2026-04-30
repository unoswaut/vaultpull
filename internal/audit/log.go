package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Path      string    `json:"path,omitempty"`
	Target    string    `json:"target,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// Logger writes structured audit entries to a destination.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger writing to the given writer.
// Pass nil to use os.Stdout.
func NewLogger(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// Log writes a single audit entry as a JSON line.
func (l *Logger) Log(event, path, target string, err error) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Path:      path,
		Target:    target,
	}
	if err != nil {
		e.Error = err.Error()
	}
	b, encErr := json.Marshal(e)
	if encErr != nil {
		return fmt.Errorf("audit: marshal: %w", encErr)
	}
	_, writeErr := fmt.Fprintf(l.w, "%s\n", b)
	return writeErr
}
