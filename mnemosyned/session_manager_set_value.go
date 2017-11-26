package mnemosyned

import (
	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type sessionManagerSetValue struct {
	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (smsv *sessionManagerSetValue) SetValue(ctx context.Context, req *mnemosynerpc.SetValueRequest) (*mnemosynerpc.SetValueResponse, error) {
	switch {
	case req.AccessToken == "":
		return nil, errMissingAccessToken
	case req.Key == "":
		return nil, grpc.Errorf(codes.InvalidArgument, "missing bag key")
	}

	if node, ok := smsv.cluster.GetOther(req.AccessToken); ok {
		smsv.logger.Debug("set value request forwarded", zap.String("to", node.Addr), zap.String("access_token", req.AccessToken))
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
