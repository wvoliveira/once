package util

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"path/filepath"
	"strings"
)

func LogHandler(logLevel slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				// The value of slog.SourceKey is a *slog.Source
				// We need to assert its type to access the File field.
				if source, ok := a.Value.Any().(*slog.Source); ok {
					source.File = lastPartPath(source.File)
					return slog.Any(a.Key, source) // Return the modified source
				}
			}
			return a // Return other attributes unchanged
		},
	}
}

func lastPartPath(path string) string {
	clean := filepath.Clean(path)
	parts := strings.Split(clean, string(filepath.Separator))

	if len(parts) >= 2 {
		return filepath.Join(parts[len(parts)-2:]...)
	}
	return clean
}

func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
