package mnemosyned

import (
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const (
	// DefaultTTL is session time to live default value.
	DefaultTTL = 24 * time.Minute
	// DefaultTTC is time to clear default value.
	DefaultTTC = 1 * time.Minute
	// DefaultCacheSize determines how big cache should be at the beginning.
	DefaultCacheSize = 100000
)

var (
	errSessionNotFound    = grpc.Errorf(codes.NotFound, "mnemosyned: session not found")
	errMissingAccessToken = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing access token")
	errMissingSubjectID   = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing subject accessToken")
	errMissingSession     = grpc.Errorf(codes.InvalidArgument, "mnemosyned: missing session")
)

type sessionManagerOpts struct {
	addr       string
	cluster    *cluster.Cluster
	ttc        time.Duration
	logger     *zap.Logger
	storage    storage
	monitoring *monitoring
}

type sessionManager struct {
	ttc     time.Duration
	logger  *zap.Logger
	monitor monitoringCleanup
	storage storage

	sessionManagerList
	sessionManagerGet
	sessionManagerStart
	sessionManagerAbandon
	sessionManagerExists
	sessionManagerDelete
	sessionManagerSetValue
}

func newSessionManager(opts sessionManagerOpts) (*sessionManager, error) {
	cache := &cache{
		ttl:     5 * time.Second,
		monitor: opts.monitoring.cache,
		data:    make(map[uint64]*cacheEntry, DefaultCacheSize),
	}
	return &sessionManager{
		ttc:     opts.ttc,
		logger:  opts.logger,
		storage: opts.storage,
		monitor: opts.monitoring.cleanup,
		sessionManagerList: sessionManagerList{
			storage: opts.storage,
		},
		sessionManagerGet: sessionManagerGet{
			storage: opts.storage,
			cache:   cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerStart: sessionManagerStart{
			storage: opts.storage,
			cache:   cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerAbandon: sessionManagerAbandon{
			storage: opts.storage,
			cache:   cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerExists: sessionManagerExists{
			storage: opts.storage,
			cache:   cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerSetValue: sessionManagerSetValue{
			storage: opts.storage,
			cache:   cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerDelete: sessionManagerDelete{
			storage: opts.storage,
			cache:   cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
	}, nil
}

// Get implements RPCServer interface.
func (sm *sessionManager) Context(ctx context.Context, req *empty.Empty) (*mnemosynerpc.ContextResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing metadata in context, access token cannot be retrieved")
	}

	if len(md[mnemosyne.AccessTokenMetadataKey]) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing access token in metadata")
	}

	at := md[mnemosyne.AccessTokenMetadataKey][0]

	res, err := sm.Get(ctx, &mnemosynerpc.GetRequest{AccessToken: at})
	if err != nil {
		return nil, err
	}
	return &mnemosynerpc.ContextResponse{
		Session: res.Session,
	}, nil
}

func (sm *sessionManager) cleanup(done chan struct{}) {
	logger := sm.logger.Named("cleanup")
	sm.logger.Info("cleanup routing started")
InfLoop:
	for {
		select {
		case <-time.After(sm.ttc):
			t := time.Now()
			logger.Debug("session cleanup start", zap.Time("start_at", t))
			affected, err := sm.storage.Delete(context.Background(), "", "", "", nil, &t)
			if err != nil {
				if sm.monitor.enabled {
					sm.monitor.errors.Inc()
				}
				logger.Error("session cleanup failure", zap.Error(err), zap.Time("expire_at_to", t))
				return
			}

			logger.Debug("session cleanup success", zap.Int64("count", affected), zap.Duration("elapsed", time.Since(t)))
		case <-done:
			logger.Info("cleanup routing terminated")
			break InfLoop
		}
	}
}
