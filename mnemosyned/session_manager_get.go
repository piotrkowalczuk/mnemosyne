package mnemosyned

import (
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/jump"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
)

type sessionManagerGet struct {
	spanner

	storage storage.Storage
	cache   *cache.Cache
	cluster *cluster.Cluster
	logger  *zap.Logger
}

func (smg *sessionManagerGet) Get(ctx context.Context, req *mnemosynerpc.GetRequest) (*mnemosynerpc.GetResponse, error) {
	span, ctx := smg.span(ctx, "session-manager.get")
	defer span.Finish()

	if req.AccessToken == "" {
		return nil, errMissingAccessToken
	}
	if node, ok := smg.cluster.GetOther(req.AccessToken); ok {
		if cluster.IsInternalRequest(ctx) {
			return nil, status.Errorf(codes.FailedPrecondition,
				"it should be final destination of get request (%s), but found another node for it: %s",
				req.GetAccessToken(),
				node.Addr,
			)
		}
		smg.logger.Debug("get request forwarded", zap.String("remote_addr", node.Addr), zap.String("access_token", req.AccessToken))
		return node.Client.Get(ctx, req)
	}
	var (
		ses *mnemosynerpc.Session
		err error
	)

	hs := jump.Sum64(req.AccessToken)
	entry, ok := smg.cache.Read(hs)
	if !ok || (!entry.Refresh && time.Since(entry.Exp) > smg.cache.TTL) {
		if ok {
			smg.cache.Refresh(hs)
		}
		ses, err = smg.storage.Get(ctx, req.AccessToken)
		if err != nil {
			if err == storage.ErrSessionNotFound && ok {
				smg.cache.Del(hs)
			}
			return nil, err
		}
		smg.cache.Put(hs, *ses)
	} else {
		ses = &entry.Ses
	}

	return &mnemosynerpc.GetResponse{
		Session: ses,
	}, nil
}
