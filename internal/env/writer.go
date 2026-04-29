package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	filePath string
	backupOnRotate bool
}

// NewWriter creates a new Writer for the given file path.
func NewWriter(filePath string, backupOnRotate bool) *Writer {
	return &Writer{
		filePath:       filePath,
		backupOnRotate: backupOnRotate,
	}
}

// Write serializes the provided secrets map into the .env file.
// If backupOnRotate is enabled and the file already exists, a backup is created first.
func (w *Writer) Write(secrets map[string]string) error {
	if w.backupOnRotate {
		if err := w.backup(); err != nil {
			return fmt.Errorf("env writer: backup failed: %w", err)
		}
	}

	f, err := os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("env writer: open file: %w", err)
	}
	defer f.Close()

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, escape(secrets[k])))
	}

	if _, err := f.WriteString(sb.String()); err != nil {
		return fmt.Errorf("env writer: write: %w", err)
	}
	return nil
}

// backup copies the existing .env file to <filePath>.bak if it exists.
func (w *Writer) backup() error {
	data, err := os.ReadFile(w.filePath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.WriteFile(w.filePath+".bak", data, 0600)
}

// escape wraps values containing spaces or special characters in double quotes.
func escape(value string) string {
	if strings.ContainsAny(value, " \t\n\r#") {
		escaped := strings.ReplaceAll(value, `"`, `\"`)
		return `"` + escaped + `"`
	}
	return value
}
