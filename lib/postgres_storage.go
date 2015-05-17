package lib

import (
	"database/sql"
	"time"

	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

type PostgresStorage struct {
	db        *sql.DB
	tableName string
	generator RandomBytesGenerator
}

// NewPostgresStorage ...
func NewPostgresStorage(db *sql.DB, tableName string) *PostgresStorage {
	return &PostgresStorage{
		db:        db,
		tableName: tableName,
		generator: &SystemRandomBytesGenerator{},
	}
}

// Session ...
func (ps *PostgresStorage) Get(id SessionID) (*Session, error) {
	query := `
		SELECT data, expire_at
		FROM ` + ps.tableName + `
		WHERE id = $1
		LIMIT 1
	`

	data := []byte{}
	session := &Session{}
	err := ps.db.QueryRow(query, id.String()).Scan(
		&data,
		&session.ExpireAt,
	)
	if err != nil {
		return nil, err
	}

	session.ID = id
	session.Data, err = DecodeSessionDataFromJSON(data)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Session ...
func (ps *PostgresStorage) Exists(id SessionID) (bool, error) {
	var exists bool

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM ` + ps.tableName + `
			WHERE id = $1
		)
	`

	err := ps.db.QueryRow(query, id.String()).Scan(
		&exists,
	)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Session ...
func (ps *PostgresStorage) New(data SessionData) (*Session, error) {
	buf, err := ps.generator.GenerateRandomBytes(128)
	if err != nil {
		return nil, err
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	id := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(id, buf)

	expire_at := time.Now().Add(30 * time.Minute)
	session := &Session{
		ID:       SessionID(hex.EncodeToString(id)),
		Data:     data,
		ExpireAt: &expire_at,
	}

	if err := ps.save(session); err != nil {
		return nil, err
	}

	return session, nil
}

// Session ...
func (ps *PostgresStorage) save(session *Session) error {
	query := `
		INSERT INTO ` + ps.tableName + ` (id, data, expire_at)
		VALUES ($1, $2, $3)
	`

	encodedData, err := session.Data.EncodeToJSON()
	if err != nil {
		return err
	}

	_, err = ps.db.Exec(
		query,
		session.ID.String(),
		encodedData,
		session.ExpireAt,
	)

	return err
}

// Abandon ...
func (ps *PostgresStorage) Abandon(id SessionID) error {
	query := `
		DELETE FROM ` + ps.tableName + `
		WHERE id = $1
	`

	_, err := ps.db.Exec(query, id.String())
	if err != nil {
		return err
	}

	return nil
}

// SetData ...
func (ps *PostgresStorage) SetData(entry SessionDataEntry) (*Session, error) {
	session := &Session{
		ID: entry.ID,
	}
	data := []byte{}
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

	err = ps.db.QueryRow(selectQuery, session.ID.String()).Scan(
		&data,
		&session.ExpireAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	session.Data, err = DecodeSessionDataFromJSON(data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	session.Data.Set(entry.Key, entry.Value)

	encodedData, err := session.Data.EncodeToJSON()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = ps.db.Exec(
		updateQuery,
		session.ID.String(),
		encodedData,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return session, nil
}
