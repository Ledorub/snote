package internal

import "time"

type Note struct {
	ID                uint64
	Content           string
	CreatedAt         time.Time
	ExpiresAt         time.Time
	ExpiresAtTimeZone *time.Location
	KeyHash           []byte
}
