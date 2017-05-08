package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
)

var Timeout = errors.New("postgres connection timout")

type Opts struct {
	Logger         log.Logger
	Retry, Timeout time.Duration
}

// Init ...
func Init(address string, opts Opts) (*sql.DB, error) {
	timeout, retry := opts.Timeout, opts.Retry
	if timeout == time.Duration(0) {
		timeout = 10 * time.Second
	}
	if retry == time.Duration(0) {
		retry = 1 * time.Second
	}

	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	username := ""
	if u.User != nil {
		username = u.User.Username()
	}

	sklog.Debug(opts.Logger, "postgres connection attempt", "postgres_host", u.Host, "postgres_user", username)

	db, err := sql.Open("postgres", address)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failure: %s", err.Error())
	}

	// Otherwise 1 second cooldown is going to be multiplied by number of tests.
	ctx, cancel := context.WithTimeout(context.Background(), retry)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		cancel := time.NewTimer(timeout)

	PingLoop:
		for {
			select {
			case <-time.After(retry):
				ctx, cancel := context.WithTimeout(context.Background(), retry)
				if err := db.PingContext(ctx); err != nil {
					sklog.Debug(opts.Logger, "postgres connection ping failure", "postgres_host", u.Host, "postgres_user", username)

					cancel()
					continue PingLoop
				}
				sklog.Info(opts.Logger, "postgres connection has been established", "postgres_host", u.Host, "postgres_user", username)

				cancel()
				break PingLoop
			case <-cancel.C:
				return nil, Timeout
			}
		}
	}

	sklog.Info(opts.Logger, "postgres connection has been established", "address", address)

	return db, nil
}
