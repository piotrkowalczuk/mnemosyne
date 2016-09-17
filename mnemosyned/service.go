package mnemosyned

import (
	"database/sql"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func initPrometheus(namespace string, enabled bool, constLabels prometheus.Labels) *monitoring {
	generalErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   "rpc",
			Name:        "general_errors_total",
			Help:        "Total number of errors that happen during execution (other than grpc and postgres).",
			ConstLabels: constLabels,
		},
		monitoringGeneralLabels,
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

	if enabled {
		prometheus.MustRegisterOrGet(generalErrors)
		prometheus.MustRegisterOrGet(rpcRequests)
		prometheus.MustRegisterOrGet(rpcDuration)
		prometheus.MustRegisterOrGet(rpcErrors)
		prometheus.MustRegisterOrGet(postgresQueries)
		prometheus.MustRegisterOrGet(postgresErrors)
	}

	return &monitoring{
		enabled: enabled,
		general: monitoringGeneral{
			enabled: enabled,
			errors:  generalErrors,
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
	}
}

func initPostgres(address string, logger log.Logger) (*sql.DB, error) {
	postgres, err := sql.Open("postgres", address)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failure: %s", err.Error())
	}
	sklog.Info(logger, "postgres connection has been established", "address", address)

	return postgres, nil
}

func initBolt(path string) (bolt.DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("BoltDB opening failure: %s", err.Error())
	}
	sklog.Info(logger, "BoltDB database file openened", "path", path)

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
