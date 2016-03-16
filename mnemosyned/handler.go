package main

import (
	"errors"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/piotrkowalczuk/mnemosyne"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type handlerFunc func(logger log.Logger, storage Storage, monitor monitoringRPC) *handler

type handler struct {
	logger  log.Logger
	storage Storage
	monitor monitoringRPC
}

func newHandlerFunc(endpoint string) handlerFunc {
	return func(logger log.Logger, storage Storage, monitor monitoringRPC) *handler {
		return &handler{
			logger:  log.NewContext(logger).With("endpoint", endpoint),
			storage: storage,
			monitor: monitoringRPC{
				errors:   monitor.errors.With(metrics.Field{Key: "endpoint", Value: endpoint}),
				requests: monitor.requests.With(metrics.Field{Key: "endpoint", Value: endpoint}),
			},
		}
	}
}

func (h *handler) context(ctx context.Context) (*mnemosyne.Session, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, errors.New("mnemosyned: missing metadata in context, session token cannot be retrieved")
	}

	if len(md[mnemosyne.AccessTokenMetadataKey]) == 0 {
		return nil, errors.New("mnemosyned: missing sesion token in metadata")
	}

	token := mnemosyne.DecodeAccessToken([]byte(md[mnemosyne.AccessTokenMetadataKey][0]))

	h.logger = log.NewContext(h.logger).With("access_token", token.Encode())

	return h.storage.Get(&token)
}

func (h *handler) get(ctx context.Context, req *mnemosyne.GetRequest) (*mnemosyne.Session, error) {
	if req.AccessToken == nil {
		return nil, mnemosyne.ErrMissingToken
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken.Encode())

	return h.storage.Get(req.AccessToken)
}

func (h *handler) list(ctx context.Context, req *mnemosyne.ListRequest) ([]*mnemosyne.Session, error) {
	expireAtFrom := req.ExpireAtFrom.Time()
	expireAtTo := req.ExpireAtTo.Time()

	h.logger = log.NewContext(h.logger).With(
		"offset", req.Offset,
		"limit", req.Limit,
		"expire_at_from", expireAtFrom.String(),
		"expire_at_to", expireAtTo.String(),
	)

	return h.storage.List(req.Offset, req.Limit, &expireAtFrom, &expireAtTo)
}

func (h *handler) start(ctx context.Context, req *mnemosyne.StartRequest) (*mnemosyne.Session, error) {
	if req.SubjectId == "" {
		return nil, mnemosyne.ErrMissingSubjectID
	}

	h.logger = log.NewContext(h.logger).With("subject_id", req.SubjectId)

	ses, err := h.storage.Start(req.SubjectId, req.Bag)
	if err != nil {
		return nil, err
	}

	h.logger = log.NewContext(h.logger).With("access_token", ses.AccessToken.Encode(), "expire_at", ses.ExpireAt.Time().Format(time.RFC3339))

	return ses, nil
}

func (h *handler) exists(ctx context.Context, req *mnemosyne.ExistsRequest) (bool, error) {
	if req.AccessToken == nil {
		return false, mnemosyne.ErrMissingToken
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken)

	exists, err := h.storage.Exists(req.AccessToken)
	if err != nil {
		return false, err
	}

	h.logger = log.NewContext(h.logger).With("exists", exists)

	return exists, nil
}

func (h *handler) abandon(ctx context.Context, req *mnemosyne.AbandonRequest) (bool, error) {
	if req.AccessToken == nil {
		return false, mnemosyne.ErrMissingToken
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken.Encode())

	abandoned, err := h.storage.Abandon(req.AccessToken)
	if err != nil {
		return false, err
	}

	return abandoned, nil
}

func (h *handler) setValue(ctx context.Context, req *mnemosyne.SetValueRequest) (map[string]string, error) {
	switch {
	case req.AccessToken == nil:
		return nil, mnemosyne.ErrMissingToken
	case req.Key == "":
		return nil, grpc.Errorf(codes.InvalidArgument, "mnemosyne: missing bag key")
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken.Encode(), "key", req.Key, "value", req.Value)

	bag, err := h.storage.SetValue(req.AccessToken, req.Key, req.Value)
	if err != nil {
		return nil, err
	}

	return bag, nil
}

func (h *handler) delete(ctx context.Context, req *mnemosyne.DeleteRequest) (int64, error) {
	expireAtFrom := req.ExpireAtFrom.Time()
	expireAtTo := req.ExpireAtTo.Time()

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken.Encode(), "expire_at_from", expireAtFrom, "expire_at_to", expireAtTo)

	affected, err := h.storage.Delete(req.AccessToken, &expireAtFrom, &expireAtTo)
	if err != nil {
		return 0, err
	}

	h.logger = log.NewContext(h.logger).With("affected", affected)

	return affected, nil
}
