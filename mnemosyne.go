package mnemosyne

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

// Context implements sklog.Contexter interface.
func (gr *GetRequest) Context() []interface{} {
	return []interface{}{"id", gr.Id}
}

// Context implements sklog.Contexter interface.
func (lr *ListRequest) Context() []interface{} {
	return []interface{}{
		"offset", lr.Offset,
		"limit", lr.Limit,
		"expire_at_from", lr.ExpireAtFrom,
		"expire_at_to", lr.ExpireAtTo,
	}
}

// Context implements sklog.Contexter interface.
func (er *ExistsRequest) Context() []interface{} {
	return []interface{}{"id", er.Id}
}

// Context implements sklog.Contexter interface.
func (er *CreateRequest) Context() (ctx []interface{}) {
	for key, value := range er.Data {
		ctx = append(ctx, "data_"+key, value)
	}

	return
}

// Context implements sklog.Contexter interface.
func (ar *AbandonRequest) Context() []interface{} {
	return []interface{}{
		"id", ar.Id,
	}
}

// Context implements sklog.Contexter interface.
func (sdr *SetDataRequest) Context() []interface{} {
	return []interface{}{
		"id", sdr.Id,
		"key", sdr.Key,
		"value", sdr.Value,
	}
}

// Context implements sklog.Contexter interface.
func (dr *DeleteRequest) Context() []interface{} {
	return []interface{}{
		"id", dr.Id,
		"expire_at_from", dr.ExpireAtFrom,
		"expire_at_to", dr.ExpireAtTo,
	}
}

// SetValue ...
func (s *Session) SetValue(key, value string) {
	if s.Data == nil {
		s.Data = make(map[string]string)
	}

	s.Data[key] = value
}

// Value ...
func (s *Session) Value(key string) string {
	if s.Data == nil {
		s.Data = make(map[string]string)
	}

	return s.Data[key]
}

// ParseTime ...
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// Value implements driver.Valuer interface.
func (i ID) Value() (driver.Value, error) {
	return i.Key + ":" + i.Hash, nil
}

// Scan implements sql.Scanner interface.
func (i *ID) Scan(src interface{}) error {
	var id *ID
	var err error

	switch s := src.(type) {
	case []byte:
		id, err = NewIDFromBytes(s)
	case string:
		id, err = NewIDFromString(s)
	default:
		return errors.New("mnemosyne: id supports scan only from slice of bytes and string")
	}
	if err != nil {
		return err
	}

	*i = *id

	return nil
}

// NewIDFromString ...
func NewIDFromString(s string) (*ID, error) {
	parts := strings.Split(s, ":")

	if len(parts) != 2 {
		return nil, errors.New("mnemosyne: id cannot be allocated, given string has wrong format")
	}

	return &ID{
		Key:  parts[0],
		Hash: parts[1],
	}, nil
}

// NewIDFromBytes ...
func NewIDFromBytes(b []byte) (*ID, error) {
	parts := bytes.Split(b, []byte{':'})

	if len(parts) != 2 {
		return nil, errors.New("mnemosyne: id cannot be allocated, given byte slice has wrong format")
	}

	return &ID{
		Key:  string(parts[0]),
		Hash: string(parts[1]),
	}, nil
}

// NewIDRandom ...
func NewIDRandom(g RandomBytesGenerator, k string) (*ID, error) {
	buf, err := g.GenerateRandomBytes(128)
	if err != nil {
		return nil, err
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	id := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(id, buf)

	return &ID{
		Key:  k,
		Hash: hex.EncodeToString(id),
	}, nil
}
