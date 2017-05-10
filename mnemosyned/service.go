package mnemosyned

import (
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func initPrometheus(namespace string, enabled bool, constLabels prometheus.Labels) *monitoring {
	cleanupErrors := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   "cleanup",
			Name:        "errors_total",
			Help:        "Total number of errors that happen during cleanup.",
			ConstLabels: constLabels,
		},
	)
	rpcRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   "rpc",
			Name:        "requests_total",
			Help:        "Total number of RPC requests made.",
			ConstLabels: constLabels,
		},
		monitoringRPCLabels,
	)
	rpcDuration := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:   namespace,
			Subsystem:   "rpc",
			Name:        "request_duration_microseconds",
			Help:        "The RPC request latencies in microseconds.",
			ConstLabels: constLabels,
		},
		[]string{"handler", "code"},
	)
	rpcErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   "rpc",
			Name:        "errors_total",
			Help:        "Total number of errors that happen during RPC calles.",
			ConstLabels: constLabels,
		},
		monitoringRPCLabels,
	)

	postgresQueries := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   "storage",
			Name:        "postgres_queries_total",
			Help:        "Total number of SQL queries made.",
			ConstLabels: constLabels,
		},
		monitoringPostgresLabels,
	)
	postgresErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   "storage",
			Name:        "postgres_errors_total",
			Help:        "Total number of errors that happen during SQL queries.",
			ConstLabels: constLabels,
		},
		monitoringPostgresLabels,
	)
	cacheHits := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   "cache",
		Name:        "hits_total",
		Help:        "Total number of cache hits.",
		ConstLabels: constLabels,
	})
	cacheMisses := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   "cache",
		Name:        "misses_total",
		Help:        "Total number of cache misses.",
		ConstLabels: constLabels,
	})
	cacheRefresh := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Subsystem:   "cache",
		Name:        "refresh_total",
		Help:        "Total number of times cache refresh.",
		ConstLabels: constLabels,
	})

	if enabled {
		cleanupErrors = prometheus.MustRegisterOrGet(cleanupErrors).(prometheus.Counter)
		rpcRequests = prometheus.MustRegisterOrGet(rpcRequests).(*prometheus.CounterVec)
		rpcDuration = prometheus.MustRegisterOrGet(rpcDuration).(*prometheus.SummaryVec)
		rpcErrors = prometheus.MustRegisterOrGet(rpcErrors).(*prometheus.CounterVec)
		postgresQueries = prometheus.MustRegisterOrGet(postgresQueries).(*prometheus.CounterVec)
		postgresErrors = prometheus.MustRegisterOrGet(postgresErrors).(*prometheus.CounterVec)
		cacheHits = prometheus.MustRegisterOrGet(cacheHits).(prometheus.Counter)
		cacheMisses = prometheus.MustRegisterOrGet(cacheMisses).(prometheus.Counter)
		cacheRefresh = prometheus.MustRegisterOrGet(cacheRefresh).(prometheus.Counter)
	}

	return &monitoring{
		enabled: enabled,
		cleanup: monitoringCleanup{
			enabled: enabled,
			errors:  cleanupErrors,
		},
		rpc: monitoringRPC{
			enabled:  enabled,
			duration: rpcDuration,
			requests: rpcRequests,
			errors:   rpcErrors,
		},
		postgres: monitoringPostgres{
			enabled: enabled,
			queries: postgresQueries,
			errors:  postgresErrors,
		},
		cache: monitoringCache{
			enabled: enabled,
			hits:    cacheHits,
			misses:  cacheMisses,
			refresh: cacheRefresh,
		},
	}
}

func initStorage(isTest bool, s storage) (storage, error) {
	if isTest {
		if err := s.TearDown(); err != nil {
			return nil, err
		}
	}
	if err := s.Setup(); err != nil {
		return nil, err
	}

	return s, nil
}

func initCluster(l *zap.Logger, addr string, seeds ...string) (*cluster.Cluster, error) {
	csr, err := cluster.New(cluster.Opts{
		Listen: addr,
		Seeds:  seeds,
		Logger: l,
	})
	if err != nil {
		return nil, err
	}

	l.Debug("cluster initialized", zap.Strings("seeds", seeds), zap.String("listen", addr))
	return csr, nil
}
