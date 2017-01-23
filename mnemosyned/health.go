package mnemosyned

import (
	"database/sql"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
)

type healthHandler struct {
	logger   log.Logger
	postgres *sql.DB
}

func (hh *healthHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if hh.postgres != nil {
		if err := hh.postgres.Ping(); err != nil {
			sklog.Debug(hh.logger, "health check failure due to postgres connection")
			http.Error(rw, "postgres ping failure", http.StatusServiceUnavailable)
			return
		}
	}

	sklog.Debug(hh.logger, "successful health check")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("1"))
}
