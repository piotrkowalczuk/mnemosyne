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

// Storage combines API that needs to be implemented by any storage to be replaceable.
type Storage interface {
	Setup() error
	TearDown() error
	Start(string, string, map[string]string) (*mnemosynerpc.Session, error)
	Abandon(*mnemosynerpc.AccessToken) (bool, error)
	Get(*mnemosynerpc.AccessToken) (*mnemosynerpc.Session, error)
	List(int64, int64, *time.Time, *time.Time) ([]*mnemosynerpc.Session, error)
	Exists(*mnemosynerpc.AccessToken) (bool, error)
	Delete(*mnemosynerpc.AccessToken, *time.Time, *time.Time) (int64, error)
	SetValue(*mnemosynerpc.AccessToken, string, string) (map[string]string, error)
}
