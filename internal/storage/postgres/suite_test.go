package postgres_test

import (
	"database/sql"
	"flag"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/mnemosyne/internal/service/postgres"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	storagepq "github.com/piotrkowalczuk/mnemosyne/internal/storage/postgres"
	"go.uber.org/zap"
)

var testPostgresAddress string

func TestMain(m *testing.M) {
	flag.StringVar(&testPostgresAddress, "postgres.address", getStringEnvOr("MNEMOSYNED_POSTGRES_ADDRESS", "postgres://localhost/test?sslmode=disable"), "")
	flag.Parse()

	os.Exit(m.Run())
}

func getStringEnvOr(env, or string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return or
}

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
	ps.db, err = postgres.Init(testPostgresAddress, postgres.Opts{
		Logger: ps.logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	if ps.store, err = storage.Init(storagepq.NewStorage(storagepq.StorageOpts{
		Table:  "session",
		Schema: "mnemosyne",
		Conn:   ps.db,
		TTL:    storage.DefaultTTL,
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
