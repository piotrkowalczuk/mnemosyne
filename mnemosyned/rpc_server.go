package main

import (
	"errors"

	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type rpcServer struct {
	logger  log.Logger
	monitor *monitoring
	storage Storage
}

func (rs *rpcServer) Context(ctx context.Context, req *mnemosyne.Empty) (*mnemosyne.Session, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "context",
	}
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, errors.New("mnemosyned: missing metadata in context, session token cannot be retrieved")
	}

	if len(md[mnemosyne.TokenMetadataKey]) == 0 {
		return nil, errors.New("mnemosyned: missing sesion token in metadata")
	}

	token := mnemosyne.DecodeToken([]byte(md[mnemosyne.TokenMetadataKey][0]))

	ses, err := rs.storage.Get(&token)
	if err != nil {
		return nil, rs.error(err, field, nil)
	}

	sklog.Debug(rs.logger, "session has been retrieved", "endpoint", "context", "token", md[mnemosyne.TokenMetadataKey][0])

	return ses, nil
}

// Get implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Get(ctx context.Context, req *mnemosyne.GetRequest) (*mnemosyne.GetResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "get",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ses, err := rs.storage.Get(req.Token)
	if err != nil {
		if err == errSessionNotFound {
			return nil, mnemosyne.ErrSessionNotFound
		}
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session has been retrieved", req.Context()...)

	return &mnemosyne.GetResponse{
		Session: ses,
	}, nil
}

// List implements mnemosyne.RPCServer interface.
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

// Start implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Start(ctx context.Context, req *mnemosyne.StartRequest) (*mnemosyne.StartResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "create",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	if req.SubjectId == "" {
		return nil, errors.New("mnemosyned: session cannot be started, subject id is missing")
	}
	ses, err := rs.storage.Start(req.SubjectId, req.Bag)
	if err != nil {
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "new session has been created", "token", ses.Token, "expire_at", ses.ExpireAt.Time().Format(time.RFC3339))

	return &mnemosyne.StartResponse{
		Session: ses,
	}, nil
}

// Exists implements mnemosyne.RPCServer interface.
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

// Abandon implements mnemosyne.RPCServer interface.
func (rs *rpcServer) Abandon(ctx context.Context, req *mnemosyne.AbandonRequest) (*mnemosyne.AbandonResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "abandon",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	ab, err := rs.storage.Abandon(req.Token)
	if err != nil {
		if err == errSessionNotFound {
			return nil, mnemosyne.ErrSessionNotFound
		}
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session has been abandoned", "token", req.Token, "abadoned", ab)

	return &mnemosyne.AbandonResponse{
		Abandoned: ab,
	}, nil
}

// SetValue implements mnemosyne.RPCServer interface.
func (rs *rpcServer) SetValue(ctx context.Context, req *mnemosyne.SetValueRequest) (*mnemosyne.SetValueResponse, error) {
	field := metrics.Field{
		Key:   "method",
		Value: "set_value",
	}
	rs.monitor.rpc.requests.With(field).Add(1)

	bag, err := rs.storage.SetValue(req.Token, req.Key, req.Value)
	if err != nil {
		if err == errSessionNotFound {
			return nil, mnemosyne.ErrSessionNotFound
		}
		return nil, rs.error(err, field, req)
	}

	sklog.Debug(rs.logger, "session value has been set", req.Context()...)

	return &mnemosyne.SetValueResponse{
		Bag: bag,
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

// Delete implements mnemosyne.RPCServer interface.
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
	if err == nil {
		return nil
	}

	rs.monitor.rpc.errors.With(field).Add(1)

	switch err {
	case errSessionNotFound:
		return mnemosyne.ErrSessionNotFound
	default:
		sklog.Error(rs.logger, err, ctx.Context()...)
		return grpc.Errorf(codes.Internal, err.Error())
	}
}
