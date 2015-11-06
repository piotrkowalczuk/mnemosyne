package main

import (
	"errors"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
)

const (
	storageEngineInMemory = "in_memory"
	storageEnginePostgres = "postgres"
	storageEngineRedis    = "redis"
)

var (
	storage            Storage
	errSessionNotFound = errors.New("mnemosyne: session not found")
)

// Storage combines API that needs to be implemented by any storage to be replaceable.
type Storage interface {
	Setup() error
	TearDown() error
	Get(*mnemosyne.ID) (*mnemosyne.Session, error)
	List(int64, int64, *time.Time, *time.Time) ([]*mnemosyne.Session, error)
	Exists(*mnemosyne.ID) (bool, error)
	Create(map[string]string) (*mnemosyne.Session, error)
	Abandon(*mnemosyne.ID) (bool, error)
	SetData(*mnemosyne.ID, string, string) (*mnemosyne.Session, error)
	Delete(*mnemosyne.ID, *time.Time, *time.Time) (int64, error)
}

func initStorage(fn func() (Storage, error), logger log.Logger) {
	s, err := fn()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	err = s.Setup()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	storage = s
}
