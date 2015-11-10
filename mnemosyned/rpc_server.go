package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
)

type rpcServer struct {
	logger  log.Logger
	monitor *monitoring
	storage Storage
}

// Get ...
func (rs *rpcServer) Get(ctx context.Context, req *mnemosyne.GetRequest) (*mnemosyne.GetResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "get",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ses, err := rs.storage.Get(req.Id)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session has been retrieved", req.Context()...)

	return &mnemosyne.GetResponse{
		Session: ses,
	}, nil
}

// List ...
// TODO: implement
func (rs *rpcServer) List(ctx context.Context, req *mnemosyne.ListRequest) (*mnemosyne.ListResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "list",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	return nil, nil
}

// Create ...
func (rs *rpcServer) Create(ctx context.Context, req *mnemosyne.CreateRequest) (*mnemosyne.CreateResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "create",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ses, err := rs.storage.Create(req.Data)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "new session has been created", "id", ses.Id, "expire_at", ses.ExpireAt)

	return &mnemosyne.CreateResponse{
		Session: ses,
	}, nil
}

// Exists ...
func (rs *rpcServer) Exists(ctx context.Context, req *mnemosyne.ExistsRequest) (*mnemosyne.ExistsResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "exists",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ex, err := rs.storage.Exists(req.Id)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session existance has been checked", "id", req.Id, "exists", ex)

	return &mnemosyne.ExistsResponse{
		Exists: ex,
	}, nil
}

// Abandon ...
func (rs *rpcServer) Abandon(ctx context.Context, req *mnemosyne.AbandonRequest) (*mnemosyne.AbandonResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "abandon",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ab, err := rs.storage.Abandon(req.Id)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session has been abandoned", "id", req.Id, "abadoned", ab)

	return &mnemosyne.AbandonResponse{
		Abandoned: ab,
	}, nil
}

// SetData ...
func (rs *rpcServer) SetData(ctx context.Context, req *mnemosyne.SetDataRequest) (*mnemosyne.SetDataResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "set_data",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ses, err := rs.storage.SetData(req.Id, req.Key, req.Value)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session data has been set", req.Context()...)

	return &mnemosyne.SetDataResponse{
		Session: ses,
	}, err
}

// Delete ...
func (rs *rpcServer) Delete(ctx context.Context, req *mnemosyne.DeleteRequest) (*mnemosyne.DeleteResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "delete",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	expireAtFrom := req.ExpireAtFromTime()
	expireAtTo := req.ExpireAtToTime()

	count, err := rs.storage.Delete(req.Id, &expireAtFrom, &expireAtTo)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session(s) has been deleted", append(req.Context(), "count", count))

	return &mnemosyne.DeleteResponse{
		Count: count,
	}, err
}

func (rs *rpcServer) error(err error, field metrics.Field, ctx sklog.Contexter) error {
	rs.monitor.rpc.errors.With(field).Add(1)
	sklog.Error(rs.logger, err, ctx.Context()...)

	return err
}
