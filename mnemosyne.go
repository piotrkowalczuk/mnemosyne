package mnemosyne

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/crypto/sha3"
	"golang.org/x/net/context"
)

const (
	contextKeyToken = "context_key_mnemosyne_token"
)

// NewTokenContext returns a new Context that carries Token value.
func NewTokenContext(ctx context.Context, t Token) context.Context {
	return context.WithValue(ctx, contextKeyToken, t)
}

// TokenFromContext returns the Token value stored in context, if any.
func TokenFromContext(ctx context.Context) (Token, bool) {
	t, ok := ctx.Value(contextKeyToken).(Token)
	return t, ok
}

// Mnemosyne ...
type Mnemosyne interface {
	Get(context.Context) (*Session, error)
	Exists(context.Context) (bool, error)
	Start(context.Context, map[string]string) (*Session, error)
	Abandon(context.Context) error
	SetValue(context.Context, string, string) (map[string]string, error)
	//	DeleteValue(context.Context, string) (*Session, error)
	//	Clear(context.Context) error
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
func (m *mnemosyne) Get(ctx context.Context) (*Session, error) {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return nil, errors.New("mnemosyne: session cannot be retrieved, missing session token in the context")
	}

	res, err := m.client.Get(ctx, &GetRequest{Token: &token})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Exists implements Mnemosyne interface.
func (m *mnemosyne) Exists(ctx context.Context) (bool, error) {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return false, errors.New("mnemosyne: session existance cannot be checked, missing session token in the context")
	}
	res, err := m.client.Exists(ctx, &ExistsRequest{Token: &token})

	if err != nil {
		return false, err
	}

	return res.Exists, nil
}

// Create implements Mnemosyne interface.
func (m *mnemosyne) Start(ctx context.Context, data map[string]string) (*Session, error) {
	res, err := m.client.Start(ctx, &StartRequest{Bag: data})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Abandon implements Mnemosyne interface.
func (m *mnemosyne) Abandon(ctx context.Context) error {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return errors.New("mnemosyne: session cannot be abandoned, missing session token in the context")
	}
	_, err := m.client.Abandon(ctx, &AbandonRequest{Token: &token})

	return err
}

// SetData implements Mnemosyne interface.
func (m *mnemosyne) SetValue(ctx context.Context, key, value string) (map[string]string, error) {
	token, ok := TokenFromContext(ctx)
	if !ok {
		return nil, errors.New("mnemosyne: session value cannot be set, missing session token in the context")
	}
	res, err := m.client.SetValue(ctx, &SetValueRequest{
		Token: &token,
		Key:   key,
		Value: value,
	})

	if err != nil {
		return nil, err
	}

	return res.Bag, nil
}

//// DeleteValue implements Mnemosyne interface.
//func (m *mnemosyne) DeleteValue(ctx context.Context, key string) (*Session, error) {
//	token, ok := TokenFromContext(ctx)
//	if !ok {
//		return nil, errors.New("mnemosyne: session value cannot be deleted, missing session token in the context")
//	}
//	res, err := m.client.DeleteValue(ctx, &DeleteValueRequest{
//		Token: &token,
//		Key:   key,
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	return res.Session, nil
//}

//// Clear ...
//func (m *mnemosyne) Clear(ctx context.Context) error {
//	token, ok := TokenFromContext(ctx)
//	if !ok {
//		return errors.New("mnemosyne: session bag cannot be cleared, missing session token in the context")
//	}
//	_, err := m.client.Clear(ctx, &ClearRequest{Token: &token})
//
//	return err
//}

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

// Context implements sklog.Contexter interface.
func (er *ExistsRequest) Context() []interface{} {
	return []interface{}{"token", er.Token}
}

// Context implements sklog.Contexter interface.
func (er *StartRequest) Context() (ctx []interface{}) {
	for key, value := range er.Bag {
		ctx = append(ctx, "bag_"+key, value)
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
func (svr *SetValueRequest) Context() []interface{} {
	return []interface{}{
		"token", svr.Token,
		"key", svr.Key,
		"value", svr.Value,
	}
}

// Context implements sklog.Contexter interface.
func (dvr *DeleteValueRequest) Context() []interface{} {
	return []interface{}{
		"token", dvr.Token,
		"key", dvr.Key,
	}
}

// Context implements sklog.Contexter interface.
func (cr *ClearRequest) Context() []interface{} {
	return []interface{}{
		"token", cr.Token,
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
	var token Token

	switch s := src.(type) {
	case []byte:
		token = NewToken(string(s))
	case string:
		token = NewToken(s)
	default:
		return errors.New("mnemosyne: token supports scan only from slice of bytes and string")
	}

	*t = token

	return nil
}

// NewToken parse string and allocates new token instance if ok. Expected token has format <key>:<hash>.
// If given string does not satisfy such pattern,
// entire string (excluding extremely situated colons) will be threaten like a hash.
func NewToken(s string) (t Token) {
	parts := strings.Split(s, ":")

	if len(parts) == 1 {
		t = Token{
			Hash: parts[0],
		}

		return
	}

	if parts[1] == "" {
		t = Token{
			Hash: parts[0],
		}

		return
	}

	t = Token{
		Key:  parts[0],
		Hash: parts[1],
	}

	return
}

// NewTokenRandom ...
func NewTokenRandom(g RandomBytesGenerator, k string) (t Token, err error) {
	var buf []byte
	buf, err = g.GenerateRandomBytes(128)
	if err != nil {
		return
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	hash := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(hash, buf)

	t = Token{
		Key:  k,
		Hash: hex.EncodeToString(hash),
	}
	return
}

// TokenContextMiddleware puts token taken from header into current context.
func TokenContextMiddleware(header string) func(fn func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request) {
	return func(fn func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request) {
		return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
			ctx = NewTokenContext(ctx, NewToken(r.Header.Get(header)))

			fn(ctx, rw, r)
		}
	}
}
