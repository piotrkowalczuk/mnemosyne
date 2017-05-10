package mnemosyned

import (
	"database/sql"
	"testing"

	"go.uber.org/zap"

	"github.com/piotrkowalczuk/mnemosyne/internal/service/postgres"
)

type postgresSuite struct {
	db     *sql.DB
	logger *zap.Logger
	store  storage
}

func (ps *postgresSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres suite ignored in short mode")
	}

	var err error

	ps.logger = zap.L()
	//ps.logger = sklog.NewHumaneLogger(os.Stdout, sklog.DefaultHTTPFormatter)
	ps.db, err = postgres.Init(testPostgresAddress, postgres.Opts{
		Logger: ps.logger,
	})
	if err != nil {
		t.Fatal(err)
	}

	if ps.store, err = initStorage(true, newPostgresStorage("session", "mnemosyne", ps.db, &monitoring{}, DefaultTTL)); err != nil {
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
