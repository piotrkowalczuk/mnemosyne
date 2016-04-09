package mnemosyne

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/golang/protobuf/ptypes"
)

var (
	tmpKey = []byte(hex.EncodeToString([]byte("1")))
)

type postgresStorage struct {
	db        *sql.DB
	table     string
	generator RandomBytesGenerator
	monitor   *monitoring
}

func newPostgresStorage(tn string, db *sql.DB, m *monitoring) Storage {
	return &postgresStorage{
		db:        db,
		table:     tn,
		generator: &SystemRandomBytesGenerator{},
		monitor:   m,
	}
}

// Create implements Storage interface.
func (ps *postgresStorage) Start(subjectID string, bag map[string]string) (*Session, error) {
	accessToken, err := RandomAccessToken(ps.generator, tmpKey)
	if err != nil {
		return nil, err
	}

	entity := &sessionEntity{
		AccessToken: accessToken,
		SubjectID:   subjectID,
		Bag:         bagpack(bag),
	}

	if err := ps.save(entity); err != nil {
		return nil, err
	}

	return newSessionFromSessionEntity(entity)
}

func (ps *postgresStorage) save(entity *sessionEntity) (err error) {
	query := `
		INSERT INTO mnemosyne.` + ps.table + ` (access_token, subject_id, bag, expire_at)
		VALUES ($1, $2, $3, NOW() + '30 minutes'::interval)
		RETURNING expire_at

	`
	field := metrics.Field{Key: "query", Value: query}

	err = ps.db.QueryRow(
		query,
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
func (ps *postgresStorage) Get(accessToken *AccessToken) (*Session, error) {
	var entity sessionEntity
	query := `
		SELECT subject_id, bag, expire_at
		FROM mnemosyne.` + ps.table + `
		WHERE access_token = $1
		LIMIT 1
	`
	field := metrics.Field{Key: "query", Value: query}

	err := ps.db.QueryRow(query, accessToken).Scan(
		&entity.SubjectID,
		&entity.Bag,
		&entity.ExpireAt,
	)
	if err != nil {
		ps.incError(field)
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return nil, err
	}

	expireAt, err := ptypes.TimestampProto(entity.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &Session{
		AccessToken: accessToken,
		SubjectId:   entity.SubjectID,
		Bag:         entity.Bag,
		ExpireAt:    expireAt,
	}, nil
}

// List implements Storage interface.
func (ps *postgresStorage) List(offset, limit int64, expiredAtFrom, expiredAtTo *time.Time) ([]*Session, error) {
	if limit == 0 {
		return nil, errors.New("mnemosyne: cannot retrieve list of sessions, limit needs to be higher than 0")
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

	field := metrics.Field{Key: "query", Value: query}

	rows, err := ps.db.Query(query, args...)
	if err != nil {
		ps.incError(field)
		return nil, err
	}
	defer rows.Close()

	ps.incQueries(field)

	sessions := make([]*Session, 0, limit)
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
		sessions = append(sessions, &Session{
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
func (ps *postgresStorage) Exists(accessToken *AccessToken) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM mnemosyne.` + ps.table + ` WHERE access_token = $1)`
	field := metrics.Field{Key: "query", Value: query}

	err = ps.db.QueryRow(query, *accessToken).Scan(
		&exists,
	)
	if err != nil {
		ps.incError(field)
	}
	ps.incQueries(field)

	return
}

// Abandon ...
func (ps *postgresStorage) Abandon(accessToken *AccessToken) (bool, error) {
	query := `DELETE FROM mnemosyne.` + ps.table + ` WHERE access_token = $1`
	field := metrics.Field{Key: "query", Value: query}

	result, err := ps.db.Exec(query, *accessToken)
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
		return false, errSessionNotFound
	}

	return true, nil
}

// SetData implements Storage interface.
func (ps *postgresStorage) SetValue(accessToken *AccessToken, key, value string) (map[string]string, error) {
	var err error
	if accessToken == nil {
		return nil, ErrMissingAccessToken
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
		ps.incError(metrics.Field{Key: "query", Value: selectQuery})
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return nil, err
	}
	ps.incQueries(metrics.Field{Key: "query", Value: selectQuery})

	entity.Bag.Set(key, value)

	_, err = tx.Exec(updateQuery, *accessToken, entity.Bag)
	if err != nil {
		ps.incError(metrics.Field{Key: "query", Value: updateQuery})
		tx.Rollback()
		return nil, err
	}
	ps.incQueries(metrics.Field{Key: "query", Value: updateQuery})

	tx.Commit()

	return entity.Bag, nil
}

// Delete implements Storage interface.
func (ps *postgresStorage) Delete(accessToken *AccessToken, expiredAtFrom, expiredAtTo *time.Time) (int64, error) {
	if accessToken == nil && expiredAtFrom == nil && expiredAtTo == nil {
		return 0, errors.New("mnemosyne: session cannot be deleted, no where parameter provided")
	}

	where, args := ps.where(accessToken, expiredAtFrom, expiredAtTo)
	query := "DELETE FROM mnemosyne." + ps.table + " WHERE " + where
	field := metrics.Field{Key: "query", Value: query}

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
	_, err := ps.db.Exec(`
		CREATE SCHEMA IF NOT EXISTS mnemosyne;
		CREATE TABLE IF NOT EXISTS mnemosyne.` + ps.table + ` (
			access_token BYTEA PRIMARY KEY,
			subject_id TEXT NOT NULL,
			bag bytea NOT NULL,
			expire_at TIMESTAMPTZ NOT NULL
		);
	`)

	return err
}

// TearDown implements Storage interface.
func (ps *postgresStorage) TearDown() error {
	_, err := ps.db.Exec(`DROP SCHEMA mnemosyne CASCADE`)

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

func (ps *postgresStorage) where(accessToken *AccessToken, expiredAtFrom, expiredAtTo *time.Time) (string, []interface{}) {
	switch {
	case accessToken != nil && expiredAtFrom == nil && expiredAtTo == nil:
		return "access_token = $1", []interface{}{accessToken}
	case accessToken == nil && expiredAtFrom != nil && expiredAtTo == nil:
		return "expire_at > $1", []interface{}{expiredAtFrom}
	case accessToken == nil && expiredAtFrom == nil && expiredAtTo != nil:
		return "expire_at < $1", []interface{}{expiredAtTo}
	case accessToken != nil && expiredAtFrom != nil && expiredAtTo == nil:
		return "access_token = $1 AND expire_at > $2", []interface{}{accessToken, expiredAtFrom}
	case accessToken != nil && expiredAtFrom == nil && expiredAtTo != nil:
		return "access_token = $1 AND expire_at < $2", []interface{}{accessToken, expiredAtTo}
	case accessToken == nil && expiredAtFrom != nil && expiredAtTo != nil:
		return "expire_at > $1 AND expire_at < $2", []interface{}{expiredAtFrom, expiredAtTo}
	case accessToken != nil && expiredAtFrom != nil && expiredAtTo != nil:
		return "access_token = $1 AND expire_at > $2 AND expire_at < $3", []interface{}{accessToken, expiredAtFrom, expiredAtTo}
	default:
		return "", nil
	}
}

type sessionEntity struct {
	AccessToken AccessToken `json:"accessToken"`
	SubjectID   string      `json:"subjectId"`
	Bag         bagpack     `json:"bag"`
	ExpireAt    time.Time   `json:"expireAt"`
}

func newSessionFromSessionEntity(entity *sessionEntity) (*Session, error) {
	expireAt, err := ptypes.TimestampProto(entity.ExpireAt)
	if err != nil {
		return nil, err
	}
	return &Session{
		AccessToken: &entity.AccessToken,
		SubjectId:   entity.SubjectID,
		Bag:         entity.Bag,
		ExpireAt:    expireAt,
	}, nil
}
