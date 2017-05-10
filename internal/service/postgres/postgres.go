package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"go.uber.org/zap"
)

var Timeout = errors.New("postgres connection timout")

type Opts struct {
	Logger         *zap.Logger
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

	opts.Logger.Debug("postgres connection attempt", zap.String("postgres_host", u.Host), zap.String("postgres_user", username))

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
					opts.Logger.Debug("postgres connection ping failure", zap.String("postgres_host", u.Host), zap.String("postgres_user", username))

					cancel()
					continue PingLoop
				}
				opts.Logger.Info("postgres connection has been established", zap.String("postgres_host", u.Host), zap.String("postgres_user", username))

				cancel()
				break PingLoop
			case <-cancel.C:
				return nil, Timeout
			}
		}
	}

	opts.Logger.Info("postgres connection has been established", zap.String("address", address))

	return db, nil
}
