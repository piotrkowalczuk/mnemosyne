package mnemosyne

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/sha3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	// TokenContextKey is used by Mnemosyne internally to retrieve session token from context.Context.
	TokenContextKey = "mnemosyne_token"
	// TokenMetadataKey is used by Mnemosyne to retrieve session token from gRPC metadata object.
	TokenMetadataKey = "mnemosyne_token"
)

var (
	// ErrSessionNotFound can be returned by any endpoint if session does not exists.
	ErrSessionNotFound = grpc.Errorf(codes.NotFound, "mnemosyne: session not found")
)

//// NewTokenContext returns a new Context that carries Token value.
//func NewTokenContext(ctx context.Context, t Token) context.Context {
//	return context.WithValue(ctx, TokenContextKey, t)
//}
//
//// TokenFromContext returns the Token value stored in context, if any.
//func TokenFromContext(ctx context.Context) (Token, bool) {
//	t, ok := ctx.Value(TokenContextKey).(Token)
//
//	return t, ok
//}

// Mnemosyne ...
type Mnemosyne interface {
	FromContext(context.Context) (*Session, error)
	Get(context.Context, Token) (*Session, error)
	Exists(context.Context, Token) (bool, error)
	Start(context.Context, string, map[string]string) (*Session, error)
	Abandon(context.Context, Token) error
	SetValue(context.Context, Token, string, string) (map[string]string, error)
	//	DeleteValue(context.Context, string) (*Session, error)
	//	Clear(context.Context) error
}

type mnemosyne struct {
	metadata []string
	client   RPCClient
}

// MnemosyneOpts ...
type MnemosyneOpts struct {
	Metadata []string
}

// New allocates new mnemosyne instance.
func New(conn *grpc.ClientConn, options MnemosyneOpts) Mnemosyne {
	return &mnemosyne{
		client: NewRPCClient(conn),
	}
}

// FromContext implements Mnemosyne interface.
func (m *mnemosyne) FromContext(ctx context.Context) (*Session, error) {
	return m.client.Context(ctx, nil)
}

// Get implements Mnemosyne interface.
func (m *mnemosyne) Get(ctx context.Context, token Token) (*Session, error) {
	res, err := m.client.Get(ctx, &GetRequest{Token: &token})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Exists implements Mnemosyne interface.
func (m *mnemosyne) Exists(ctx context.Context, token Token) (bool, error) {
	res, err := m.client.Exists(ctx, &ExistsRequest{Token: &token})

	if err != nil {
		return false, err
	}

	return res.Exists, nil
}

// Create implements Mnemosyne interface.
func (m *mnemosyne) Start(ctx context.Context, subjectID string, data map[string]string) (*Session, error) {
	res, err := m.client.Start(ctx, &StartRequest{
		SubjectId: subjectID,
		Bag:       data,
	})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Abandon implements Mnemosyne interface.
func (m *mnemosyne) Abandon(ctx context.Context, token Token) error {
	_, err := m.client.Abandon(ctx, &AbandonRequest{Token: &token})

	return err
}

// SetData implements Mnemosyne interface.
func (m *mnemosyne) SetValue(ctx context.Context, token Token, key, value string) (map[string]string, error) {
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
	return []interface{}{"token", gr.Token.Encode()}
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
	return []interface{}{"token", er.Token.Encode()}
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
		"token", ar.Token.Encode(),
	}
}

// Context implements sklog.Contexter interface.
func (svr *SetValueRequest) Context() []interface{} {
	return []interface{}{
		"token", svr.Token.Encode(),
		"bag_key", svr.Key,
		"bag_value", svr.Value,
	}
}

// Context implements sklog.Contexter interface.
func (dvr *DeleteValueRequest) Context() []interface{} {
	return []interface{}{
		"token", dvr.Token.Encode(),
		"bag_key", dvr.Key,
	}
}

// Context implements sklog.Contexter interface.
func (cr *ClearRequest) Context() []interface{} {
	return []interface{}{
		"token", cr.Token.Encode(),
	}
}

// Context implements sklog.Contexter interface.
func (dr *DeleteRequest) Context() []interface{} {
	return []interface{}{
		"token", dr.Token.Encode(),
		"expire_at_from", dr.ExpireAtFrom,
		"expire_at_to", dr.ExpireAtTo,
	}
}

// ParseTime ...
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// EncodeToken ...
func EncodeToken(key, hash []byte) Token {
	t := Token{
		Key:  make([]byte, hex.EncodedLen(len(key))),
		Hash: make([]byte, hex.EncodedLen(len(hash))),
	}
	hex.Encode(t.Key, key)
	hex.Encode(t.Hash, hash)

	return t
}

// EncodeTokenString ...
func EncodeTokenString(key, hash string) (t Token) {
	return EncodeToken([]byte(key), []byte(hash))
}

// Value implements driver.Valuer interface.
func (t Token) Value() (driver.Value, error) {
	return t.Encode(), nil
}

// Scan implements sql.Scanner interface.
func (t *Token) Scan(src interface{}) error {
	var (
		token Token
		err   error
	)

	switch s := src.(type) {
	case []byte:
		token = DecodeToken(s)
	case string:
		if token, err = DecodeTokenString(s); err != nil {
			return err
		}

	default:
		return errors.New("mnemosyne: token supports scan only from slice of bytes and string")
	}

	*t = token

	return nil
}

// Encode ...
func (t *Token) Encode() []byte {
	return append(append(t.Key, ':'), t.Hash...)
}

// IsEmpty
func (t *Token) IsEmpty() bool {
	if t == nil {
		return true
	}
	return len(t.Hash) == 0
}

// DecodeToken parse string and allocates new token instance if ok. Expected token has format <key>:<hash>.
// If given string does not satisfy such pattern,
// entire string (excluding extremely situated colons) will be threaten like a hash.
func DecodeToken(s []byte) (t Token) {
	parts := bytes.Split(s, []byte(":"))

	if len(parts) == 1 {
		return Token{
			Hash: bytes.TrimSpace(parts[0]),
		}
	}

	if len(parts[1]) == 0 {
		return Token{
			Hash: bytes.TrimSpace(parts[0]),
		}
	}

	return Token{
		Key:  bytes.TrimSpace(parts[0]),
		Hash: bytes.TrimSpace(parts[1]),
	}
}

// DecodeTokenString ...
func DecodeTokenString(s string) (Token, error) {
	hx, err := hex.DecodeString(s)
	if err != nil {
		return Token{}, err
	}

	return DecodeToken(hx), nil
}

// RandomToken ...
func RandomToken(g RandomBytesGenerator, k []byte) (t Token, err error) {
	var buf []byte
	buf, err = g.GenerateRandomBytes(128)
	if err != nil {
		return
	}

	// A hash needs to be 64 bytes long to have 256-bit collision resistance.
	hash := make([]byte, 64)
	// Compute a 64-byte hash of buf and put it in h.
	sha3.ShakeSum256(hash, buf)

	t = EncodeToken(k, hash)
	return
}

//// TokenContextMiddleware puts token taken from header into current context.
//func TokenContextMiddleware(header string) func(fn func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request) {
//	return func(fn func(context.Context, http.ResponseWriter, *http.Request)) func(context.Context, http.ResponseWriter, *http.Request) {
//		return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
//			token := r.Header.Get(header)
//			ctx = NewTokenContext(ctx, DecodeToken(token))
//
//			rw.Header().Set(header, token)
//			fn(ctx, rw, r)
//		}
//	}
//}
