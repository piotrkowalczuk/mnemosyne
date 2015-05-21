package lib

import (
	"database/sql"
	"strings"
	"time"

	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

const (
	PostgresEngine               = "postgres"
	PostgresTableNamePlaceholder = "%%TABLE_NAME%%"
	PostgresSchema               = `
		CREATE TABLE %%TABLE_NAME%% (
			id character varying(128) NOT NULL,
			data json NOT NULL,
			expire_at timestamp with time zone NOT NULL
		);

		ALTER TABLE ONLY %%TABLE_NAME%%
			ADD CONSTRAINT mnemosyne_session_pkey PRIMARY KEY (id);

		ALTER TABLE ONLY %%TABLE_NAME%%
			ADD CONSTRAINT mnemosyne_session_id_key UNIQUE (id);
    `
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

func (ps *PostgresStorage) Init() (err error) {
	sql := strings.Replace(PostgresSchema, PostgresTableNamePlaceholder, ps.tableName, -1)

	_, err = ps.db.Exec(sql)

	return
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

	session := &Session{
		ID:   SessionID(hex.EncodeToString(id)),
		Data: data,
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
		VALUES ($1, $2, NOW() + '30 minutes'::interval)
		RETURNING expire_at

	`

	encodedData, err := session.Data.EncodeToJSON()
	if err != nil {
		return err
	}

	err = ps.db.QueryRow(
		query,
		session.ID.String(),
		encodedData,
	).Scan(
		&session.ExpireAt,
	)

	return err
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
		if err == sql.ErrNoRows {
			return nil, ErrSessionNotFound
		}
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

// Abandon ...
func (ps *PostgresStorage) Abandon(id SessionID) error {
	query := `
		DELETE FROM ` + ps.tableName + `
		WHERE id = $1
	`

	result, err := ps.db.Exec(query, id.String())
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrSessionNotFound
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

	err = tx.QueryRow(selectQuery, session.ID.String()).Scan(
		&data,
		&session.ExpireAt,
	)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, ErrSessionNotFound
		}
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

	_, err = tx.Exec(
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

// Cleanup
func (ps *PostgresStorage) Cleanup(until *time.Time) (int64, error) {
	var err error
	var result sql.Result
	delUntil := `
		DELETE FROM ` + ps.tableName + `
		WHERE expire_at < $1
	`
	delAll := `TRUNCATE TABLE ` + ps.tableName

	if until == nil {
		result, err = ps.db.Exec(delAll)
	} else {
		result, err = ps.db.Exec(delUntil, until)
	}

	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
