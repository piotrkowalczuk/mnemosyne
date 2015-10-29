package main

import (
	"errors"
	"time"

	"github.com/piotrkowalczuk/mnemosyne/shared"
)

const (
	storageEngineInMemory = "in_memory"
	storageEnginePostgres = "postgres"
	storageEngineRedis    = "redis"
)

var (
	errSessionNotFound = errors.New("mnemosyne: session not found")
)

// Storage combines API that needs to be implemented by any storage to be replaceable.
type Storage interface {
	Setup() error
	TearDown() error
	Get(*shared.ID) (*shared.Session, error)
	List(int64, int64, *time.Time, *time.Time) (*shared.Session, error)
	Exists(*shared.ID) (bool, error)
	Create(map[string]string) (*shared.Session, error)
	Abandon(*shared.ID) (bool, error)
	SetData(*shared.ID, string, string) (*shared.Session, error)
	Delete(*shared.ID, *time.Time, *time.Time) (int64, error)
}
