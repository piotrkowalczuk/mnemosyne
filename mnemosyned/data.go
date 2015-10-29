package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Data ...
type Data map[string]string

// DecodeSessionDataFromJSON ...
func DecodeSessionDataFromJSON(b []byte) (Data, error) {
	var data Data
	err := json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// EncodeToJSON ...
func (d *Data) EncodeToJSON() ([]byte, error) {
	if d == nil {
		return []byte{}, nil
	}

	return json.Marshal(d)
}

// Scan satisfy sql.Scanner interface.
func (d *Data) Scan(src interface{}) error {
	var err error
	var data Data

	switch t := src.(type) {
	case []byte:
		data, err = DecodeSessionDataFromJSON(t)
	default:
		return errors.New("mnemosyne: unsuported data source type")
	}

	if err != nil {
		return err
	}

	*d = data

	return nil
}

// Value satisfy driver.Valuer interface.
func (d Data) Value() (driver.Value, error) {
	return d.EncodeToJSON()
}

// Set ...
func (d *Data) Set(key, value string) {
	(*d)[key] = value
}

// Get ...
func (d *Data) Get(key string) string {
	return (*d)[key]
}
