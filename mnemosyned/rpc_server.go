package mnemosyned

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	// DefaultTTL is session time to live default value.
	DefaultTTL = 24 * time.Minute
	// DefaultTTC is time to clear default value.
	DefaultTTC = 1 * time.Minute
)

var (
	errSessionNotFound    = grpc.Errorf(codes.NotFound, "mnemosyned: session not found")
	errMissingAccessToken = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing access token")
	errMissingSubjectID   = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing subject id")
)

type sessionManager struct {
	ttc     time.Duration
	logger  log.Logger
	monitor *monitoring
	storage storage
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

func newSessionManager(l log.Logger, s storage, m *monitoring, ttc time.Duration) *sessionManager {
	return &sessionManager{
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
func (sm *sessionManager) Context(ctx context.Context, req *empty.Empty) (*mnemosynerpc.ContextResponse, error) {
	h := sm.alloc.context(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	ses, err := h.context(ctx)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session has been retrieved (by context)")

	return &mnemosynerpc.ContextResponse{
		Session: ses,
	}, nil
}

// Get implements RPCServer interface.
func (sm *sessionManager) Get(ctx context.Context, req *mnemosynerpc.GetRequest) (*mnemosynerpc.GetResponse, error) {
	h := sm.alloc.get(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	ses, err := h.get(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session has been retrieved (by token)")

	return &mnemosynerpc.GetResponse{
		Session: ses,
	}, nil
}

// List implements RPCServer interface.
func (sm *sessionManager) List(ctx context.Context, req *mnemosynerpc.ListRequest) (*mnemosynerpc.ListResponse, error) {
	h := sm.alloc.list(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	sessions, err := h.list(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session list has been retrieved")

	return &mnemosynerpc.ListResponse{
		Sessions: sessions,
	}, nil
}

// Start implements RPCServer interface.
func (sm *sessionManager) Start(ctx context.Context, req *mnemosynerpc.StartRequest) (*mnemosynerpc.StartResponse, error) {
	h := sm.alloc.start(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	ses, err := h.start(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session has been started")

	return &mnemosynerpc.StartResponse{
		Session: ses,
	}, nil
}

// Exists implements RPCServer interface.
func (sm *sessionManager) Exists(ctx context.Context, req *mnemosynerpc.ExistsRequest) (*mnemosynerpc.ExistsResponse, error) {
	h := sm.alloc.exists(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	exists, err := h.exists(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session presence has been checked")

	return &mnemosynerpc.ExistsResponse{
		Exists: exists,
	}, nil
}

// Abandon implements RPCServer interface.
func (sm *sessionManager) Abandon(ctx context.Context, req *mnemosynerpc.AbandonRequest) (*mnemosynerpc.AbandonResponse, error) {
	h := sm.alloc.abandon(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	abandoned, err := h.abandon(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session has been abandoned")

	return &mnemosynerpc.AbandonResponse{
		Abandoned: abandoned,
	}, nil
}

// SetValue implements RPCServer interface.
func (sm *sessionManager) SetValue(ctx context.Context, req *mnemosynerpc.SetValueRequest) (*mnemosynerpc.SetValueResponse, error) {
	h := sm.alloc.setValue(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	bag, err := h.setValue(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session bag value has been set")

	return &mnemosynerpc.SetValueResponse{
		Bag: bag,
	}, nil
}

// Delete implements RPCServer interface.
func (sm *sessionManager) Delete(ctx context.Context, req *mnemosynerpc.DeleteRequest) (*mnemosynerpc.DeleteResponse, error) {
	h := sm.alloc.delete(sm.loggerBackground(ctx), sm.storage, sm.monitor.rpc)

	affected, err := h.delete(ctx, req)
	if err != nil {
		return nil, err
	}

	sklog.Debug(h.logger, "session value has been deleted")

	return &mnemosynerpc.DeleteResponse{
		Count: affected,
	}, nil
}

func (sm *sessionManager) cleanup(done chan struct{}) {
	logger := log.NewContext(sm.logger).WithPrefix("module", "cleanup")
	sklog.Info(sm.logger, "cleanup routing started")
InfLoop:
	for {
		select {
		case <-time.After(sm.ttc):
			t := time.Now()
			sklog.Debug(logger, "session cleanup start", "start_at", t.Format(time.RFC3339))
			affected, err := sm.storage.Delete("", nil, &t)
			if err != nil {
				if sm.monitor.enabled {
					sm.monitor.general.errors.With(prometheus.Labels{"action": "cleanup"}).Inc()
				}
				sklog.Error(logger, fmt.Errorf("session cleanup failure: %s", err.Error()), "expire_at_to", t)
				return
			}

			sklog.Debug(logger, "session cleanup success", "count", affected, "elapsed", time.Now().Sub(t))
		case <-done:
			sklog.Info(logger, "cleanup routing terminated")
			break InfLoop
		}
	}
}

func (sm *sessionManager) loggerBackground(ctx context.Context, keyval ...interface{}) log.Logger {
	l := log.NewContext(sm.logger).With(keyval...)
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
