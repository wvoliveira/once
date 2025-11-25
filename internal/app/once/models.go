package once

import "time"

type Content struct {
	ID        string
	Text      string
	FileID    string
	FileName  string
	MimeType  string
	ExpiresAt time.Time
	CreatedAt time.Time
}
