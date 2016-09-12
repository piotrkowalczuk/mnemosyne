package mnemosyned

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/golang/protobuf/ptypes"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type handlerFunc func(logger log.Logger, storage storage, monitor monitoringRPC) *handler

type handler struct {
	logger  log.Logger
	storage storage
	monitor monitoringRPC
}

func newHandlerFunc(endpoint string) handlerFunc {
	return func(logger log.Logger, storage storage, monitor monitoringRPC) *handler {
		h := &handler{
			logger:  log.NewContext(logger).With("endpoint", endpoint),
			storage: storage,
		}
		return h
	}
}

func (h *handler) get(ctx context.Context, req *mnemosynerpc.GetRequest) (*mnemosynerpc.Session, error) {
	if req.AccessToken == "" {
		return nil, errMissingAccessToken
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken)

	ses, err := h.storage.Get(req.AccessToken)
	if err != nil {
		if err == errSessionNotFound {
			return nil, grpc.Errorf(codes.NotFound, "session does not exists")
		}
		return nil, err
	}
	return ses, nil
}

func (h *handler) list(ctx context.Context, req *mnemosynerpc.ListRequest) ([]*mnemosynerpc.Session, error) {
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

func (h *handler) start(ctx context.Context, req *mnemosynerpc.StartRequest) (*mnemosynerpc.Session, error) {
	if req.Session.SubjectId == "" {
		return nil, errMissingSubjectID
	}

	h.logger = log.NewContext(h.logger).With("subject_id", req.Session.SubjectId)

	ses, err := h.storage.Start(req.Session.AccessToken, req.Session.SubjectId, req.Session.SubjectClient, req.Session.Bag)
	if err != nil {
		return nil, err
	}

	expireAt, err := ptypes.Timestamp(ses.ExpireAt)
	if err != nil {
		return nil, err
	}
	h.logger = log.NewContext(h.logger).With("access_token", ses.AccessToken, "expire_at", expireAt)

	return ses, nil
}

func (h *handler) exists(ctx context.Context, req *mnemosynerpc.ExistsRequest) (bool, error) {
	if req.AccessToken == "" {
		return false, errMissingAccessToken
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken)

	exists, err := h.storage.Exists(req.AccessToken)
	if err != nil {
		return false, err
	}

	h.logger = log.NewContext(h.logger).With("exists", exists)

	return exists, nil
}

func (h *handler) abandon(ctx context.Context, req *mnemosynerpc.AbandonRequest) (bool, error) {
	if req.AccessToken == "" {
		return false, errMissingAccessToken
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken)

	abandoned, err := h.storage.Abandon(req.AccessToken)
	if err != nil {
		return false, err
	}

	return abandoned, nil
}

func (h *handler) setValue(ctx context.Context, req *mnemosynerpc.SetValueRequest) (map[string]string, error) {
	switch {
	case req.AccessToken == "":
		return nil, errMissingAccessToken
	case req.Key == "":
		return nil, grpc.Errorf(codes.InvalidArgument, "missing bag key")
	}

	h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken, "key", req.Key, "value", req.Value)

	bag, err := h.storage.SetValue(req.AccessToken, req.Key, req.Value)
	if err != nil {
		return nil, err
	}

	return bag, nil
}

func (h *handler) delete(ctx context.Context, req *mnemosynerpc.DeleteRequest) (int64, error) {
	var (
		expireAtFrom, expireAtTo *time.Time
	)
	if req.AccessToken == "" && req.ExpireAtFrom == nil && req.ExpireAtTo == nil {
		return 0, grpc.Errorf(codes.InvalidArgument, "none of expected arguments was provided")
	}
	if req.AccessToken != "" {
		h.logger = log.NewContext(h.logger).With("access_token", req.AccessToken)
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
