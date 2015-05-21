package lib

import (
	"errors"
	"time"
)

const (
	// SessionStorageEngineInMemory ...
	SessionStorageEngineInMemory = "in_memory"
	// SessionStorageEnginePostgres ...
	SessionStorageEnginePostgres = "postgres"
	// SessionStorageEngineRedis ...
	SessionStorageEngineRedis = "redis"
)

var (
	// ErrSessionNotFound ...
	ErrSessionNotFound = errors.New("mnemosyne: session not found")
)

// SessionStorage ...
type SessionStorage interface {
	Get(SessionID) (*Session, error)
	Exists(SessionID) (bool, error)
	New(SessionData) (*Session, error)
	Abandon(SessionID) error
	SetData(SessionDataEntry) (*Session, error)
	Cleanup(*time.Time) (int64, error)
}
