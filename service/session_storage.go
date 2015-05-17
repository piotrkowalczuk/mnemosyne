package service

import (
	"errors"

	"github.com/go-soa/mnemosyne/lib"
)

var (
	// ErrSessionStorageEngineNotSupported ...
	ErrSessionStorageEngineNotSupported = errors.New("service: session storage engine not supported")
	// SessionStorage ...
	SessionStorage lib.SessionStorage
)

// SessionStorageConfig  ...
type SessionStorageConfig struct {
	Engine   string         `xml:"engine"`
	Postgres PostgresConfig `xml:"postgres"`
}

// InitSessionStorage ...
func InitSessionStorage(config SessionStorageConfig) {
	switch config.Engine {
	case lib.SessionStorageEngineInMemory:
		//		SessionStorage = lib.NewInMemorySessionStorage()
	case lib.SessionStorageEnginePostgres:
		SessionStorage = lib.NewPostgresStorage(Postgres, config.Postgres.TableName)
	case lib.SessionStorageEngineRedis:
	default:
		Logger.Fatal(ErrSessionStorageEngineNotSupported)
	}
}
