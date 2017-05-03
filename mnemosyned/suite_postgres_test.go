package mnemosyned

import (
	"database/sql"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne/internal/service/postgres"
	"github.com/piotrkowalczuk/sklog"
)

type postgresSuite struct {
	db     *sql.DB
	logger log.Logger
	store  storage
}

func (ps *postgresSuite) setup(t *testing.T) {
	if testing.Short() {
		t.Skip("postgres suite ignored in short mode")
	}

	var err error

	ps.logger = sklog.NewTestLogger(t)
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
