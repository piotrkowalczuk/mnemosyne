package lib

import (
	"errors"
	"time"
)

var (
	// ErrSessionDoesNotExists ...
	ErrSessionDoesNotExists = errors.New("session: session does not exists")
)

// Session ...
type Session struct {
	ID       SessionID
	Data     SessionData
	ExpireAt *time.Time
}

// NewSession ...
func NewSession(id SessionID, data SessionData, expireAt *time.Time) *Session {
	return &Session{
		ID:       id,
		Data:     data,
		ExpireAt: expireAt,
	}
}
