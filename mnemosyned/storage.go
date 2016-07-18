package mnemosyned

import (
	"time"

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
	Start(string, string, map[string]string) (*mnemosynerpc.Session, error)
	Abandon(string) (bool, error)
	Get(string) (*mnemosynerpc.Session, error)
	List(int64, int64, *time.Time, *time.Time) ([]*mnemosynerpc.Session, error)
	Exists(string) (bool, error)
	Delete(string, *time.Time, *time.Time) (int64, error)
	SetValue(string, string, string) (map[string]string, error)
}
