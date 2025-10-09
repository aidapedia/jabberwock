package session

import "time"

type Session struct {
	ID        string
	UserID    int64
	TokenID   string
	UserAgent string
	IP        string
}

type LoginAttempt struct {
	Attempt     int       `json:"attempt"`
	BlockTime   time.Time `json:"block_time"`
	RefreshTime time.Time `json:"refresh_time"`
}
