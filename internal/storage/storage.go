package storage

import (
	"errors"
	"time"

	"context"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// DefaultTTL is session time to live default value.
	DefaultTTL = 24 * time.Minute
	// DefaultTTC is time to clear default value.
	DefaultTTC = 1 * time.Minute
)

var (
	ErrSessionNotFound    = errors.New("storage: session not found")
	ErrMissingAccessToken = errors.New("storage: missing access token")
	ErrMissingSubjectID   = errors.New("storage: missing subject accessToken")
	ErrMissingSession     = errors.New("storage: missing session")
)

const (
	// EngineInMemory is not implemented yet.
	EngineInMemory = "in_memory"
	// EnginePostgres keeps session within postgres database.
	EnginePostgres = "postgres"
	// EngineRedis is not implemented yet.
	EngineRedis = "redis"
)

// Storage combines API that needs to be implemented by any storage to be replaceable.
type Storage interface {
	Setup() error
	TearDown() error
	Start(context.Context, string, string, string, string, map[string]string) (*mnemosynerpc.Session, error)
	Abandon(context.Context, string) (bool, error)
	Get(context.Context, string) (*mnemosynerpc.Session, error)
	List(context.Context, int64, int64, *time.Time, *time.Time) ([]*mnemosynerpc.Session, error)
	Exists(context.Context, string) (bool, error)
	Delete(context.Context, string, string, string, *time.Time, *time.Time) (int64, error)
	SetValue(context.Context, string, string, string) (map[string]string, error)
}

// InstrumentedStorage combines Storage and prometheus Collector interface.
type InstrumentedStorage interface {
	Storage
	prometheus.Collector
}

func Init(s Storage, isTest bool) (Storage, error) {
	if isTest {
		if err := s.TearDown(); err != nil {
			return nil, err
		}
	}
	if err := s.Setup(); err != nil {
		return nil, err
	}

	return s, nil
}
