package main

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-kit/kit/metrics"
	"github.com/piotrkowalczuk/mnemosyne"
)

const (
	postgresTableNamePlaceholder = "%%TABLE_NAME%%"
)

type postgresStorage struct {
	db        *sql.DB
	tableName string
	generator mnemosyne.RandomBytesGenerator
	monitor   *monitoring
}

func newPostgresStorage(tn string, db *sql.DB, m *monitoring) Storage {
	return &postgresStorage{
		db:        db,
		tableName: tn,
		generator: &mnemosyne.SystemRandomBytesGenerator{},
		monitor:   m,
	}
}

func initPostgresStorage(tn string, db *sql.DB, m *monitoring) func() (Storage, error) {
	return func() (Storage, error) {
		return newPostgresStorage(tn, db, m), nil
	}
}

// Create ...
func (ps *postgresStorage) Create(data map[string]string) (*mnemosyne.Session, error) {
	token, err := mnemosyne.NewTokenRandom(ps.generator, "1")
	if err != nil {
		return nil, err
	}

	entity := &SessionEntity{
		Token: token,
		Data:  Data(data),
	}

	if err := ps.save(entity); err != nil {
		return nil, err
	}

	return newSessionFromSessionEntity(entity), nil
}

func (ps *postgresStorage) save(entity *SessionEntity) error {
	query := `
		INSERT INTO ` + ps.tableName + ` (token, data, expire_at)
		VALUES ($1, $2, NOW() + '30 minutes'::interval)
		RETURNING expire_at

	`
	field := metrics.Field{Key: "query", Value: query}

	encodedData, err := entity.Data.EncodeToJSON()
	if err != nil {
		ps.monitor.postgres.errors.With(field).Add(1)
		return err
	}

	err = ps.db.QueryRow(
		query,
		entity.Token,
		encodedData,
	).Scan(
		&entity.ExpireAt,
	)
	ps.monitor.postgres.queries.With(field).Add(1)

	return err
}

// Get ...
func (ps *postgresStorage) Get(token *mnemosyne.Token) (*mnemosyne.Session, error) {
	var data Data
	var expireAt time.Time
	query := `
		SELECT data, expire_at
		FROM ` + ps.tableName + `
		WHERE token = $1
		LIMIT 1
	`
	field := metrics.Field{Key: "query", Value: query}

	err := ps.db.QueryRow(query, token).Scan(
		&data,
		&expireAt,
	)
	if err != nil {
		ps.monitor.postgres.errors.With(field).Add(1)
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return nil, err
	}

	return &mnemosyne.Session{
		Token:    token,
		Data:     data,
		ExpireAt: mnemosyne.TimeToTimestamp(expireAt),
	}, nil
}

// List satisfy Storage interface.
func (ps *postgresStorage) List(offset, limit int64, expiredAtFrom, expiredAtTo *time.Time) ([]*mnemosyne.Session, error) {
	if limit == 0 {
		return nil, errors.New("mnemosyned: cannot retrieve list of sessions, limit needs to be higher than 0")
	}

	args := []interface{}{offset, limit}
	query := "SELECT token, data, expire_at FROM " + ps.tableName

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
		ps.monitor.postgres.errors.With(field).Add(1)
		return nil, err
	}
	ps.monitor.postgres.queries.With(field).Add(1)

	sessions := make([]*mnemosyne.Session, 0, limit)
	for rows.Next() {
		var token mnemosyne.Token
		var data Data
		var expireAt time.Time

		err = rows.Scan(
			&token,
			&data,
			&expireAt,
		)
		if err != nil {
			ps.monitor.postgres.errors.With(field).Add(1)
			return nil, err
		}

		sessions = append(sessions, &mnemosyne.Session{
			Token:    &token,
			Data:     data,
			ExpireAt: mnemosyne.TimeToTimestamp(expireAt),
		})
	}
	if rows.Err() != nil {
		ps.monitor.postgres.errors.With(field).Add(1)
		return nil, rows.Err()
	}

	return sessions, nil
}

// Exists ...
func (ps *postgresStorage) Exists(token *mnemosyne.Token) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM ` + ps.tableName + ` WHERE token = $1)`
	field := metrics.Field{Key: "query", Value: query}

	err = ps.db.QueryRow(query, token).Scan(
		&exists,
	)
	if err != nil {
		ps.monitor.postgres.errors.With(field).Add(1)
	}
	ps.monitor.postgres.queries.With(field).Add(1)

	return
}

// Abandon ...
func (ps *postgresStorage) Abandon(token *mnemosyne.Token) (bool, error) {
	query := `DELETE FROM ` + ps.tableName + ` WHERE token = $1`
	field := metrics.Field{Key: "query", Value: query}

	result, err := ps.db.Exec(query, token)
	if err != nil {
		ps.monitor.postgres.errors.With(field).Add(1)
		return false, err
	}

	ps.monitor.postgres.queries.With(field).Add(1)

	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if affected == 0 {
		return false, errSessionNotFound
	}

	return true, nil
}

// SetData ...
func (ps *postgresStorage) SetData(token *mnemosyne.Token, key, value string) (*mnemosyne.Session, error) {
	var dataEncoded []byte
	var err error

	entity := &SessionEntity{
		Token: token,
	}
	selectQuery := `
		SELECT data, expire_at
		FROM ` + ps.tableName + `
		WHERE token = $1
		FOR UPDATE
	`
	updateQuery := `
		UPDATE ` + ps.tableName + `
		SET
			data = $2
		WHERE token = $1
	`

	tx, err := ps.db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.QueryRow(selectQuery, token).Scan(
		&dataEncoded,
		&entity.ExpireAt,
	)
	if err != nil {
		ps.monitor.postgres.errors.With(metrics.Field{Key: "query", Value: selectQuery}).Add(1)
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return nil, err
	}
	ps.monitor.postgres.queries.With(metrics.Field{Key: "query", Value: selectQuery}).Add(1)

	entity.Data, err = DecodeSessionDataFromJSON(dataEncoded)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	entity.Data.Set(key, value)

	dataEncoded, err = entity.Data.EncodeToJSON()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(updateQuery, token, dataEncoded)
	if err != nil {
		ps.monitor.postgres.errors.With(metrics.Field{Key: "query", Value: updateQuery}).Add(1)
		tx.Rollback()
		return nil, err
	}
	ps.monitor.postgres.queries.With(metrics.Field{Key: "query", Value: updateQuery}).Add(1)

	tx.Commit()

	return newSessionFromSessionEntity(entity), nil
}

// Delete
func (ps *postgresStorage) Delete(token *mnemosyne.Token, expiredAtFrom, expiredAtTo *time.Time) (int64, error) {
	if token == nil && expiredAtFrom == nil && expiredAtTo == nil {
		return 0, errors.New("mnemosyned: session cannot be deleted, no where parameter provided")
	}

	where, args := ps.where(token, expiredAtFrom, expiredAtTo)
	query := "DELETE FROM " + ps.tableName + " WHERE " + where
	field := metrics.Field{Key: "query", Value: query}

	result, err := ps.db.Exec(query, args...)
	if err != nil {
		ps.monitor.postgres.errors.With(field).Add(1)
		return 0, err
	}
	ps.monitor.postgres.queries.With(field).Add(1)

	return result.RowsAffected()
}

// Setup prepares storage for usage.
func (ps *postgresStorage) Setup() error {
	sql := strings.Replace(`
		CREATE TABLE IF NOT EXISTS %%TABLE_NAME%% (
			token character varying(255) PRIMARY KEY,
			data json NOT NULL,
			expire_at timestamp with time zone NOT NULL
		)
    `, postgresTableNamePlaceholder, ps.tableName, -1)

	_, err := ps.db.Exec(sql)

	return err
}

// TearDown ...
func (ps *postgresStorage) TearDown() error {
	sql := strings.Replace(`DROP TABLE %%TABLE_NAME%%`, postgresTableNamePlaceholder, ps.tableName, -1)

	_, err := ps.db.Exec(sql)

	return err
}

func (ps *postgresStorage) where(token *mnemosyne.Token, expiredAtFrom, expiredAtTo *time.Time) (string, []interface{}) {
	switch {
	case token != nil && expiredAtFrom == nil && expiredAtTo == nil:
		return "token = $1", []interface{}{token}
	case token == nil && expiredAtFrom != nil && expiredAtTo == nil:
		return "expire_at > $1", []interface{}{expiredAtFrom}
	case token == nil && expiredAtFrom == nil && expiredAtTo != nil:
		return "expire_at < $1", []interface{}{expiredAtTo}
	case token != nil && expiredAtFrom != nil && expiredAtTo == nil:
		return "token = $1 AND expire_at > $2", []interface{}{token, expiredAtFrom}
	case token != nil && expiredAtFrom == nil && expiredAtTo != nil:
		return "token = $1 AND expire_at < $2", []interface{}{token, expiredAtTo}
	case token == nil && expiredAtFrom != nil && expiredAtTo != nil:
		return "expire_at > $1 AND expire_at < $2", []interface{}{expiredAtFrom, expiredAtTo}
	case token != nil && expiredAtFrom != nil && expiredAtTo != nil:
		return "token = $1 AND expire_at > $2 AND expire_at < $3", []interface{}{token, expiredAtFrom, expiredAtTo}
	default:
		return "", nil
	}
}
