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
	// setup storage
	Setup() error
	// teardown storage
	TearDown() error
	// start user session
	Start(string, string, map[string]string) (*mnemosynerpc.Session, error)
	// abandon user session
	Abandon(string) (bool, error)
	// get user session information
	Get(string) (*mnemosynerpc.Session, error)
	// get list of sessions
	List(int64, int64, *time.Time, *time.Time) ([]*mnemosynerpc.Session, error)
	// check if user session exists
	Exists(string) (bool, error)
	// delete user session
	Delete(string, *time.Time, *time.Time) (int64, error)
	// set value in user session
	SetValue(string, string, string) (map[string]string, error)
}
