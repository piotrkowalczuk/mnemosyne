package main

import (
	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"golang.org/x/net/context"
)

type rpcServer struct {
	logger  log.Logger
	storage Storage
}

// Get ...
func (rs *rpcServer) Get(ctx context.Context, req *mnemosyne.GetRequest) (*mnemosyne.GetResponse, error) {
	ses, err := rs.storage.Get(req.Id)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "session has been retrieved", req.Context()...)

	return &mnemosyne.GetResponse{
		Session: ses,
	}, nil
}

// List ...
func (rs *rpcServer) List(ctx context.Context, req *mnemosyne.ListRequest) (*mnemosyne.ListResponse, error) {
	return nil, nil
}

// Create ...
func (rs *rpcServer) Create(ctx context.Context, req *mnemosyne.CreateRequest) (*mnemosyne.CreateResponse, error) {
	ses, err := rs.storage.Create(req.Data)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "new session has been created", "id", ses.Id, "expire_at", ses.ExpireAt)

	return &mnemosyne.CreateResponse{
		Session: ses,
	}, nil
}

// Exists ...
func (rs *rpcServer) Exists(ctx context.Context, req *mnemosyne.ExistsRequest) (*mnemosyne.ExistsResponse, error) {
	ex, err := rs.storage.Exists(req.Id)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "session existance has been checked", "id", req.Id, "exists", ex)

	return &mnemosyne.ExistsResponse{
		Exists: ex,
	}, nil
}

// Abandon ...
func (rs *rpcServer) Abandon(ctx context.Context, req *mnemosyne.AbandonRequest) (*mnemosyne.AbandonResponse, error) {
	ab, err := rs.storage.Abandon(req.Id)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "session has been abandoned", "id", req.Id, "abadoned", ab)

	return &mnemosyne.AbandonResponse{
		Abandoned: ab,
	}, nil
}

// SetData ...
func (rs *rpcServer) SetData(ctx context.Context, req *mnemosyne.SetDataRequest) (*mnemosyne.SetDataResponse, error) {
	ses, err := rs.storage.SetData(req.Id, req.Key, req.Value)
	if err != nil {
		sklog.Error(rs.logger, err, req.Context()...)

		return nil, err
	}

	sklog.Debug(rs.logger, "session data has been set", req.Context()...)

	return &mnemosyne.SetDataResponse{
		Session: ses,
	}, err
}

// Delete ...
func (rs *rpcServer) Delete(ctx context.Context, req *mnemosyne.DeleteRequest) (*mnemosyne.DeleteResponse, error) {
	expireAtFrom, err := mnemosyne.ParseTime(req.ExpireAtFrom)
	if err != nil {
		return nil, rs.error(err, req)
	}
	expireAtTo, err := mnemosyne.ParseTime(req.ExpireAtTo)
	if err != nil {
		return nil, rs.error(err, req)
	}

	count, err := rs.storage.Delete(req.Id, &expireAtFrom, &expireAtTo)
	if err != nil {
		return nil, rs.error(err, req)
	}

	sklog.Debug(rs.logger, "session(s) has been deleted", append(req.Context(), "count", count))

	return &mnemosyne.DeleteResponse{
		Count: count,
	}, err
}

func (rs *rpcServer) error(err error, ctx sklog.Contexter) error {
	sklog.Error(rs.logger, err, ctx.Context()...)

	return err
}
