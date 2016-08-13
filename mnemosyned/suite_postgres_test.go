package mnemosyned

import (
	"database/sql"
	"testing"

	"github.com/go-kit/kit/log"
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
	ps.db, err = initPostgres(testPostgresAddress, ps.logger)
	if err != nil {
		t.Fatal(err)
	}

	if ps.store, err = initStorage(true, newPostgresStorage("session", ps.db, &monitoring{}, DefaultTTL), ps.logger); err != nil {
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
