package mnemosyned

import (
	"time"

	"context"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
)

const (
	// StorageEngineInMemory is not implemented yet.
	StorageEngineInMemory = "in_memory"
	// StorageEnginePostgres keeps session within postgres database.
	StorageEnginePostgres = "postgres"
	// StorageEngineRedis is not implemented yet.
	StorageEngineRedis = "redis"
)

// storage combines API that needs to be implemented by any storage to be replaceable.
type storage interface {
	Setup() error
	TearDown() error
	Start(context.Context, string, string, string, string, map[string]string) (*mnemosynerpc.Session, error)
	Abandon(context.Context, string) (bool, error)
	Get(context.Context, string) (*mnemosynerpc.Session, error)
	List(context.Context, int64, int64, *time.Time, *time.Time) ([]*mnemosynerpc.Session, error)
	Exists(context.Context, string) (bool, error)
	Delete(context.Context, string, string, *time.Time, *time.Time) (int64, error)
	SetValue(context.Context, string, string, string) (map[string]string, error)
}
