package mnemosyned

import (
	"github.com/opentracing/opentracing-go/log"
	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type sessionManagerSetValue struct {
	spanner

	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (smsv *sessionManagerSetValue) SetValue(ctx context.Context, req *mnemosynerpc.SetValueRequest) (*mnemosynerpc.SetValueResponse, error) {
	span, ctx := smsv.span(ctx, "session-manager.set-value")
	defer span.Finish()

	switch {
	case req.AccessToken == "":
		return nil, errMissingAccessToken
	case req.Key == "":
		return nil, status.Errorf(codes.InvalidArgument, "missing bag key")
	}

	if node, ok := smsv.cluster.GetOther(req.AccessToken); ok {
		if cluster.IsInternalRequest(ctx) {
			span.LogFields(
				log.String("error", "recursive internal call"),
				log.String("addr", node.Addr),
			)
			return nil, status.Errorf(codes.FailedPrecondition,
				"it should be final destination of set value request (%s), but found another node for it: %s",
				req.GetAccessToken(),
				node.Addr,
			)
		}
		smsv.logger.Debug("set value request forwarded", zap.String("remote_addr", node.Addr), zap.String("access_token", req.AccessToken))
		return node.Client.SetValue(ctx, req)
	}

	bag, err := smsv.storage.SetValue(ctx, req.AccessToken, req.Key, req.Value)
	if err != nil {
		return nil, err
	}

	return &mnemosynerpc.SetValueResponse{
		Bag: bag,
	}, nil
}
