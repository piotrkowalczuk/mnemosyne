package mnemosyned

import (
	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type rpcServer struct {
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

func newRPCServer(l log.Logger, s Storage, m *monitoring) *rpcServer {
	return &rpcServer{
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
func (rs *rpcServer) Context(ctx context.Context, req *empty.Empty) (*mnemosyne.Session, error) {
	h := rs.alloc.context(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	ses, err := h.context(ctx)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session has been retrieved (by context)")

	return ses, nil
}

// Get implements RPCServer interface.
func (rs *rpcServer) Get(ctx context.Context, req *mnemosyne.GetRequest) (*mnemosyne.GetResponse, error) {
	h := rs.alloc.get(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	ses, err := h.get(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session has been retrieved (by token)")

	return &mnemosyne.GetResponse{
		Session: ses,
	}, nil
}

// List implements RPCServer interface.
func (rs *rpcServer) List(ctx context.Context, req *mnemosyne.ListRequest) (*mnemosyne.ListResponse, error) {
	h := rs.alloc.list(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	sessions, err := h.list(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session list has been retrieved")

	return &mnemosyne.ListResponse{
		Sessions: sessions,
	}, nil
}

// Start implements RPCServer interface.
func (rs *rpcServer) Start(ctx context.Context, req *mnemosyne.StartRequest) (*mnemosyne.StartResponse, error) {
	h := rs.alloc.start(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	ses, err := h.start(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session has been started")

	return &mnemosyne.StartResponse{
		Session: ses,
	}, nil
}

// Exists implements RPCServer interface.
func (rs *rpcServer) Exists(ctx context.Context, req *mnemosyne.ExistsRequest) (*mnemosyne.ExistsResponse, error) {
	h := rs.alloc.exists(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	exists, err := h.exists(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session presence has been checked")

	return &mnemosyne.ExistsResponse{
		Exists: exists,
	}, nil
}

// Abandon implements RPCServer interface.
func (rs *rpcServer) Abandon(ctx context.Context, req *mnemosyne.AbandonRequest) (*mnemosyne.AbandonResponse, error) {
	h := rs.alloc.abandon(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	abandoned, err := h.abandon(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session has been abandoned")

	return &mnemosyne.AbandonResponse{
		Abandoned: abandoned,
	}, nil
}

// SetValue implements RPCServer interface.
func (rs *rpcServer) SetValue(ctx context.Context, req *mnemosyne.SetValueRequest) (*mnemosyne.SetValueResponse, error) {
	h := rs.alloc.setValue(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	bag, err := h.setValue(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session bag value has been set")

	return &mnemosyne.SetValueResponse{
		Bag: bag,
	}, nil
}

// Delete implements RPCServer interface.
func (rs *rpcServer) Delete(ctx context.Context, req *mnemosyne.DeleteRequest) (*mnemosyne.DeleteResponse, error) {
	h := rs.alloc.delete(rs.loggerBackground(ctx), rs.storage, rs.monitor.rpc)
	if h.monitor.enabled {
		h.monitor.requests.Add(1)
	}

	affected, err := h.delete(ctx, req)
	if err != nil {
		return nil, h.error(err)
	}

	sklog.Debug(h.logger, "session value has been deleted")

	return &mnemosyne.DeleteResponse{
		Count: affected,
	}, nil
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
