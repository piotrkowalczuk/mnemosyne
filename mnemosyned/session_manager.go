package mnemosyned

import (
	"time"

	"github.com/opentracing/opentracing-go/log"

	"github.com/opentracing/opentracing-go"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/mnemosyne/internal/cache"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/constant"
	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingAccessToken = status.Errorf(codes.InvalidArgument, "mnemosyned: missing access token")
	errMissingSubjectID   = status.Errorf(codes.InvalidArgument, "mnemosyned: missing subject accessToken")
	errMissingSession     = status.Errorf(codes.InvalidArgument, "mnemosyned: missing session")
)

type sessionManagerOpts struct {
	addr    string
	cluster *cluster.Cluster
	cache   *cache.Cache
	ttc     time.Duration
	logger  *zap.Logger
	storage storage.Storage
	tracer  opentracing.Tracer
}

type sessionManager struct {
	ttc     time.Duration
	logger  *zap.Logger
	storage storage.Storage
	tracer  opentracing.Tracer
	// monitoring
	cleanupErrorsTotal prometheus.Counter

	sessionManagerList
	sessionManagerGet
	sessionManagerStart
	sessionManagerAbandon
	sessionManagerExists
	sessionManagerDelete
	sessionManagerSetValue
}

func newSessionManager(opts sessionManagerOpts) (*sessionManager, error) {
	spanner := spanner{tracer: opts.tracer}

	return &sessionManager{
		ttc:     opts.ttc,
		logger:  opts.logger,
		storage: opts.storage,
		tracer:  opts.tracer,
		cleanupErrorsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: constant.Subsystem,
				Subsystem: "cleanup",
				Name:      "errors_total",
				Help:      "Total number of errors that happen during cleanup.",
			},
		),
		sessionManagerList: sessionManagerList{
			spanner: spanner,
			storage: opts.storage,
		},
		sessionManagerGet: sessionManagerGet{
			spanner: spanner,
			storage: opts.storage,
			cache:   opts.cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerStart: sessionManagerStart{
			spanner: spanner,
			storage: opts.storage,
			cache:   opts.cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerAbandon: sessionManagerAbandon{
			spanner: spanner,
			storage: opts.storage,
			cache:   opts.cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerExists: sessionManagerExists{
			spanner: spanner,
			storage: opts.storage,
			cache:   opts.cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerSetValue: sessionManagerSetValue{
			spanner: spanner,
			storage: opts.storage,
			cache:   opts.cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
		sessionManagerDelete: sessionManagerDelete{
			storage: opts.storage,
			cache:   opts.cache,
			cluster: opts.cluster,
			logger:  opts.logger,
		},
	}, nil
}

// Context gets implements RPCServer interface.
func (sm *sessionManager) Context(ctx context.Context, req *empty.Empty) (*mnemosynerpc.ContextResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata in context, access token cannot be retrieved")
	}

	if len(md[mnemosyne.AccessTokenMetadataKey]) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing access token in metadata")
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
		span := sm.tracer.StartSpan("sessionManager.cleanup")

		select {
		case <-time.After(sm.ttc):
			t := time.Now()
			logger.Debug("session cleanup start", zap.Time("start_at", t))
			affected, err := sm.storage.Delete(opentracing.ContextWithSpan(context.Background(), span), "", "", "", nil, &t)
			if err != nil {
				sm.cleanupErrorsTotal.Inc()
				logger.Error("session cleanup failure", zap.Error(err), zap.Time("expire_at_to", t))
				span.LogFields(log.String("event", err.Error()))
				span.Finish()
				return
			}

			logger.Debug("session cleanup success", zap.Int64("count", affected), zap.Duration("elapsed", time.Since(t)))
		case <-done:
			logger.Info("cleanup routing terminated")
			span.Finish()
			break InfLoop
		}
		span.Finish()
	}
}

// Collect implements prometheus Collector interface.
func (sm *sessionManager) Collect(in chan<- prometheus.Metric) {
	sm.cleanupErrorsTotal.Collect(in)
}

// Describe implements prometheus Collector interface.
func (sm *sessionManager) Describe(in chan<- *prometheus.Desc) {
	sm.cleanupErrorsTotal.Describe(in)
}

type spanner struct {
	tracer opentracing.Tracer
}

func (s *spanner) span(ctx context.Context, operationName string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	var span opentracing.Span
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
		span = s.getTracer().StartSpan(operationName, opts...)
	} else {
		span = s.getTracer().StartSpan(operationName, opts...)
	}
	return span, opentracing.ContextWithSpan(ctx, span)
}

func (s *spanner) getTracer() opentracing.Tracer {
	if s.tracer == nil {
		return opentracing.GlobalTracer()
	}
	return s.tracer
}
