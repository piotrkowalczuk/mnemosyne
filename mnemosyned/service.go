package mnemosyned

import (
	"database/sql"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
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
		prometheus.MustRegister(generalErrors)
		prometheus.MustRegister(rpcRequests)
		prometheus.MustRegister(rpcDuration)
		prometheus.MustRegister(rpcErrors)
		prometheus.MustRegister(postgresQueries)
		prometheus.MustRegister(postgresErrors)
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

func initStorage(env string, s storage, l log.Logger) (storage, error) {
	if env == EnvironmentTest {
		if err := s.TearDown(); err != nil {
			return nil, err
		}
	}
	if err := s.Setup(); err != nil {
		return nil, err
	}

	return s, nil
	//	switch e := err.(type) {
	//	case *pq.Error:
	//		sklog.Fatal(l, fmt.Errorf("storage setup failure: %s", e.Error()),
	//			"code", e.Code,
	//			"constraint", e.Constraint,
	//			"internal_query", e.InternalQuery,
	//			"column", e.Column,
	//			"detail", e.Detail,
	//			"hint", e.Hint,
	//			"line", e.Line,
	//			"schema", e.Schema,
	//		)
	//	default:
	//		sklog.Fatal(l, fmt.Errorf("storage setup failure: %s", e.Error()))
	//	}
	//}
	//
	//return s
}
