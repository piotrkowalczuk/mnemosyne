package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-soa/mnemosyne/shared"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
)

type rpcServer struct {
	logger  log.Logger
	storage Storage
}

// Get ...
func (rs *rpcServer) Get(ctx context.Context, req *shared.GetRequest) (*shared.GetResponse, error) {
	ses, err := rs.storage.Get(req.Id)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "session has been retrieved", req.Context()...)

	return &shared.GetResponse{
		Session: ses,
	}, nil
}

func (rs *rpcServer) List(ctx context.Context, req *shared.ListRequest) (*shared.ListResponse, error) {
	return nil, nil
}

// Create ...
func (rs *rpcServer) Create(ctx context.Context, req *shared.CreateRequest) (*shared.CreateResponse, error) {
	ses, err := rs.storage.Create(req.Data)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "new session has been created", "id", ses.Id, "expire_at", ses.ExpireAt)

	return &shared.CreateResponse{
		Session: ses,
	}, nil
}

// Exists ...
func (rs *rpcServer) Exists(ctx context.Context, req *shared.ExistsRequest) (*shared.ExistsResponse, error) {
	ex, err := rs.storage.Exists(req.Id)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "session existance has been checked", "id", req.Id, "exists", ex)

	return &shared.ExistsResponse{
		Exists: ex,
	}, nil
}

// Abandon ...
func (rs *rpcServer) Abandon(ctx context.Context, req *shared.AbandonRequest) (*shared.AbandonResponse, error) {
	ab, err := rs.storage.Abandon(req.Id)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "session has been abandoned", "id", req.Id, "abadoned", ab)

	return &shared.AbandonResponse{
		Abandoned: ab,
	}, nil
}

// SetData ...
func (rs *rpcServer) SetData(ctx context.Context, req *shared.SetDataRequest) (*shared.SetDataResponse, error) {
	ses, err := rs.storage.SetData(req.Id, req.Key, req.Value)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "session data has been set", req.Context()...)

	return &shared.SetDataResponse{
		Session: ses,
	}, err
}

// Delete ...
func (rs *rpcServer) Delete(ctx context.Context, req *shared.DeleteRequest) (*shared.DeleteResponse, error) {
	expireAtFrom, err := shared.ParseTime(req.ExpireAtFrom)
	if err != nil {
		return nil, rs.error(err, req)
	}
	expireAtTo, err := shared.ParseTime(req.ExpireAtTo)
	if err != nil {
		return nil, rs.error(err, req)
	}

	count, err := rs.storage.Delete(req.Id, &expireAtFrom, &expireAtTo)
	if err != nil {
		return nil, rs.error(err, req)
	}

	sklog.Debug(rs.logger, "session(s) has been deleted", append(req.Context(), "count", count))

	return &shared.DeleteResponse{
		Count: count,
	}, err
}

func (rs *rpcServer) error(err error, ctx sklog.Contexter) error {
	sklog.Error(rs.logger, err, ctx.Context()...)

	return err
}
