package mnemosyned

import (
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type sessionManagerStart struct {
	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (sms *sessionManagerStart) Start(ctx context.Context, req *mnemosynerpc.StartRequest) (*mnemosynerpc.StartResponse, error) {
	if req.Session == nil {
		return nil, errMissingSession
	}
	if req.Session.AccessToken == "" {
		var err error
		req.Session.AccessToken, err = mnemosyne.RandomAccessToken()
		if err != nil {
			return nil, grpc.Errorf(codes.Internal, "access token generation failure: %s", err.Error())
		}
	}

	if node, ok := sms.cluster.GetOther(req.Session.AccessToken); ok {
		sms.logger.Debug("start request forwarded", zap.String("to", node.Addr), zap.String("access_token", req.Session.AccessToken))
		return node.Client.Start(ctx, req)
	}

	if req.Session.SubjectId == "" {
		return nil, errMissingSubjectID
	}

	ses, err := sms.storage.Start(ctx,
		req.Session.AccessToken,
		req.Session.RefreshToken,
		req.Session.SubjectId,
		req.Session.SubjectClient,
		req.Session.Bag,
	)
	if err != nil {
		return nil, err
	}

	return &mnemosynerpc.StartResponse{
		Session: ses,
	}, nil
}
