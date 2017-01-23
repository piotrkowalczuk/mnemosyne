package mnemosyned

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/prometheus/client_golang/prometheus"
)

type postgresStorage struct {
	db                                             *sql.DB
	schema                                         string
	table                                          string
	ttl                                            time.Duration
	monitor                                        *monitoring
	querySave, queryGet, queryExists, queryAbandon string
}

func newPostgresStorage(tb, schema string, db *sql.DB, m *monitoring, ttl time.Duration) storage {
	return &postgresStorage{
		db:      db,
		table:   tb,
		schema:  schema,
		ttl:     ttl,
		monitor: m,
		querySave: `INSERT INTO ` + schema + ` .` + tb + ` (access_token, subject_id, subject_client, bag)
			VALUES ($1, $2, $3, $4)
			RETURNING expire_at`,
		queryGet: fmt.Sprintf(`UPDATE `+schema+` .`+tb+`
			SET expire_at = (NOW() + '%d seconds')
			WHERE access_token = $1
			RETURNING subject_id, subject_client, bag, expire_at`, int64(ttl.Seconds())),
		queryExists:  `SELECT EXISTS(SELECT 1 FROM ` + schema + ` .` + tb + ` WHERE access_token = $1)`,
		queryAbandon: `DELETE FROM ` + schema + ` .` + tb + ` WHERE access_token = $1`,
	}
}

// Start implements storage interface.
func (ps *postgresStorage) Start(ctx context.Context, at, sid, sc string, b map[string]string) (*mnemosynerpc.Session, error) {
	ent := &sessionEntity{
		AccessToken:   at,
		SubjectID:     sid,
		SubjectClient: sc,
		Bag:           bag(b),
	}

	if err := ps.save(ctx, ent); err != nil {
		return nil, err
	}

	return ent.session()
}

func (ps *postgresStorage) save(ctx context.Context, entity *sessionEntity) (err error) {
	labels := prometheus.Labels{"query": "save"}
	err = ps.db.QueryRowContext(
		ctx,
		ps.querySave,
		entity.AccessToken,
		entity.SubjectID,
		entity.SubjectClient,
		entity.Bag,
	).Scan(
		&entity.ExpireAt,
	)
	ps.incQueries(labels)
	if err != nil {
		ps.incError(labels)
	}
	return
}

// Get implements storage interface.
func (ps *postgresStorage) Get(ctx context.Context, accessToken string) (*mnemosynerpc.Session, error) {
	var entity sessionEntity
	labels := prometheus.Labels{"query": "get"}

	err := ps.db.QueryRowContext(ctx, ps.queryGet, accessToken).Scan(
		&entity.SubjectID,
		&entity.SubjectClient,
		&entity.Bag,
		&entity.ExpireAt,
	)
	ps.incQueries(labels)
	if err != nil {
		ps.incError(labels)
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return nil, err
	}

	expireAt, err := ptypes.TimestampProto(entity.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.Session{
		AccessToken:   string(accessToken),
		SubjectId:     entity.SubjectID,
		SubjectClient: entity.SubjectClient,
		Bag:           entity.Bag,
		ExpireAt:      expireAt,
	}, nil
}

// List implements storage interface.
func (ps *postgresStorage) List(ctx context.Context, offset, limit int64, expiredAtFrom, expiredAtTo *time.Time) ([]*mnemosynerpc.Session, error) {
	if limit == 0 {
		return nil, errors.New("cannot retrieve list of sessions, limit needs to be higher than 0")
	}

	args := []interface{}{offset, limit}
	query := "SELECT access_token, subject_id, subject_client, bag, expire_at FROM " + ps.schema + "." + ps.table + " "
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

	rows, err := ps.db.QueryContext(ctx, query, args...)
	ps.incQueries(labels)
	if err != nil {
		ps.incError(labels)
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*mnemosynerpc.Session, 0, limit)
	for rows.Next() {
		var entity sessionEntity

		err = rows.Scan(
			&entity.AccessToken,
			&entity.SubjectID,
			&entity.SubjectClient,
			&entity.Bag,
			&entity.ExpireAt,
		)
		if err != nil {
			ps.incError(labels)
			return nil, err
		}

		expireAt, err := ptypes.TimestampProto(entity.ExpireAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &mnemosynerpc.Session{
			AccessToken:   entity.AccessToken,
			SubjectId:     entity.SubjectID,
			SubjectClient: entity.SubjectClient,
			Bag:           entity.Bag,
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
func (ps *postgresStorage) Exists(ctx context.Context, accessToken string) (exists bool, err error) {
	labels := prometheus.Labels{"query": "exists"}

	err = ps.db.QueryRowContext(ctx, ps.queryExists, accessToken).Scan(
		&exists,
	)
	if err != nil {
		ps.incError(labels)
	}
	ps.incQueries(labels)

	return
}

// Abandon implements storage interface.
func (ps *postgresStorage) Abandon(ctx context.Context, accessToken string) (bool, error) {
	labels := prometheus.Labels{"query": "abandon"}

	result, err := ps.db.ExecContext(ctx, ps.queryAbandon, accessToken)
	ps.incQueries(labels)
	if err != nil {
		ps.incError(labels)
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected == 0 {
		return false, errSessionNotFound
	}

	return true, nil
}

// SetValue implements storage interface.
func (ps *postgresStorage) SetValue(ctx context.Context, accessToken string, key, value string) (map[string]string, error) {
	var err error
	if accessToken == "" {
		return nil, errMissingAccessToken
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

	err = tx.QueryRowContext(ctx, selectQuery, accessToken).Scan(
		&entity.Bag,
	)
	ps.incQueries(prometheus.Labels{"query": "set_value_select"})
	if err != nil {
		ps.incError(prometheus.Labels{"query": "set_value_select"})
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return nil, err
	}

	entity.Bag.set(key, value)

	_, err = tx.ExecContext(ctx, updateQuery, accessToken, entity.Bag)
	ps.incQueries(prometheus.Labels{"query": "set_value_update"})
	if err != nil {
		ps.incError(prometheus.Labels{"query": "set_value_update"})
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return entity.Bag, nil
}

// Delete implements storage interface.
func (ps *postgresStorage) Delete(ctx context.Context, accessToken string, expiredAtFrom, expiredAtTo *time.Time) (int64, error) {
	if accessToken == "" && expiredAtFrom == nil && expiredAtTo == nil {
		return 0, errors.New("session cannot be deleted, no where parameter provided")
	}

	where, args := ps.where(accessToken, expiredAtFrom, expiredAtTo)
	query := "DELETE FROM " + ps.schema + "." + ps.table + " WHERE " + where
	labels := prometheus.Labels{"query": "delete"}

	result, err := ps.db.Exec(query, args...)
	ps.incQueries(labels)
	if err != nil {
		ps.incError(labels)
		return 0, err
	}

	return result.RowsAffected()
}

// Setup implements storage interface.
func (ps *postgresStorage) Setup() error {
	query := fmt.Sprintf(`
		CREATE SCHEMA IF NOT EXISTS %s;
		CREATE TABLE IF NOT EXISTS %s.%s (
			access_token BYTEA PRIMARY KEY,
			subject_id TEXT NOT NULL,
			subject_client TEXT,
			bag bytea NOT NULL,
			expire_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + '%d seconds')

		);
		CREATE INDEX ON %s.%s (subject_id);
		CREATE INDEX ON %s.%s (expire_at DESC);
	`, ps.schema, ps.schema, ps.table, int64(ps.ttl.Seconds()), ps.schema, ps.table, ps.schema, ps.table)
	_, err := ps.db.Exec(query)

	return err
}

// TearDown implements storage interface.
func (ps *postgresStorage) TearDown() error {
	_, err := ps.db.Exec(`DROP SCHEMA IF EXISTS ` + ps.schema + ` CASCADE`)

	return err
}

func (ps *postgresStorage) incQueries(field prometheus.Labels) {
	if ps.monitor.postgres.enabled {
		ps.monitor.postgres.queries.With(field).Inc()
	}
}

func (ps *postgresStorage) incError(field prometheus.Labels) {
	if ps.monitor.postgres.enabled {
		ps.monitor.postgres.errors.With(field).Inc()
	}
}

func (ps *postgresStorage) where(accessToken string, expiredAtFrom, expiredAtTo *time.Time) (string, []interface{}) {
	switch {
	case accessToken != "" && expiredAtFrom == nil && expiredAtTo == nil:
		return " access_token = $1", []interface{}{accessToken}
	case accessToken == "" && expiredAtFrom != nil && expiredAtTo == nil:
		return " expire_at > $1", []interface{}{expiredAtFrom}
	case accessToken == "" && expiredAtFrom == nil && expiredAtTo != nil:
		return " expire_at < $1", []interface{}{expiredAtTo}
	case accessToken != "" && expiredAtFrom != nil && expiredAtTo == nil:
		return " access_token = $1 AND expire_at > $2", []interface{}{accessToken, expiredAtFrom}
	case accessToken != "" && expiredAtFrom == nil && expiredAtTo != nil:
		return " access_token = $1 AND expire_at < $2", []interface{}{accessToken, expiredAtTo}
	case accessToken == "" && expiredAtFrom != nil && expiredAtTo != nil:
		return " expire_at > $1 AND expire_at < $2", []interface{}{expiredAtFrom, expiredAtTo}
	case accessToken != "" && expiredAtFrom != nil && expiredAtTo != nil:
		return " access_token = $1 AND expire_at > $2 AND expire_at < $3", []interface{}{accessToken, expiredAtFrom, expiredAtTo}
	default:
		return " ", nil
	}
}

type sessionEntity struct {
	AccessToken   string    `json:"accessToken"`
	SubjectID     string    `json:"subjectId"`
	SubjectClient string    `json:"subjectClient"`
	Bag           bag       `json:"bag"`
	ExpireAt      time.Time `json:"expireAt"`
}

func (se *sessionEntity) session() (*mnemosynerpc.Session, error) {
	expireAt, err := ptypes.TimestampProto(se.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.Session{
		AccessToken:   se.AccessToken,
		SubjectId:     se.SubjectID,
		SubjectClient: se.SubjectClient,
		Bag:           se.Bag,
		ExpireAt:      expireAt,
	}, nil
}
