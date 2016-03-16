package mnemosyne

import (
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	// AccessTokenContextKey is used by Mnemosyne internally to retrieve session token from context.Context.
	AccessTokenContextKey = "mnemosyne_token"
	// AccessTokenMetadataKey is used by Mnemosyne to retrieve session token from gRPC metadata object.
	AccessTokenMetadataKey = "mnemosyne_token"
)

var (
	// ErrSessionNotFound can be returned by any endpoint if session does not exists.
	ErrSessionNotFound = grpc.Errorf(codes.NotFound, "mnemosyne: session not found")
	// ErrMissingToken can be returned by any endpoint that expects token in request.
	ErrMissingToken = grpc.Errorf(codes.InvalidArgument, "mnemosyne: missing token")
	// ErrMissingSubjectID can be returned by start endpoint if subject was not provided.
	ErrMissingSubjectID = grpc.Errorf(codes.InvalidArgument, "mnemosyne: missing subject id")
)

// Token implements oauth2.TokenSource interface.
func (s *Session) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: s.AccessToken.Encode(),
		Expiry:      s.ExpireAt.Time(),
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
	Start(context.Context, string, map[string]string) (*Session, error)
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
	return m.client.Context(ctx, &Empty{})
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

// Context implements sklog.Contexter interface.
func (gr *GetRequest) Context() []interface{} {
	return []interface{}{"token", gr.AccessToken.Bytes()}
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
	return []interface{}{"token", er.AccessToken.Bytes()}
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
		"token", ar.AccessToken.Bytes(),
	}
}

// Context implements sklog.Contexter interface.
func (svr *SetValueRequest) Context() []interface{} {
	return []interface{}{
		"token", svr.AccessToken.Bytes(),
		"bag_key", svr.Key,
		"bag_value", svr.Value,
	}
}

// Context implements sklog.Contexter interface.
func (dvr *DeleteValueRequest) Context() []interface{} {
	return []interface{}{
		"token", dvr.AccessToken.Bytes(),
		"bag_key", dvr.Key,
	}
}

// Context implements sklog.Contexter interface.
func (cr *ClearRequest) Context() []interface{} {
	return []interface{}{
		"token", cr.AccessToken.Bytes(),
	}
}

// Context implements sklog.Contexter interface.
func (dr *DeleteRequest) Context() []interface{} {
	return []interface{}{
		"token", dr.AccessToken.Bytes(),
		"expire_at_from", dr.ExpireAtFrom,
		"expire_at_to", dr.ExpireAtTo,
	}
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
