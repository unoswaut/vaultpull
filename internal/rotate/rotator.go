package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Rotator handles backup rotation of .env files before overwriting.
type Rotator struct {
	MaxBackups int
}

// New returns a Rotator with the given max backup count.
// If maxBackups is zero or negative, it defaults to 3.
func New(maxBackups int) *Rotator {
	if maxBackups <= 0 {
		maxBackups = 3
	}
	return &Rotator{MaxBackups: maxBackups}
}

// Rotate creates a timestamped backup of the file at path if it exists,
// then prunes old backups so that at most MaxBackups copies are kept.
func (r *Rotator) Rotate(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	dir := filepath.Dir(path)
	base := filepath.Base(path)
	timestamp := time.Now().UTC().Format("20060102T150405Z")
	backupName := fmt.Sprintf("%s.%s.bak", base, timestamp)
	backupPath := filepath.Join(dir, backupName)

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("rotate: read %s: %w", path, err)
	}

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("rotate: write backup %s: %w", backupPath, err)
	}

	return r.pruneBackups(dir, base)
}

// pruneBackups removes the oldest backups when the count exceeds MaxBackups.
func (r *Rotator) pruneBackups(dir, base string) error {
	pattern := filepath.Join(dir, base+".*.bak")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("rotate: glob %s: %w", pattern, err)
	}

	// filepath.Glob returns sorted results; oldest timestamps sort first.
	for len(matches) > r.MaxBackups {
		if err := os.Remove(matches[0]); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("rotate: remove %s: %w", matches[0], err)
		}
		matches = matches[1:]
	}
	return nil
}
