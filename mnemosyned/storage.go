package mnemosyned

import (
	"errors"
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
)

const (
	StorageEngineInMemory = "in_memory"
	StorageEnginePostgres = "postgres"
	StorageEngineRedis    = "redis"
)

var (
	SessionNotFound = errors.New("session not found")
)

// Storage combines API that needs to be implemented by any storage to be replaceable.
type Storage interface {
	Setup() error
	TearDown() error

	Start(string, map[string]string) (*mnemosyne.Session, error)
	Abandon(*mnemosyne.AccessToken) (bool, error)
	Get(*mnemosyne.AccessToken) (*mnemosyne.Session, error)
	List(int64, int64, *time.Time, *time.Time) ([]*mnemosyne.Session, error)
	Exists(*mnemosyne.AccessToken) (bool, error)
	Delete(*mnemosyne.AccessToken, *time.Time, *time.Time) (int64, error)

	SetValue(*mnemosyne.AccessToken, string, string) (map[string]string, error)
	//	DeleteValue(*mnemosyne.AccessToken, string) (*mnemosyne.Session, error)
	//	Clear(*mnemosyne.AccessToken) (*mnemosyne.Session, error)
}
