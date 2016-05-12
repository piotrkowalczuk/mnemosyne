package mnemosyned

import (
	"bytes"
	"database/sql/driver"
	"encoding/gob"
	"errors"
)

type bagpack map[string]string

// Scan satisfy sql.Scanner interface.
func (b *bagpack) Scan(src interface{}) (err error) {
	switch t := src.(type) {
	case []byte:
		err = gob.NewDecoder(bytes.NewReader(t)).Decode(b)
	default:
		return errors.New("unsuported data source type")
	}

	return
}

// Value satisfy driver.Valuer interface.
func (bp bagpack) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(bp)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Set implements Bag interface.
func (bp *bagpack) Set(key, value string) {
	(*bp)[key] = value
}

// Get implements Bag interface.
func (bp *bagpack) Get(key string) string {
	return (*bp)[key]
}

// Has implements Bag interface.
func (bp *bagpack) Has(key string) bool {
	_, ok := (*bp)[key]

	return ok
}
