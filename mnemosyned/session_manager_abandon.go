package mnemosyned

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/jump"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type sessionManagerAbandon struct {
	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (sma *sessionManagerAbandon) Abandon(ctx context.Context, req *mnemosynerpc.AbandonRequest) (*wrappers.BoolValue, error) {
	if req.AccessToken == "" {
		return nil, errMissingAccessToken
	}

	if node, ok := sma.cluster.GetOther(req.AccessToken); ok {
		sma.logger.Debug("abandon request forwarded", zap.String("to", node.Addr), zap.String("access_token", req.AccessToken))
		return node.Client.Abandon(ctx, req)
	}

	sma.cache.Del(jump.Sum64(req.AccessToken))
	abandoned, err := sma.storage.Abandon(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}

	return &wrappers.BoolValue{Value: abandoned}, nil
}
