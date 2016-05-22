package mnemosyned

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	DefaultTTL = 24 * time.Minute
	DefaultTTC = 1 * time.Minute
)

var (
	// ErrSessionNotFound can be returned by any endpoint if session does not exists.
	ErrSessionNotFound = grpc.Errorf(codes.NotFound, "mnemosyned: session not found")
	// mnemosynerpc.ErrMissingAccessToken can be returned by any endpoint that expects access token in request.
	ErrMissingAccessToken = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing access token")
	// ErrMissingSubjectID can be returned by start endpoint if subject was not provided.
	ErrMissingSubjectID = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing subject id")
)

type rpcServer struct {
	ttc     time.Duration
	logger  log.Logger
	monitor *monitoring
	storage Storage
	alloc   struct {
		abandon  handlerFunc
		context  handlerFunc
		delete   handlerFunc
		exists   handlerFunc
		get      handlerFunc
		list     handlerFunc
		setValue handlerFunc
		start    handlerFunc
	}
}

func newRPCServer(l log.Logger, s Storage, m *monitoring, ttc time.Duration) *rpcServer {
	return &rpcServer{
		ttc: ttc,
		alloc: struct {
			abandon  handlerFunc
			context  handlerFunc
			delete   handlerFunc
			exists   handlerFunc
			get      handlerFunc
			list     handlerFunc
			setValue handlerFunc
			start    handlerFunc
		}{
			abandon:  newHandlerFunc("abandon"),
			context:  newHandlerFunc("context"),
			delete:   newHandlerFunc("delete"),
			exists:   newHandlerFunc("exists"),
			get:      newHandlerFunc("get"),
			list:     newHandlerFunc("list"),
			setValue: newHandlerFunc("set_value"),
			start:    newHandlerFunc("start"),
		},
		logger:  l,
		storage: s,
		monitor: m,
	}
}

// Get implements RPCServer interface.
func (rs *rpcServer) Context(ctx context.Context, req *empty.Empty) (*mnemosynerpc.ContextResponse, error) {
	h := rs.alloc.context(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	ses, err := h.context(ctx)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session has been retrieved (by context)")

	return &mnemosynerpc.ContextResponse{
		Session: ses,
	}, nil
}

// Get implements RPCServer interface.
func (rs *rpcServer) Get(ctx context.Context, req *mnemosynerpc.GetRequest) (*mnemosynerpc.GetResponse, error) {
	h := rs.alloc.get(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	ses, err := h.get(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session has been retrieved (by token)")

	return &mnemosynerpc.GetResponse{
		Session: ses,
	}, nil
}

// List implements RPCServer interface.
func (rs *rpcServer) List(ctx context.Context, req *mnemosynerpc.ListRequest) (*mnemosynerpc.ListResponse, error) {
	h := rs.alloc.list(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	sessions, err := h.list(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session list has been retrieved")

	return &mnemosynerpc.ListResponse{
		Sessions: sessions,
	}, nil
}

// Start implements RPCServer interface.
func (rs *rpcServer) Start(ctx context.Context, req *mnemosynerpc.StartRequest) (*mnemosynerpc.StartResponse, error) {
	h := rs.alloc.start(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	ses, err := h.start(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session has been started")

	return &mnemosynerpc.StartResponse{
		Session: ses,
	}, nil
}

// Exists implements RPCServer interface.
func (rs *rpcServer) Exists(ctx context.Context, req *mnemosynerpc.ExistsRequest) (*mnemosynerpc.ExistsResponse, error) {
	h := rs.alloc.exists(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	exists, err := h.exists(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session presence has been checked")

	return &mnemosynerpc.ExistsResponse{
		Exists: exists,
	}, nil
}

// Abandon implements RPCServer interface.
func (rs *rpcServer) Abandon(ctx context.Context, req *mnemosynerpc.AbandonRequest) (*mnemosynerpc.AbandonResponse, error) {
	h := rs.alloc.abandon(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	abandoned, err := h.abandon(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session has been abandoned")

	return &mnemosynerpc.AbandonResponse{
		Abandoned: abandoned,
	}, nil
}

// SetValue implements RPCServer interface.
func (rs *rpcServer) SetValue(ctx context.Context, req *mnemosynerpc.SetValueRequest) (*mnemosynerpc.SetValueResponse, error) {
	h := rs.alloc.setValue(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	bag, err := h.setValue(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session bag value has been set")

	return &mnemosynerpc.SetValueResponse{
		Bag: bag,
	}, nil
}

// Delete implements RPCServer interface.
func (rs *rpcServer) Delete(ctx context.Context, req *mnemosynerpc.DeleteRequest) (*mnemosynerpc.DeleteResponse, error) {
	h := rs.alloc.delete(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	affected, err := h.delete(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session value has been deleted")

	return &mnemosynerpc.DeleteResponse{
		Count: affected,
	}, nil
}

func (rs *rpcServer) cleanup(done chan struct{}) {
	sklog.Info(rs.logger, "cleanup routing started")
InfLoop:
	for {
		select {
		case <-time.After(rs.ttc):
			t := time.Now()
			sklog.Debug(rs.logger, "session cleanup start", "start_at", t.Format(time.RFC3339))
			affected, err := rs.storage.Delete(nil, nil, &t)
			if err != nil {
				if rs.monitor.enabled {
					rs.monitor.general.errors.Add(1)
				}
				sklog.Error(rs.logger, fmt.Errorf("session cleanup failure: %s", err.Error()), "expire_at_to", t)
				return
			}

			sklog.Debug(rs.logger, "session cleanup success", "count", affected, "elapsed", time.Now().Sub(t))
		case <-done:
			sklog.Info(rs.logger, "cleanup routing terminated")
			break InfLoop
		}
	}
}

func (rs *rpcServer) loggerBackground(ctx context.Context, keyval ...interface{}) log.Logger {
	l := log.NewContext(rs.logger).With(keyval...)
	if md, ok := metadata.FromContext(ctx); ok {
		if rid, ok := md["request_id"]; ok && len(rid) >= 1 {
			l = l.With("request_id", rid[0])
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		l = l.With("peer_address", p.Addr.String())
	}

	return l
}
