package mnemosyned

import (
	"errors"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
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
		h := &handler{
			logger:  log.NewContext(logger).With("endpoint", endpoint),
			storage: storage,
		}
		if monitor.enabled {
			h.monitor = monitoringRPC{
				errors:   monitor.errors.With(metrics.Field{Key: "endpoint", Value: endpoint}),
				requests: monitor.requests.With(metrics.Field{Key: "endpoint", Value: endpoint}),
			}
		}
		return h
	}
}

func (h *handler) error(err error) error {
	if err == nil {
		return nil
	}

	if h.monitor.enabled {
		h.monitor.errors.Add(1)
	}
	sklog.Error(h.logger, errors.New(grpc.ErrorDesc(err)), "grpc_code", grpc.Code(err))

	switch err {
	case ErrMissingAccessToken, ErrMissingSubjectID, ErrSessionNotFound:
		return err
	}
	code := grpc.Code(err)
	switch code {
	case codes.Unknown:
		return grpc.Errorf(codes.Internal, "mnemosyned: %s", grpc.ErrorDesc(err))
	default:
		return grpc.Errorf(code, "mnemosyned: %s", grpc.ErrorDesc(err))
	}
}

func (h *handler) context(ctx context.Context) (*mnemosyne.Session, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing metadata in context, session token cannot be retrieved")
	}

	if len(md[mnemosyne.AccessTokenMetadataKey]) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing sesion token in metadata")
	}

	token := mnemosyne.DecodeAccessToken([]byte(md[mnemosyne.AccessTokenMetadataKey][0]))

	h.logger = log.NewContext(h.logger).With("access_token", token.Encode())

	ses, err := h.storage.Get(&token)
	if err != nil {
		if err == ErrSessionNotFound {
			return nil, grpc.Errorf(codes.NotFound, "session (context) does not exists")
		}
		return nil, err
	}
	return ses, nil
}

func (h *handler) get(ctx context.Context, req *mnemosyne.GetRequest) (*mnemosyne.Session, error) {
	if req.AccessToken == nil {
		return nil, ErrMissingAccessToken
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken.Encode())

	ses, err := h.storage.Get(req.AccessToken)
	if err != nil {
		if err == ErrSessionNotFound {
			return nil, grpc.Errorf(codes.NotFound, "session (get) does not exists")
		}
		return nil, err
	}
	return ses, nil
}

func (h *handler) list(ctx context.Context, req *mnemosyne.ListRequest) ([]*mnemosyne.Session, error) {
	var (
		expireAtFrom, expireAtTo *time.Time
	)
	h.logger = log.NewContext(h.logger).With("offset", req.Offset, "limit", req.Limit)
	if req.ExpireAtFrom != nil {
		eaf, err := ptypes.Timestamp(req.ExpireAtFrom)
		if err != nil {
			return nil, err
		}
		expireAtFrom = &eaf
		h.logger = log.NewContext(h.logger).With("expire_at_from", eaf)
	}
	if req.ExpireAtTo != nil {
		eat, err := ptypes.Timestamp(req.ExpireAtTo)
		if err != nil {
			return nil, err
		}
		expireAtTo = &eat
		h.logger = log.NewContext(h.logger).With("expire_at_to", eat)
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	return h.storage.List(req.Offset, req.Limit, expireAtFrom, expireAtTo)
}

func (h *handler) start(ctx context.Context, req *mnemosyne.StartRequest) (*mnemosyne.Session, error) {
	if req.SubjectId == "" {
		return nil, ErrMissingSubjectID
	}

	h.logger = log.NewContext(h.logger).With("subject_id", req.SubjectId)

	ses, err := h.storage.Start(req.SubjectId, req.Bag)
	if err != nil {
		return nil, err
	}

	expireAt, err := ptypes.Timestamp(ses.ExpireAt)
	if err != nil {
		return nil, err
	}
	h.logger = log.NewContext(h.logger).With("access_token", ses.AccessToken.Encode(), "expire_at", expireAt)

	return ses, nil
}

func (h *handler) exists(ctx context.Context, req *mnemosyne.ExistsRequest) (bool, error) {
	if req.AccessToken == nil {
		return false, ErrMissingAccessToken
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
		return false, ErrMissingAccessToken
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
		return nil, ErrMissingAccessToken
	case req.Key == "":
		return nil, grpc.Errorf(codes.InvalidArgument, "missing bag key")
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken.Encode(), "key", req.Key, "value", req.Value)

	bag, err := h.storage.SetValue(req.AccessToken, req.Key, req.Value)
	if err != nil {
		return nil, err
	}

	return bag, nil
}

func (h *handler) delete(ctx context.Context, req *mnemosyne.DeleteRequest) (int64, error) {
	var (
		expireAtFrom, expireAtTo *time.Time
	)
	if req.AccessToken == nil && req.ExpireAtFrom == nil && req.ExpireAtTo == nil {
		return 0, grpc.Errorf(codes.InvalidArgument, "none of expected arguments was provided")
	}
	if req.AccessToken != nil {
		h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken.Encode())
	}
	if req.ExpireAtFrom != nil {
		eaf, err := ptypes.Timestamp(req.ExpireAtFrom)
		if err != nil {
			return 0, err
		}
		expireAtFrom = &eaf
		h.logger = log.NewContext(h.logger).With("expire_at_from", eaf)
	}
	if req.ExpireAtTo != nil {
		eat, err := ptypes.Timestamp(req.ExpireAtTo)
		if err != nil {
			return 0, err
		}
		expireAtTo = &eat
		h.logger = log.NewContext(h.logger).With("expire_at_to", eat)
	}

	affected, err := h.storage.Delete(req.AccessToken, expireAtFrom, expireAtTo)
	if err != nil {
		return 0, err
	}

	h.logger = log.NewContext(h.logger).With("affected", affected)

	return affected, nil
}
