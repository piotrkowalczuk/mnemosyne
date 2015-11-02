package main

import (
	"database/sql"
	"encoding/hex"
	"strings"
	"time"

	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/crypto/sha3"
)

const (
	postgresTableNamePlaceholder = "%%TABLE_NAME%%"
)

type postgresStorage struct {
	db        *sql.DB
	tableName string
	generator RandomBytesGenerator
}

func newPostgresStorage(db *sql.DB, tableName string) Storage {
	return &postgresStorage{
		db:        db,
		tableName: tableName,
		generator: &SystemRandomBytesGenerator{},
	}
}

// Create ...
func (ps *postgresStorage) Create(data map[string]string) (*mnemosyne.Session, error) {
	buf, err := ps.generator.GenerateRandomBytes(128)
	if err != nil {
		return nil, err
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	id := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(id, buf)

	entity := &SessionEntity{
		ID: &mnemosyne.ID{
			Key:  "1", // TODO: implement partitioning
			Hash: hex.EncodeToString(id),
		},
		Data: Data(data),
	}

	if err := ps.save(entity); err != nil {
		return nil, err
	}

	return newSessionFromSessionEntity(entity), nil
}

func (ps *postgresStorage) save(entity *SessionEntity) error {
	query := `
		INSERT INTO ` + ps.tableName + ` (id, data, expire_at)
		VALUES ($1, $2, NOW() + '30 minutes'::interval)
		RETURNING expire_at

	`

	encodedData, err := entity.Data.EncodeToJSON()
	if err != nil {
		return err
	}

	err = ps.db.QueryRow(
		query,
		entity.ID.Hash,
		encodedData,
	).Scan(
		&entity.ExpireAt,
	)

	return err
}

// Get ...
func (ps *postgresStorage) Get(id *mnemosyne.ID) (*mnemosyne.Session, error) {
	var data Data
	var expireAt time.Time
	query := `
		SELECT data, expire_at
		FROM ` + ps.tableName + `
		WHERE id = $1
		LIMIT 1
	`

	err := ps.db.QueryRow(query, id.Hash).Scan(
		&data,
		&expireAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return nil, err
	}

	return &mnemosyne.Session{
		Id:       id,
		Data:     data,
		ExpireAt: expireAt.Format(time.RFC3339),
	}, nil
}

// List ...
func (ps *postgresStorage) List(offset, limit int64, expiredAtFrom, expiredAtTo *time.Time) (*mnemosyne.Session, error) {
	return nil, nil
}

// Exists ...
func (ps *postgresStorage) Exists(id *mnemosyne.ID) (exists bool, err error) {
	query := `SELECT EXISTS(SELECT 1 FROM ` + ps.tableName + ` WHERE id = $1)`

	err = ps.db.QueryRow(query, id.Hash).Scan(
		&exists,
	)

	return
}

// Abandon ...
func (ps *postgresStorage) Abandon(id *mnemosyne.ID) (bool, error) {
	query := `DELETE FROM ` + ps.tableName + ` WHERE id = $1`

	result, err := ps.db.Exec(query, id.Hash)
	if err != nil {
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

// SetData ...
func (ps *postgresStorage) SetData(id *mnemosyne.ID, key, value string) (*mnemosyne.Session, error) {
	var dataEncoded []byte
	var err error

	entity := &SessionEntity{
		ID: id,
	}
	selectQuery := `
		SELECT data, expire_at
		FROM ` + ps.tableName + `
		WHERE id = $1
		FOR UPDATE
	`
	updateQuery := `
		UPDATE ` + ps.tableName + `
		SET
			data = $2
		WHERE id = $1
	`

	tx, err := ps.db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.QueryRow(selectQuery, id.Hash).Scan(
		&dataEncoded,
		&entity.ExpireAt,
	)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, errSessionNotFound
		}
		return nil, err
	}

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

	_, err = tx.Exec(
		updateQuery,
		id.Hash,
		dataEncoded,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return newSessionFromSessionEntity(entity), nil
}

// Delete
// TODO: implement properly, works partially
func (ps *postgresStorage) Delete(id *mnemosyne.ID, expiredAtFrom, expiredAtTo *time.Time) (int64, error) {
	var err error
	var result sql.Result
	delUntil := `
		DELETE FROM ` + ps.tableName + `
		WHERE expire_at < $1
	`
	delAll := `TRUNCATE TABLE ` + ps.tableName

	if expiredAtFrom == nil {
		result, err = ps.db.Exec(delAll)
	} else {
		result, err = ps.db.Exec(delUntil, expiredAtFrom)
	}

	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Setup prepares storage for usage.
func (ps *postgresStorage) Setup() error {
	sql := strings.Replace(`
		CREATE TABLE IF NOT EXISTS %%TABLE_NAME%% (
			id character varying(128) PRIMARY KEY,
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
