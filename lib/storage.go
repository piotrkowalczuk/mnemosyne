package lib

const (
	// SessionStorageEngineInMemory ...
	SessionStorageEngineInMemory = "in_memory"
	// SessionStorageEnginePostgres ...
	SessionStorageEnginePostgres = "postgres"
	// SessionStorageEngineRedis ...
	SessionStorageEngineRedis = "redis"
)

// SessionStorage ...
type SessionStorage interface {
	Get(SessionID) (*Session, error)
	Exists(SessionID) (bool, error)
	New(SessionData) (*Session, error)
	Abandon(SessionID) error
	SetData(SessionDataEntry) (*Session, error)
}
