package session

import "time"

type Session struct {
	ID        int64
	Token     string
	UserID    int64
	UserAgent string
	IP        string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type LoginAttempt struct {
	Attempt     int       `json:"attempt"`
	BlockTime   time.Time `json:"block_time"`
	RefreshTime time.Time `json:"refresh_time"`
}
