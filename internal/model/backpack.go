package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"errors"
)

// Bag is a simple abstraction on the top of a map.
// It can be stored in a SQL database as a gob BLOB.
type Bag map[string]string

// Scan satisfy sql.Scanner interface.
func (b *Bag) Scan(src interface{}) (err error) {
	switch t := src.(type) {
	case []byte:
		err = gob.NewDecoder(bytes.NewReader(t)).Decode(b)
	default:
		return errors.New("unsupported data source type")
	}

	return
}

// Value satisfy driver.Valuer interface.
func (b Bag) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(b)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Set assigns the value to a given key.
func (b *Bag) Set(key, value string) {
	(*b)[key] = value
}

// Get returns a value for the given key.
func (b *Bag) Get(key string) string {
	return (*b)[key]
}

// Has checks if the given key is set.
func (b *Bag) Has(key string) bool {
	_, ok := (*b)[key]

	return ok
}
