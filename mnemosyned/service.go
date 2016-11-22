package mnemosyned

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
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

func initPostgres(address string, logger log.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", address)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failure: %s", err.Error())
	}

	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	username := ""
	if u.User != nil {
		username = u.User.Username()
	}

	// Otherwise 1 second cooldown is going to be multiplied by number of tests.
	if err := db.Ping(); err != nil {
		cancel := time.NewTimer(10 * time.Second)

	PingLoop:
		for {
			select {
			case <-time.After(1 * time.Second):
				if err := db.Ping(); err != nil {
					sklog.Debug(logger, "postgres connection ping failure", "postgres_host", u.Host, "postgres_user", username)
					continue PingLoop
				}
				sklog.Info(logger, "postgres connection has been established", "postgres_host", u.Host, "postgres_user", username)
				break PingLoop
			case <-cancel.C:
				return nil, errors.New("postgres connection timout")
			}
		}
	}

	sklog.Info(logger, "postgres connection has been established", "address", address)

	return db, nil
}

func initStorage(isTest bool, s storage, l log.Logger) (storage, error) {
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

func initCluster(l log.Logger, addr string, seeds ...string) (*cluster.Cluster, error) {
	csr, err := cluster.New(addr, seeds...)
	if err != nil {
		return nil, err
	}

	sklog.Debug(l, "cluster initialized", "seeds", len(seeds), "listen", addr)
	return csr, nil
}
