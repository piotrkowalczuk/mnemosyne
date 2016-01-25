package main

import (
	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

// Get implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Context(ctx context.Context, req *mnemosyne.Empty) (*mnemosyne.Session, error) {
	h := rs.alloc.context(rs.logger, rs.storage, rs.monitor.rpc)
	h.monitor.requests.Add(1)

	ses, err := h.context(ctx)
	if err != nil {
		h.monitor.errors.Add(1)
		sklog.Error(h.logger, err)

		return nil, rs.error(err)
	}

	sklog.Debug(h.logger, "session has been retrieved (by context)")

	return ses, nil
}

// Get implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Get(ctx context.Context, req *mnemosyne.GetRequest) (*mnemosyne.GetResponse, error) {
	h := rs.alloc.get(rs.logger, rs.storage, rs.monitor.rpc)
	h.monitor.requests.Add(1)

	ses, err := h.get(ctx, req)
	if err != nil {
		h.monitor.errors.Add(1)
		sklog.Error(h.logger, err)

		return nil, rs.error(err)
	}

	sklog.Debug(h.logger, "session has been retrieved (by token)")

	return &mnemosyne.GetResponse{
		Session: ses,
	}, nil
}

// List implements mnemosyne.RPCServer interface.
func (rs *rpcServer) List(ctx context.Context, req *mnemosyne.ListRequest) (*mnemosyne.ListResponse, error) {
	h := rs.alloc.list(rs.logger, rs.storage, rs.monitor.rpc)
	h.monitor.requests.Add(1)

	sessions, err := h.list(ctx, req)
	if err != nil {
		h.monitor.errors.Add(1)
		sklog.Error(h.logger, err)

		return nil, rs.error(err)
	}

	sklog.Debug(h.logger, "session list has been retrieved")

	return &mnemosyne.ListResponse{
		Sessions: sessions,
	}, nil
}

// Start implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Start(ctx context.Context, req *mnemosyne.StartRequest) (*mnemosyne.StartResponse, error) {
	h := rs.alloc.start(rs.logger, rs.storage, rs.monitor.rpc)
	h.monitor.requests.Add(1)

	ses, err := h.start(ctx, req)
	if err != nil {
		h.monitor.errors.Add(1)
		sklog.Error(h.logger, err)

		return nil, rs.error(err)
	}

	sklog.Debug(h.logger, "session has been started")

	return &mnemosyne.StartResponse{
		Session: ses,
	}, nil
}

// Exists implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Exists(ctx context.Context, req *mnemosyne.ExistsRequest) (*mnemosyne.ExistsResponse, error) {
	h := rs.alloc.exists(rs.logger, rs.storage, rs.monitor.rpc)
	h.monitor.requests.Add(1)

	exists, err := h.exists(ctx, req)
	if err != nil {
		h.monitor.errors.Add(1)
		sklog.Error(h.logger, err)

		return nil, rs.error(err)
	}

	sklog.Debug(h.logger, "session presence has been checked")

	return &mnemosyne.ExistsResponse{
		Exists: exists,
	}, nil
}

// Abandon implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Abandon(ctx context.Context, req *mnemosyne.AbandonRequest) (*mnemosyne.AbandonResponse, error) {
	h := rs.alloc.abandon(rs.logger, rs.storage, rs.monitor.rpc)
	h.monitor.requests.Add(1)

	abandoned, err := h.abandon(ctx, req)
	if err != nil {
		h.monitor.errors.Add(1)
		sklog.Error(h.logger, err)

		return nil, rs.error(err)
	}

	sklog.Debug(h.logger, "session has been abandoned")

	return &mnemosyne.AbandonResponse{
		Abandoned: abandoned,
	}, nil
}

// SetValue implements mnemosyne.RPCServer interface.
func (rs *rpcServer) SetValue(ctx context.Context, req *mnemosyne.SetValueRequest) (*mnemosyne.SetValueResponse, error) {
	h := rs.alloc.setValue(rs.logger, rs.storage, rs.monitor.rpc)
	h.monitor.requests.Add(1)

	bag, err := h.setValue(ctx, req)
	if err != nil {
		h.monitor.errors.Add(1)
		sklog.Error(h.logger, err)

		return nil, rs.error(err)
	}

	sklog.Debug(h.logger, "session bag value has been set")

	return &mnemosyne.SetValueResponse{
		Bag: bag,
	}, nil
}

// Delete implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Delete(ctx context.Context, req *mnemosyne.DeleteRequest) (*mnemosyne.DeleteResponse, error) {
	h := rs.alloc.delete(rs.logger, rs.storage, rs.monitor.rpc)
	h.monitor.requests.Add(1)

	affected, err := h.delete(ctx, req)
	if err != nil {
		h.monitor.errors.Add(1)
		sklog.Error(h.logger, err)

		return nil, rs.error(err)
	}

	sklog.Debug(h.logger, "session value has been deleted")

	return &mnemosyne.DeleteResponse{
		Count: affected,
	}, nil
}

func (rs *rpcServer) error(err error) error {
	if err == nil {
		return nil
	}

	switch err {
	case errSessionNotFound:
		return mnemosyne.ErrSessionNotFound
	default:
		return grpc.Errorf(codes.Internal, err.Error())
	}
}
