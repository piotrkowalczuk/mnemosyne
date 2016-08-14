package mnemosyned

import (
	"database/sql"
	"net/http"
)

type healthHandler struct {
	postgres *sql.DB
}

func (hh *healthHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if hh.postgres != nil {
		if err := hh.postgres.Ping(); err != nil {
			http.Error(rw, "postgres ping failure", http.StatusServiceUnavailable)
			return
		}
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("1"))
}
