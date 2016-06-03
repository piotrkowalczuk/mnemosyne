package mnemosyned

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
)

var (
	tmpKey = []byte(hex.EncodeToString([]byte("1")))
)

type postgresStorage struct {
	db                                             *sql.DB
	table                                          string
	ttl                                            time.Duration
	monitor                                        *monitoring
	querySave, queryGet, queryExists, queryAbandon string
}

func newPostgresStorage(tb string, db *sql.DB, m *monitoring, ttl time.Duration) Storage {
	return &postgresStorage{
		db:      db,
		table:   tb,
		ttl:     ttl,
		monitor: m,
		querySave: `INSERT INTO mnemosyne.` + tb + ` (access_token, subject_id, subject_client, bag)
			VALUES ($1, $2, $3, $4)
			RETURNING expire_at`,
		queryGet: fmt.Sprintf(`UPDATE mnemosyne.`+tb+`
			SET expire_at = (NOW() + '%d seconds')
			WHERE access_token = $1
			RETURNING subject_id, subject_client, bag, expire_at`, int64(ttl.Seconds())),
		queryExists:  `SELECT EXISTS(SELECT 1 FROM mnemosyne.` + tb + ` WHERE access_token = $1)`,
		queryAbandon: `DELETE FROM mnemosyne.` + tb + ` WHERE access_token = $1`,
	}
}

// Create implements Storage interface.
func (ps *postgresStorage) Start(sid, sc string, b map[string]string) (*mnemosynerpc.Session, error) {
	accessToken, err := mnemosynerpc.RandomAccessToken(tmpKey)
	if err != nil {
		return nil, err
	}

	ent := &sessionEntity{
		AccessToken:   accessToken,
		SubjectID:     sid,
		SubjectClient: sc,
		Bag:           bag(b),
	}

	if err := ps.save(ent); err != nil {
		return nil, err
	}

	return ent.session()
}

func (ps *postgresStorage) save(entity *sessionEntity) (err error) {
	field := metrics.Field{Key: "query", Value: "save"}

	err = ps.db.QueryRow(
		ps.querySave,
		entity.AccessToken,
		entity.SubjectID,
		entity.SubjectClient,
		entity.Bag,
	).Scan(
		&entity.ExpireAt,
	)
	ps.incQueries(field)

	return
}

// Get implements Storage interface.
func (ps *postgresStorage) Get(accessToken *mnemosynerpc.AccessToken) (*mnemosynerpc.Session, error) {
	var entity sessionEntity
	field := metrics.Field{Key: "query", Value: "get"}

	err := ps.db.QueryRow(ps.queryGet, accessToken).Scan(
		&entity.SubjectID,
		&entity.SubjectClient,
		&entity.Bag,
		&entity.ExpireAt,
	)
	if err != nil {
		ps.incError(field)
		if err == sql.ErrNoRows {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	expireAt, err := ptypes.TimestampProto(entity.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.Session{
		AccessToken:   accessToken,
		SubjectId:     entity.SubjectID,
		SubjectClient: entity.SubjectClient,
		Bag:           entity.Bag,
		ExpireAt:      expireAt,
	}, nil
}

// List implements Storage interface.
func (ps *postgresStorage) List(offset, limit int64, expiredAtFrom, expiredAtTo *time.Time) ([]*mnemosynerpc.Session, error) {
	if limit == 0 {
		return nil, errors.New("cannot retrieve list of sessions, limit needs to be higher than 0")
	}

	args := []interface{}{offset, limit}
	query := "SELECT access_token, subject_id, subject_client, bag, expire_at FROM mnemosyne." + ps.table + " "
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

	field := metrics.Field{Key: "query", Value: "list"}

	rows, err := ps.db.Query(query, args...)
	if err != nil {
		ps.incError(field)
		return nil, err
	}
	defer rows.Close()

	ps.incQueries(field)

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
			ps.incError(field)
			return nil, err
		}

		expireAt, err := ptypes.TimestampProto(entity.ExpireAt)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &mnemosynerpc.Session{
			AccessToken:   &entity.AccessToken,
			SubjectId:     entity.SubjectID,
			SubjectClient: entity.SubjectClient,
			Bag:           entity.Bag,
			ExpireAt:      expireAt,
		})
	}
	if rows.Err() != nil {
		ps.incError(field)
		return nil, rows.Err()
	}

	return sessions, nil
}

// Exists implements Storage interface.
func (ps *postgresStorage) Exists(accessToken *mnemosynerpc.AccessToken) (exists bool, err error) {
	field := metrics.Field{Key: "query", Value: "exists"}

	err = ps.db.QueryRow(ps.queryExists, *accessToken).Scan(
		&exists,
	)
	if err != nil {
		ps.incError(field)
	}
	ps.incQueries(field)

	return
}

// Abandon ...
func (ps *postgresStorage) Abandon(accessToken *mnemosynerpc.AccessToken) (bool, error) {
	field := metrics.Field{Key: "query", Value: "abandon"}

	result, err := ps.db.Exec(ps.queryAbandon, *accessToken)
	if err != nil {
		ps.monitor.postgres.errors.With(field).Add(1)
		return false, err
	}

	ps.incQueries(field)

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected == 0 {
		return false, ErrSessionNotFound
	}

	return true, nil
}

// SetData implements Storage interface.
func (ps *postgresStorage) SetValue(accessToken *mnemosynerpc.AccessToken, key, value string) (map[string]string, error) {
	var err error
	if accessToken == nil {
		return nil, ErrMissingAccessToken
	}

	entity := &sessionEntity{
		AccessToken: *accessToken,
	}
	selectQuery := `
		SELECT bag
		FROM mnemosyne.` + ps.table + `
		WHERE access_token = $1
		FOR UPDATE
	`
	updateQuery := `
		UPDATE mnemosyne.` + ps.table + `
		SET
			bag = $2
		WHERE access_token = $1
	`

	tx, err := ps.db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.QueryRow(selectQuery, *accessToken).Scan(
		&entity.Bag,
	)
	if err != nil {
		ps.incError(metrics.Field{Key: "query", Value: "set_value_select"})
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	ps.incQueries(metrics.Field{Key: "query", Value: "set_value_select"})

	entity.Bag.set(key, value)

	_, err = tx.Exec(updateQuery, *accessToken, entity.Bag)
	if err != nil {
		ps.incError(metrics.Field{Key: "query", Value: "set_value_update"})
		tx.Rollback()
		return nil, err
	}
	ps.incQueries(metrics.Field{Key: "query", Value: "set_value_update"})

	tx.Commit()

	return entity.Bag, nil
}

// Delete implements Storage interface.
func (ps *postgresStorage) Delete(accessToken *mnemosynerpc.AccessToken, expiredAtFrom, expiredAtTo *time.Time) (int64, error) {
	if accessToken == nil && expiredAtFrom == nil && expiredAtTo == nil {
		return 0, errors.New("session cannot be deleted, no where parameter provided")
	}

	where, args := ps.where(accessToken, expiredAtFrom, expiredAtTo)
	query := "DELETE FROM mnemosyne." + ps.table + " WHERE " + where
	field := metrics.Field{Key: "query", Value: "delete"}

	result, err := ps.db.Exec(query, args...)
	if err != nil {
		ps.incError(field)
		return 0, err
	}
	ps.incQueries(field)

	return result.RowsAffected()
}

// Setup implements Storage interface.
func (ps *postgresStorage) Setup() error {
	_, err := ps.db.Exec(fmt.Sprintf(`
		CREATE SCHEMA IF NOT EXISTS mnemosyne;
		CREATE TABLE IF NOT EXISTS mnemosyne.%s (
			access_token BYTEA PRIMARY KEY,
			subject_id TEXT NOT NULL,
			subject_client TEXT,
			bag bytea NOT NULL,
			expire_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + '%d seconds')

		);
		CREATE INDEX ON mnemosyne.%s (subject_id);
		CREATE INDEX ON mnemosyne.%s (expire_at DESC);
	`, ps.table, int64(ps.ttl.Seconds()), ps.table, ps.table))

	return err
}

// TearDown implements Storage interface.
func (ps *postgresStorage) TearDown() error {
	_, err := ps.db.Exec(`DROP SCHEMA IF EXISTS mnemosyne CASCADE`)

	return err
}

func (ps *postgresStorage) incQueries(field metrics.Field) {
	if ps.monitor.postgres.enabled {
		ps.monitor.postgres.queries.With(field).Add(1)
	}
}

func (ps *postgresStorage) incError(field metrics.Field) {
	if ps.monitor.postgres.enabled {
		ps.monitor.postgres.errors.With(field).Add(1)
	}
}

func (ps *postgresStorage) where(accessToken *mnemosynerpc.AccessToken, expiredAtFrom, expiredAtTo *time.Time) (string, []interface{}) {
	switch {
	case accessToken != nil && expiredAtFrom == nil && expiredAtTo == nil:
		return " access_token = $1", []interface{}{accessToken}
	case accessToken == nil && expiredAtFrom != nil && expiredAtTo == nil:
		return " expire_at > $1", []interface{}{expiredAtFrom}
	case accessToken == nil && expiredAtFrom == nil && expiredAtTo != nil:
		return " expire_at < $1", []interface{}{expiredAtTo}
	case accessToken != nil && expiredAtFrom != nil && expiredAtTo == nil:
		return " access_token = $1 AND expire_at > $2", []interface{}{accessToken, expiredAtFrom}
	case accessToken != nil && expiredAtFrom == nil && expiredAtTo != nil:
		return " access_token = $1 AND expire_at < $2", []interface{}{accessToken, expiredAtTo}
	case accessToken == nil && expiredAtFrom != nil && expiredAtTo != nil:
		return " expire_at > $1 AND expire_at < $2", []interface{}{expiredAtFrom, expiredAtTo}
	case accessToken != nil && expiredAtFrom != nil && expiredAtTo != nil:
		return " access_token = $1 AND expire_at > $2 AND expire_at < $3", []interface{}{accessToken, expiredAtFrom, expiredAtTo}
	default:
		return " ", nil
	}
}

type sessionEntity struct {
	AccessToken   mnemosynerpc.AccessToken `json:"accessToken"`
	SubjectID     string                   `json:"subjectId"`
	SubjectClient string                   `json:"subjectClient"`
	Bag           bag                      `json:"bag"`
	ExpireAt      time.Time                `json:"expireAt"`
}

func (se *sessionEntity) session() (*mnemosynerpc.Session, error) {
	expireAt, err := ptypes.TimestampProto(se.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.Session{
		AccessToken:   &se.AccessToken,
		SubjectId:     se.SubjectID,
		SubjectClient: se.SubjectClient,
		Bag:           se.Bag,
		ExpireAt:      expireAt,
	}, nil
}
