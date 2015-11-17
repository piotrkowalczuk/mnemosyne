package mnemosyne

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/crypto/sha3"
	"golang.org/x/net/context"
)

const (
	contextKeySession = "context_key_mnemosyne_session"
)

// NewContext returns a new Context that carries Session value.
func NewContext(ctx context.Context, ses Session) context.Context {
	return context.WithValue(ctx, contextKeySession, ses)
}

// FromContext returns the Session value stored in context, if any.
func FromContext(ctx context.Context) (Session, bool) {
	c, ok := ctx.Value(contextKeySession).(Session)
	return c, ok
}

// Mnemosyne ...
type Mnemosyne interface {
	Get(context.Context, *Token) (*Session, error)
	Exists(context.Context, *Token) (bool, error)
	Create(context.Context, map[string]string) (*Session, error)
	Abandon(context.Context, *Token) (bool, error)
	SetData(context.Context, *Token, string, string) (*Session, error)
}

type mnemosyne struct {
	client RPCClient
}

// MnemosyneOpts ...
type MnemosyneOpts struct {
}

// New allocates new mnemosyne instance.
func New(conn *grpc.ClientConn, options MnemosyneOpts) Mnemosyne {
	return &mnemosyne{
		client: NewRPCClient(conn),
	}
}

// Get implements Mnemosyne interface.
func (m *mnemosyne) Get(ctx context.Context, token *Token) (*Session, error) {
	res, err := m.client.Get(ctx, &GetRequest{Token: token})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Exists implements Mnemosyne interface.
func (m *mnemosyne) Exists(ctx context.Context, token *Token) (bool, error) {
	res, err := m.client.Exists(ctx, &ExistsRequest{Token: token})

	if err != nil {
		return false, err
	}

	return res.Exists, nil
}

// Create implements Mnemosyne interface.
func (m *mnemosyne) Create(ctx context.Context, data map[string]string) (*Session, error) {
	res, err := m.client.Create(ctx, &CreateRequest{Data: data})

	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Abandon implements Mnemosyne interface.
func (m *mnemosyne) Abandon(ctx context.Context, token *Token) (bool, error) {
	res, err := m.client.Abandon(ctx, &AbandonRequest{Token: token})

	if err != nil {
		return false, err
	}

	return res.Abandoned, nil
}

// SetData implements Mnemosyne interface.
func (m *mnemosyne) SetData(ctx context.Context, token *Token, key, value string) (*Session, error) {
	res, err := m.client.SetData(ctx, &SetDataRequest{
		Token: token,
		Key:   key,
		Value: value,
	})

	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Context implements sklog.Contexter interface.
func (gr *GetRequest) Context() []interface{} {
	return []interface{}{"token", gr.Token}
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

// ExpireAtFromTime ...
func (lr *ListRequest) ExpireAtFromTime() time.Time {
	return TimestampToTime(lr.ExpireAtFrom)
}

// ExpireAtToTime ...
func (lr *ListRequest) ExpireAtToTime() time.Time {
	return TimestampToTime(lr.ExpireAtTo)
}

// Context implements sklog.Contexter interface.
func (er *ExistsRequest) Context() []interface{} {
	return []interface{}{"token", er.Token}
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
		"token", ar.Token,
	}
}

// Context implements sklog.Contexter interface.
func (sdr *SetDataRequest) Context() []interface{} {
	return []interface{}{
		"token", sdr.Token,
		"key", sdr.Key,
		"value", sdr.Value,
	}
}

// Context implements sklog.Contexter interface.
func (dr *DeleteRequest) Context() []interface{} {
	return []interface{}{
		"token", dr.Token,
		"expire_at_from", dr.ExpireAtFrom,
		"expire_at_to", dr.ExpireAtTo,
	}
}

// ExpireAtFromTime ...
func (dr *DeleteRequest) ExpireAtFromTime() time.Time {
	return TimestampToTime(dr.ExpireAtFrom)
}

// ExpireAtToTime ...
func (dr *DeleteRequest) ExpireAtToTime() time.Time {
	return TimestampToTime(dr.ExpireAtTo)
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

// ExpireAtTime ...
func (s *Session) ExpireAtTime() time.Time {
	return TimestampToTime(s.ExpireAt)
}

// ParseTime ...
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// Value implements driver.Valuer interface.
func (t Token) Value() (driver.Value, error) {
	return t.Key + ":" + t.Hash, nil
}

// Scan implements sql.Scanner interface.
func (t *Token) Scan(src interface{}) error {
	var token *Token
	var err error

	switch s := src.(type) {
	case []byte:
		token, err = NewTokenFromBytes(s)
	case string:
		token, err = NewTokenFromString(s)
	default:
		return errors.New("mnemosyne: token supports scan only from slice of bytes and string")
	}
	if err != nil {
		return err
	}

	*t = *token

	return nil
}

// NewTokenFromString parse string and allocates new token instance if ok.
func NewTokenFromString(s string) (*Token, error) {
	parts := strings.Split(s, ":")

	if len(parts) != 2 {
		return nil, errors.New("mnemosyne: token cannot be allocated, given string has wrong format")
	}

	return NewToken(parts[0], parts[1]), nil
}

// NewToken allocates new Token instance.
func NewToken(key, hash string) *Token {
	return &Token{
		Key:  key,
		Hash: hash,
	}
}

// NewTokenFromBytes ...
func NewTokenFromBytes(b []byte) (*Token, error) {
	parts := bytes.Split(b, []byte{':'})

	if len(parts) != 2 {
		return nil, errors.New("mnemosyne: token cannot be allocated, given byte slice has wrong format")
	}

	return NewToken(string(parts[0]), string(parts[1])), nil
}

// NewTokenRandom ...
func NewTokenRandom(g RandomBytesGenerator, k string) (*Token, error) {
	buf, err := g.GenerateRandomBytes(128)
	if err != nil {
		return nil, err
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	hash := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(hash, buf)

	return NewToken(k, hex.EncodeToString(hash)), nil
}
