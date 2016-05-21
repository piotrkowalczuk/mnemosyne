package mnemosyned

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne"
)

var (
	tmpKey = []byte(hex.EncodeToString([]byte("1")))
)

type postgresStorage struct {
	db                                             *sql.DB
	table                                          string
	ttl                                            time.Duration
	generator                                      mnemosyne.RandomBytesGenerator
	monitor                                        *monitoring
	querySave, queryGet, queryExists, queryAbandon string
}

func newPostgresStorage(tb string, db *sql.DB, m *monitoring, ttl time.Duration) Storage {
	return &postgresStorage{
		db:        db,
		table:     tb,
		ttl:       ttl,
		generator: &mnemosyne.SystemRandomBytesGenerator{},
		monitor:   m,
		querySave: `INSERT INTO mnemosyne.` + tb + ` (access_token, subject_id, bag)
		VALUES ($1, $2, $3)
		RETURNING expire_at`,
		queryGet: `SELECT subject_id, bag, expire_at
		FROM mnemosyne.` + tb + `
		WHERE access_token = $1
		LIMIT 1`,
		queryExists:  `SELECT EXISTS(SELECT 1 FROM mnemosyne.` + tb + ` WHERE access_token = $1)`,
		queryAbandon: `DELETE FROM mnemosyne.` + tb + ` WHERE access_token = $1`,
	}
}

// Create implements Storage interface.
func (ps *postgresStorage) Start(subjectID string, bag map[string]string) (*mnemosyne.Session, error) {
	accessToken, err := mnemosyne.RandomAccessToken(ps.generator, tmpKey)
	if err != nil {
		return nil, err
	}

	ent := &sessionEntity{
		AccessToken: accessToken,
		SubjectID:   subjectID,
		Bag:         bagpack(bag),
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
		entity.Bag,
	).Scan(
		&entity.ExpireAt,
	)
	ps.incQueries(field)

	return
}

// Get implements Storage interface.
func (ps *postgresStorage) Get(accessToken *mnemosyne.AccessToken) (*mnemosyne.Session, error) {
	var entity sessionEntity
	field := metrics.Field{Key: "query", Value: "get"}

	err := ps.db.QueryRow(ps.queryGet, accessToken).Scan(
		&entity.SubjectID,
		&entity.Bag,
		&entity.ExpireAt,
	)
	if err != nil {
		ps.incError(field)
		if err == sql.ErrNoRows {
			return nil, SessionNotFound
		}
		return nil, err
	}

	expireAt, err := ptypes.TimestampProto(entity.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &mnemosyne.Session{
		AccessToken: accessToken,
		SubjectId:   entity.SubjectID,
		Bag:         entity.Bag,
		ExpireAt:    expireAt,
	}, nil
}

// List implements Storage interface.
func (ps *postgresStorage) List(offset, limit int64, expiredAtFrom, expiredAtTo *time.Time) ([]*mnemosyne.Session, error) {
	if limit == 0 {
		return nil, errors.New("cannot retrieve list of sessions, limit needs to be higher than 0")
	}

	args := []interface{}{offset, limit}
	query := "SELECT access_token, subject_id, bag, expire_at FROM mnemosyne." + ps.table + " "

	switch {
	case expiredAtFrom != nil && expiredAtTo == nil:
		query += "expire_at > $3"
		args = append(args, expiredAtFrom)
	case expiredAtFrom == nil && expiredAtTo != nil:
		query += "expire_at < $3"
		args = append(args, expiredAtTo)
	case expiredAtFrom != nil && expiredAtTo != nil:
		query += "expire_at > $4 AND expire_at < $5"
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

	sessions := make([]*mnemosyne.Session, 0, limit)
	for rows.Next() {
		var entity sessionEntity

		err = rows.Scan(
			&entity.AccessToken,
			&entity.SubjectID,
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
		sessions = append(sessions, &mnemosyne.Session{
			AccessToken: &entity.AccessToken,
			SubjectId:   entity.SubjectID,
			Bag:         entity.Bag,
			ExpireAt:    expireAt,
		})
	}
	if rows.Err() != nil {
		ps.incError(field)
		return nil, rows.Err()
	}

	return sessions, nil
}

// Exists implements Storage interface.
func (ps *postgresStorage) Exists(accessToken *mnemosyne.AccessToken) (exists bool, err error) {
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
func (ps *postgresStorage) Abandon(accessToken *mnemosyne.AccessToken) (bool, error) {
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
		return false, SessionNotFound
	}

	return true, nil
}

// SetData implements Storage interface.
func (ps *postgresStorage) SetValue(accessToken *mnemosyne.AccessToken, key, value string) (map[string]string, error) {
	var err error
	if accessToken == nil {
		return nil, mnemosyne.ErrMissingAccessToken
	}

	entity := &sessionEntity{
		AccessToken: *accessToken,
	}
	selectQuery := `
		SELECT subject_id, bag, expire_at
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
		&entity.SubjectID,
		&entity.Bag,
		&entity.ExpireAt,
	)
	if err != nil {
		ps.incError(metrics.Field{Key: "query", Value: "set_value_select"})
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, SessionNotFound
		}
		return nil, err
	}
	ps.incQueries(metrics.Field{Key: "query", Value: "set_value_select"})

	entity.Bag.Set(key, value)

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
func (ps *postgresStorage) Delete(accessToken *mnemosyne.AccessToken, expiredAtFrom, expiredAtTo *time.Time) (int64, error) {
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

func (ps *postgresStorage) where(accessToken *mnemosyne.AccessToken, expiredAtFrom, expiredAtTo *time.Time) (string, []interface{}) {
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
	AccessToken mnemosyne.AccessToken `json:"accessToken"`
	SubjectID   string                `json:"subjectId"`
	Bag         bagpack               `json:"bag"`
	ExpireAt    time.Time             `json:"expireAt"`
}

func (se *sessionEntity) session() (*mnemosyne.Session, error) {
	expireAt, err := ptypes.TimestampProto(se.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &mnemosyne.Session{
		AccessToken: &se.AccessToken,
		SubjectId:   se.SubjectID,
		Bag:         se.Bag,
		ExpireAt:    expireAt,
	}, nil
}
