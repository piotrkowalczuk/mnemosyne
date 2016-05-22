package mnemosyne

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Mnemosyne ...
type Mnemosyne interface {
	FromContext(context.Context) (*mnemosynerpc.Session, error)
	Get(context.Context, mnemosynerpc.AccessToken) (*mnemosynerpc.Session, error)
	Exists(context.Context, mnemosynerpc.AccessToken) (bool, error)
	Start(context.Context, string, string, map[string]string) (*mnemosynerpc.Session, error)
	Abandon(context.Context, mnemosynerpc.AccessToken) error
	SetValue(context.Context, mnemosynerpc.AccessToken, string, string) (map[string]string, error)
}

type mnemosyne struct {
	metadata []string
	client   mnemosynerpc.RPCClient
}

// MnemosyneOpts ...
type MnemosyneOpts struct {
	Metadata []string
}

// New allocates new mnemosyne instance.
func New(conn *grpc.ClientConn, options MnemosyneOpts) Mnemosyne {
	return &mnemosyne{
		client: mnemosynerpc.NewRPCClient(conn),
	}
}

// FromContext implements Mnemosyne interface.
func (m *mnemosyne) FromContext(ctx context.Context) (*mnemosynerpc.Session, error) {
	res, err := m.client.Context(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return res.Session, nil
}

// Get implements Mnemosyne interface.
func (m *mnemosyne) Get(ctx context.Context, token mnemosynerpc.AccessToken) (*mnemosynerpc.Session, error) {
	res, err := m.client.Get(ctx, &mnemosynerpc.GetRequest{AccessToken: &token})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Exists implements Mnemosyne interface.
func (m *mnemosyne) Exists(ctx context.Context, token mnemosynerpc.AccessToken) (bool, error) {
	res, err := m.client.Exists(ctx, &mnemosynerpc.ExistsRequest{AccessToken: &token})

	if err != nil {
		return false, err
	}

	return res.Exists, nil
}

// Create implements Mnemosyne interface.
func (m *mnemosyne) Start(ctx context.Context, subjectID, subjectClient string, data map[string]string) (*mnemosynerpc.Session, error) {
	res, err := m.client.Start(ctx, &mnemosynerpc.StartRequest{
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
func (m *mnemosyne) Abandon(ctx context.Context, token mnemosynerpc.AccessToken) error {
	_, err := m.client.Abandon(ctx, &mnemosynerpc.AbandonRequest{AccessToken: &token})

	return err
}

// SetData implements Mnemosyne interface.
func (m *mnemosyne) SetValue(ctx context.Context, token mnemosynerpc.AccessToken, key, value string) (map[string]string, error) {
	res, err := m.client.SetValue(ctx, &mnemosynerpc.SetValueRequest{
		AccessToken: &token,
		Key:         key,
		Value:       value,
	})

	if err != nil {
		return nil, err
	}

	return res.Bag, nil
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
