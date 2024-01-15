package mnemosyned

import (
	"database/sql"
	"testing"

	"go.uber.org/zap"

	"github.com/piotrkowalczuk/mnemosyne/internal/service/postgres"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	storagepq "github.com/piotrkowalczuk/mnemosyne/internal/storage/postgres"
)

type postgresSuite struct {
	db     *sql.DB
	logger *zap.Logger
	store  storage.Storage
}

func (ps *postgresSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres suite ignored in short mode")
	}

	var err error

	ps.logger = zap.L()
	// ps.logger = sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter)
	ps.db, err = postgres.Init(testPostgresAddress, postgres.Opts{
		Logger: ps.logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	if ps.store, err = storage.Init(storagepq.NewStorage(storagepq.StorageOpts{
		Table: "session", Schema: "mnemosyne", Conn: ps.db, TTL: storage.DefaultTTL,
	}), true); err != nil {
		t.Fatal(err)
	}
}

func (ps *postgresSuite) teardown(t *testing.T) {
	var err error

	if err = ps.store.TearDown(); err != nil {
		t.Fatal(err)
	}
	if err = ps.db.Close(); err != nil {
		t.Fatal(err)
	}
}
