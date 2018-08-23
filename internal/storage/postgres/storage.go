package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne/internal/model"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/prometheus/client_golang/prometheus"
)

var monitoringPostgresLabels = []string{
	"query",
}

type Storage struct {
	db                                             *sql.DB
	schema                                         string
	table                                          string
	ttl                                            time.Duration
	querySave, queryGet, queryExists, queryAbandon string
	// monitoring
	queriesTotal    *prometheus.CounterVec
	queriesDuration *prometheus.HistogramVec
	errors          *prometheus.CounterVec
}

type StorageOpts struct {
	Conn          *sql.DB
	Schema, Table string
	Namespace     string
	TTL           time.Duration
}

func NewStorage(opts StorageOpts) storage.Storage {
	return &Storage{
		db:     opts.Conn,
		table:  opts.Table,
		schema: opts.Schema,
		ttl:    opts.TTL,
		querySave: `INSERT INTO ` + opts.Schema + ` .` + opts.Table + ` (access_token, refresh_token, subject_id, subject_client, bag)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING expire_at`,
		queryGet: fmt.Sprintf(`UPDATE `+opts.Schema+` .`+opts.Table+`
			SET expire_at = (NOW() + '%d seconds')
			WHERE access_token = $1
			RETURNING refresh_token, subject_id, subject_client, bag, expire_at`, int64(opts.TTL.Seconds())),
		queryExists:  `SELECT EXISTS(SELECT 1 FROM ` + opts.Schema + ` .` + opts.Table + ` WHERE access_token = $1)`,
		queryAbandon: `DELETE FROM ` + opts.Schema + ` .` + opts.Table + ` WHERE access_token = $1`,
		queriesTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: opts.Namespace,
				Subsystem: "storage",
				Name:      "postgres_queries_total",
				Help:      "Total number of SQL queries made.",
			},
			monitoringPostgresLabels,
		),
		queriesDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: opts.Namespace,
				Subsystem: "storage",
				Name:      "postgres_query_duration_seconds",
				Help:      "The SQL query latencies in seconds on the client side.",
			},
			monitoringPostgresLabels,
		),
		errors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: opts.Namespace,
				Subsystem: "storage",
				Name:      "postgres_errors_total",
				Help:      "Total number of errors that happen during SQL queries.",
			},
			monitoringPostgresLabels,
		),
	}
}

// Start implements storage interface.
func (ps *Storage) Start(ctx context.Context, accessToken, refreshToken, sid, sc string, b map[string]string) (*mnemosynerpc.Session, error) {
	ent := &sessionEntity{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		SubjectID:     sid,
		SubjectClient: sc,
		Bag:           model.Bag(b),
	}

	if err := ps.save(ctx, ent); err != nil {
		return nil, err
	}

	return ent.session()
}

func (ps *Storage) save(ctx context.Context, ent *sessionEntity) (err error) {
	start := time.Now()
	labels := prometheus.Labels{"query": "save"}
	err = ps.db.QueryRowContext(
		ctx,
		ps.querySave,
		ent.AccessToken,
		ent.RefreshToken,
		ent.SubjectID,
		ent.SubjectClient,
		ent.Bag,
	).Scan(
		&ent.ExpireAt,
	)
	ps.incQueries(labels, start)
	if err != nil {
		ps.incError(labels)
	}
	return
}

// Get implements storage interface.
func (ps *Storage) Get(ctx context.Context, accessToken string) (*mnemosynerpc.Session, error) {
	var entity sessionEntity
	start := time.Now()
	labels := prometheus.Labels{"query": "get"}

	err := ps.db.QueryRowContext(ctx, ps.queryGet, accessToken).Scan(
		&entity.RefreshToken,
		&entity.SubjectID,
		&entity.SubjectClient,
		&entity.Bag,
		&entity.ExpireAt,
	)
	ps.incQueries(labels, start)
	if err != nil {
		ps.incError(labels)
		if err == sql.ErrNoRows {
			return nil, storage.ErrSessionNotFound
		}
		return nil, err
	}

	expireAt, err := ptypes.TimestampProto(entity.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.Session{
		AccessToken:   accessToken,
		RefreshToken:  entity.RefreshToken,
		SubjectId:     entity.SubjectID,
		SubjectClient: entity.SubjectClient,
		Bag:           entity.Bag,
		ExpireAt:      expireAt,
	}, nil
}

// sessionManagerList implements storage interface.
func (ps *Storage) List(ctx context.Context, offset, limit int64, expiredAtFrom, expiredAtTo *time.Time) ([]*mnemosynerpc.Session, error) {
	if limit == 0 {
		return nil, errors.New("cannot retrieve list of sessions, limit needs to be higher than 0")
	}

	args := []interface{}{offset, limit}
	query := "SELECT access_token, refresh_token, subject_id, subject_client, bag, expire_at FROM " + ps.schema + "." + ps.table + " "
	if expiredAtFrom != nil || expiredAtTo != nil {
		query += " WHERE "
	}
	switch {
	case expiredAtFrom != nil && expiredAtTo == nil:
		query += "expire_at > $3"
		args = append(args, expiredAtFrom)
	case expiredAtFrom == nil && expiredAtTo != nil:
		query += "expire_at < $3"
		args = append(args, expiredAtTo)
	case expiredAtFrom != nil && expiredAtTo != nil:
		query += "expire_at > $3 AND expire_at < $4"
		args = append(args, expiredAtFrom, expiredAtTo)
	}

	query += " OFFSET $1 LIMIT $2"
	labels := prometheus.Labels{"query": "list"}

	start := time.Now()
	rows, err := ps.db.QueryContext(ctx, query, args...)
	ps.incQueries(labels, start)
	if err != nil {
		ps.incError(labels)
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*mnemosynerpc.Session, 0, limit)
	for rows.Next() {
		var ent sessionEntity

		err = rows.Scan(
			&ent.AccessToken,
			&ent.RefreshToken,
			&ent.SubjectID,
			&ent.SubjectClient,
			&ent.Bag,
			&ent.ExpireAt,
		)
		if err != nil {
			ps.incError(labels)
			return nil, err
		}

		expireAt, err := ptypes.TimestampProto(ent.ExpireAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &mnemosynerpc.Session{
			AccessToken:   ent.AccessToken,
			RefreshToken:  ent.RefreshToken,
			SubjectId:     ent.SubjectID,
			SubjectClient: ent.SubjectClient,
			Bag:           ent.Bag,
			ExpireAt:      expireAt,
		})
	}
	if rows.Err() != nil {
		ps.incError(labels)
		return nil, rows.Err()
	}

	return sessions, nil
}

// Exists implements storage interface.
func (ps *Storage) Exists(ctx context.Context, accessToken string) (exists bool, err error) {
	start := time.Now()
	labels := prometheus.Labels{"query": "exists"}

	err = ps.db.QueryRowContext(ctx, ps.queryExists, accessToken).Scan(
		&exists,
	)
	ps.incQueries(labels, start)
	if err != nil {
		ps.incError(labels)
	}

	return
}

// Abandon implements storage interface.
func (ps *Storage) Abandon(ctx context.Context, accessToken string) (bool, error) {
	start := time.Now()
	labels := prometheus.Labels{"query": "abandon"}

	result, err := ps.db.ExecContext(ctx, ps.queryAbandon, accessToken)
	ps.incQueries(labels, start)
	if err != nil {
		ps.incError(labels)
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, storage.ErrSessionNotFound
	}

	return true, nil
}

// SetValue implements storage interface.
func (ps *Storage) SetValue(ctx context.Context, accessToken string, key, value string) (map[string]string, error) {
	var err error
	if accessToken == "" {
		return nil, storage.ErrMissingAccessToken
	}

	entity := &sessionEntity{
		AccessToken: accessToken,
	}
	selectQuery := `
		SELECT bag
		FROM ` + ps.schema + `.` + ps.table + `
		WHERE access_token = $1
		FOR UPDATE
	`
	updateQuery := `
		UPDATE ` + ps.schema + `.` + ps.table + `
		SET
			bag = $2
		WHERE access_token = $1
	`

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	startSelect := time.Now()
	err = tx.QueryRowContext(ctx, selectQuery, accessToken).Scan(
		&entity.Bag,
	)
	ps.incQueries(prometheus.Labels{"query": "set_value_select"}, startSelect)
	if err != nil {
		ps.incError(prometheus.Labels{"query": "set_value_select"})
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, storage.ErrSessionNotFound
		}
		return nil, err
	}

	entity.Bag.Set(key, value)

	startUpdate := time.Now()
	_, err = tx.ExecContext(ctx, updateQuery, accessToken, entity.Bag)
	ps.incQueries(prometheus.Labels{"query": "set_value_update"}, startUpdate)
	if err != nil {
		ps.incError(prometheus.Labels{"query": "set_value_update"})
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return entity.Bag, nil
}

// Delete implements storage interface.
func (ps *Storage) Delete(ctx context.Context, subjectID, accessToken, refreshToken string, expiredAtFrom, expiredAtTo *time.Time) (int64, error) {
	where, args := ps.where(subjectID, accessToken, refreshToken, expiredAtFrom, expiredAtTo)
	if where.Len() == 0 {
		return 0, fmt.Errorf("session cannot be deleted, no where parameter provided: %s", where.String())
	}
	query := "DELETE FROM " + ps.schema + "." + ps.table + " WHERE " + where.String()
	labels := prometheus.Labels{"query": "delete"}
	start := time.Now()

	result, err := ps.db.Exec(query, args...)
	ps.incQueries(labels, start)
	if err != nil {
		ps.incError(labels)
		return 0, err
	}

	return result.RowsAffected()
}

// Setup implements storage interface.
func (ps *Storage) Setup() error {
	query := fmt.Sprintf(`
		CREATE SCHEMA IF NOT EXISTS %s;
		CREATE TABLE IF NOT EXISTS %s.%s (
			access_token BYTEA PRIMARY KEY,
			refresh_token BYTEA,
			subject_id TEXT NOT NULL,
			subject_client TEXT,
			bag bytea NOT NULL,
			expire_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + '%d seconds')

		);
		CREATE INDEX ON %s.%s (refresh_token);
		CREATE INDEX ON %s.%s (subject_id);
		CREATE INDEX ON %s.%s (expire_at DESC);
	`, ps.schema, ps.schema, ps.table, int64(ps.ttl.Seconds()),
		ps.schema, ps.table,
		ps.schema, ps.table,
		ps.schema, ps.table,
	)
	_, err := ps.db.Exec(query)

	return err
}

// TearDown implements storage interface.
func (ps *Storage) TearDown() error {
	_, err := ps.db.Exec(`DROP SCHEMA IF EXISTS ` + ps.schema + ` CASCADE`)

	return err
}

func (ps *Storage) incQueries(field prometheus.Labels, start time.Time) {
	ps.queriesTotal.With(field).Inc()
	ps.queriesDuration.With(field).Observe(time.Since(start).Seconds())
}

func (ps *Storage) incError(field prometheus.Labels) {
	ps.errors.With(field).Inc()
}

func (ps *Storage) where(subjectID, accessToken, refreshToken string, expiredAtFrom, expiredAtTo *time.Time) (*bytes.Buffer, []interface{}) {
	var count int
	buf := bytes.NewBuffer(nil)
	args := make([]interface{}, 0, 4)

	switch {
	case subjectID != "":
		count++
		fmt.Fprintf(buf, " subject_id = $%d", count)
		args = append(args, subjectID)
	case accessToken != "":
		count++
		fmt.Fprintf(buf, " access_token = $%d", count)
		args = append(args, accessToken)
	case refreshToken != "":
		count++
		fmt.Fprintf(buf, " refresh_token = $%d", count)
		args = append(args, refreshToken)
	}
	if expiredAtFrom != nil {
		if buf.Len() > 0 {
			fmt.Fprint(buf, " AND")
		}
		count++
		fmt.Fprintf(buf, " expire_at > $%d", count)
		args = append(args, expiredAtFrom)
	}
	if expiredAtTo != nil {
		if buf.Len() > 0 {
			fmt.Fprint(buf, " AND")
		}
		count++
		fmt.Fprintf(buf, " expire_at < $%d", count)
		args = append(args, expiredAtTo)
	}

	return buf, args
}

// Collect implements prometheus Collector interface.
func (c *Storage) Collect(in chan<- prometheus.Metric) {
	c.queriesTotal.Collect(in)
	c.queriesDuration.Collect(in)
	c.errors.Collect(in)
}

// Describe implements prometheus Collector interface.
func (c *Storage) Describe(in chan<- *prometheus.Desc) {
	c.queriesTotal.Describe(in)
	c.queriesDuration.Describe(in)
	c.errors.Describe(in)
}

type sessionEntity struct {
	AccessToken   string    `json:"accessToken"`
	RefreshToken  string    `json:"refreshToken"`
	SubjectID     string    `json:"subjectId"`
	SubjectClient string    `json:"subjectClient"`
	Bag           model.Bag `json:"bag"`
	ExpireAt      time.Time `json:"expireAt"`
}

func (se *sessionEntity) session() (*mnemosynerpc.Session, error) {
	expireAt, err := ptypes.TimestampProto(se.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.Session{
		AccessToken:   se.AccessToken,
		RefreshToken:  se.RefreshToken,
		SubjectId:     se.SubjectID,
		SubjectClient: se.SubjectClient,
		Bag:           se.Bag,
		ExpireAt:      expireAt,
	}, nil
}
