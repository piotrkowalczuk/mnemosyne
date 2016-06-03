package mnemosyned

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"errors"
)

type bag map[string]string

// Scan satisfy sql.Scanner interface.
func (b *bag) Scan(src interface{}) (err error) {
	switch t := src.(type) {
	case []byte:
		err = gob.NewDecoder(bytes.NewReader(t)).Decode(b)
	default:
		return errors.New("unsuported data source type")
	}

	return
}

// Value satisfy driver.Valuer interface.
func (b bag) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(b)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b *bag) set(key, value string) {
	(*b)[key] = value
}

func (b *bag) get(key string) string {
	return (*b)[key]
}

func (b *bag) has(key string) bool {
	_, ok := (*b)[key]

	return ok
}
