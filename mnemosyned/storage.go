package main

import (
	"errors"
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
)

const (
	storageEngineInMemory = "in_memory"
	storageEnginePostgres = "postgres"
	storageEngineRedis    = "redis"
)

var (
	errSessionNotFound = errors.New("mnemosyned: session not found")
)

// Storage combines API that needs to be implemented by any storage to be replaceable.
type Storage interface {
	Setup() error
	TearDown() error
	Get(*mnemosyne.Token) (*mnemosyne.Session, error)
	List(int64, int64, *time.Time, *time.Time) ([]*mnemosyne.Session, error)
	Exists(*mnemosyne.Token) (bool, error)
	Create(map[string]string) (*mnemosyne.Session, error)
	Abandon(*mnemosyne.Token) (bool, error)
	SetData(*mnemosyne.Token, string, string) (*mnemosyne.Session, error)
	Delete(*mnemosyne.Token, *time.Time, *time.Time) (int64, error)
}
