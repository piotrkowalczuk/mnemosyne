package mnemosyne

import (
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

const (
	// TokenContextKey ...
	TokenContextKey = "mnemosyne_token"
	// AccessTokenContextKey is used by Mnemosyne internally to retrieve session token from context.Context.
	AccessTokenContextKey = "mnemosyne_access_token"
	// AccessTokenMetadataKey is used by Mnemosyne to retrieve session token from gRPC metadata object.
	AccessTokenMetadataKey = "authorization"
)

// Token implements oauth2.TokenSource interface.
func (s *Session) Token() (*oauth2.Token, error) {
	var (
		err      error
		expireAt time.Time
	)
	if s.ExpireAt != nil {
		expireAt, err = ptypes.Timestamp(s.ExpireAt)
		if err != nil {
			return nil, err
		}
	}
	token := &oauth2.Token{
		AccessToken: s.AccessToken.Encode(),
		Expiry:      expireAt,
	}
	if s.Bag != nil && len(s.Bag) > 0 {
		token = token.WithExtra(bagToURLValues(s.Bag))
	}

	return token, nil
}

func bagToURLValues(b map[string]string) url.Values {
	r := make(map[string][]string, len(b))
	for k, v := range b {
		r[k] = []string{v}
	}
	return url.Values(r)
}

// Mnemosyne ...
type Mnemosyne interface {
	FromContext(context.Context) (*Session, error)
	Get(context.Context, AccessToken) (*Session, error)
	Exists(context.Context, AccessToken) (bool, error)
	Start(context.Context, string, string, map[string]string) (*Session, error)
	Abandon(context.Context, AccessToken) error
	SetValue(context.Context, AccessToken, string, string) (map[string]string, error)
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
	return m.client.Context(ctx, &empty.Empty{})
}

// Get implements Mnemosyne interface.
func (m *mnemosyne) Get(ctx context.Context, token AccessToken) (*Session, error) {
	res, err := m.client.Get(ctx, &GetRequest{AccessToken: &token})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Exists implements Mnemosyne interface.
func (m *mnemosyne) Exists(ctx context.Context, token AccessToken) (bool, error) {
	res, err := m.client.Exists(ctx, &ExistsRequest{AccessToken: &token})

	if err != nil {
		return false, err
	}

	return res.Exists, nil
}

// Create implements Mnemosyne interface.
func (m *mnemosyne) Start(ctx context.Context, subjectID, subjectClient string, data map[string]string) (*Session, error) {
	res, err := m.client.Start(ctx, &StartRequest{
		SubjectId:     subjectID,
		SubjectClient: subjectClient,
		Bag:           data,
	})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Abandon implements Mnemosyne interface.
func (m *mnemosyne) Abandon(ctx context.Context, token AccessToken) error {
	_, err := m.client.Abandon(ctx, &AbandonRequest{AccessToken: &token})

	return err
}

// SetData implements Mnemosyne interface.
func (m *mnemosyne) SetValue(ctx context.Context, token AccessToken, key, value string) (map[string]string, error) {
	res, err := m.client.SetValue(ctx, &SetValueRequest{
		AccessToken: &token,
		Key:         key,
		Value:       value,
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
