package mnemosyne

import (
	"crypto/x509"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Mnemosyne ...
type Mnemosyne interface {
	Close() error
	FromContext(ctx context.Context) (ses *mnemosynerpc.Session, err error)
	Get(ctx context.Context, token string) (ses *mnemosynerpc.Session, err error)
	Start(ctx context.Context, subjectID string, subjectClient string, bag map[string]string) (ses *mnemosynerpc.Session, err error)
	Exists(ctx context.Context, token string) (bool, error)
	Abandon(ctx context.Context, token string) error
	SetValue(ctx context.Context, token, key, value string) (bag map[string]string, err error)
}

type mnemosyne struct {
	metadata []string
	conn     *grpc.ClientConn
	client   mnemosynerpc.RPCClient
}

// MnemosyneOpts ...
type MnemosyneOpts struct {
	Metadata []string
	// Only one supported now.
	Addresses   []string
	UserAgent   string
	Certificate *x509.CertPool
	Block       bool
	Timeout     time.Duration
}

// New allocates new mnemosyne instance.
func New(opts MnemosyneOpts) (Mnemosyne, error) {
	if len(opts.Addresses) == 0 {
		return nil, errors.New("mnemosyne: at least one address needs to be provided")
	}
	if len(opts.Addresses) > 1 {
		return nil, errors.New("mnemosyne: client side load balancing is not implemented yet, only one address can be provided")
	}
	var dialOpts []grpc.DialOption
	if opts.Certificate == nil {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(opts.Certificate, "")))
	}
	if opts.Block {
		dialOpts = append(dialOpts, grpc.WithBlock())
	}
	if opts.Timeout.Nanoseconds() > 0 {
		dialOpts = append(dialOpts, grpc.WithTimeout(opts.Timeout))
	}
	if opts.UserAgent != "" {
		dialOpts = append(dialOpts, grpc.WithUserAgent(opts.UserAgent))
	}
	conn, err := grpc.Dial(opts.Addresses[0], dialOpts...)
	if err != nil {
		return nil, err
	}
	return &mnemosyne{
		metadata: opts.Metadata,
		conn:     conn,
		client:   mnemosynerpc.NewRPCClient(conn),
	}, nil
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
func (m *mnemosyne) Get(ctx context.Context, token string) (*mnemosynerpc.Session, error) {
	at := mnemosynerpc.ParseAccessToken(token)
	res, err := m.client.Get(ctx, &mnemosynerpc.GetRequest{AccessToken: &at})
	if err != nil {
		return nil, err
	}

	return res.Session, nil
}

// Exists implements Mnemosyne interface.
func (m *mnemosyne) Exists(ctx context.Context, token string) (bool, error) {
	at := mnemosynerpc.ParseAccessToken(token)
	res, err := m.client.Exists(ctx, &mnemosynerpc.ExistsRequest{AccessToken: &at})

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
func (m *mnemosyne) Abandon(ctx context.Context, token string) error {
	at := mnemosynerpc.ParseAccessToken(token)
	_, err := m.client.Abandon(ctx, &mnemosynerpc.AbandonRequest{AccessToken: &at})

	return err
}

// SetData implements Mnemosyne interface.
func (m *mnemosyne) SetValue(ctx context.Context, token, key, value string) (map[string]string, error) {
	at := mnemosynerpc.ParseAccessToken(token)
	res, err := m.client.SetValue(ctx, &mnemosynerpc.SetValueRequest{
		AccessToken: &at,
		Key:         key,
		Value:       value,
	})

	if err != nil {
		return nil, err
	}

	return res.Bag, nil
}

// Close implements io Closer interface.
func (m *mnemosyne) Close() error {
	return m.conn.Close()
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
