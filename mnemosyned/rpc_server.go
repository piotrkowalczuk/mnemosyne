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

	ses, err := rs.storage.Get(req.Token)
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

	expireAtFrom := req.ExpireAtFrom.Time()
	expireAtTo := req.ExpireAtTo.Time()

	rs.monitor.rpc.requests.With(field).Add(1)
	sessions, err := rs.storage.List(req.Offset, req.Limit, &expireAtFrom, &expireAtTo)
	if err != nil {
		return nil, rs.error(err, field, req)
	}
	return &mnemosyne.ListResponse{
		Sessions: sessions,
	}, nil
}

// Start ...
func (rs *rpcServer) Start(ctx context.Context, req *mnemosyne.StartRequest) (*mnemosyne.StartResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "create",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ses, err := rs.storage.Start(req.SubjectId, req.Bag)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "new session has been created", "token", ses.Token, "expire_at", ses.ExpireAt)

	return &mnemosyne.StartResponse{
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

	ex, err := rs.storage.Exists(req.Token)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session existance has been checked", "token", req.Token, "exists", ex)

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

	ab, err := rs.storage.Abandon(req.Token)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session has been abandoned", "token", req.Token, "abadoned", ab)

	return &mnemosyne.AbandonResponse{
		Abandoned: ab,
	}, nil
}

// SetValue ...
func (rs *rpcServer) SetValue(ctx context.Context, req *mnemosyne.SetValueRequest) (*mnemosyne.SetValueResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "set_value",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ses, err := rs.storage.SetValue(req.Token, req.Key, req.Value)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session value has been set", req.Context()...)

	return &mnemosyne.SetValueResponse{
		Session: ses,
	}, err
}

//// DeleteValue ...
//func (rs *rpcServer) DeleteValue(ctx context.Context, req *mnemosyne.DeleteValueRequest) (*mnemosyne.DeleteValueResponse, error) {
//	field := metrics.Field{
//		Key:   "method",
//		Value: "delete_value",
//	}
//	rs.monitor.rpc.requests.With(field).Add(1)
//
//	ses, err := rs.storage.DeleteValue(req.Token, req.Key)
//	if err != nil {
//		return nil, rs.error(err, field, req)
//	}
//
//	sklog.Debug(rs.logger, "session value has been deleted", req.Context()...)
//
//	return &mnemosyne.DeleteValueResponse{
//		Session: ses,
//	}, err
//}

// Delete ...
func (rs *rpcServer) Delete(ctx context.Context, req *mnemosyne.DeleteRequest) (*mnemosyne.DeleteResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "delete",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	expireAtFrom := req.ExpireAtFrom.Time()
	expireAtTo := req.ExpireAtTo.Time()

	count, err := rs.storage.Delete(req.Token, &expireAtFrom, &expireAtTo)
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
