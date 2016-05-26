package mnemosyned

import (
	"database/sql"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/piotrkowalczuk/sklog"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	MonitoringEnginePrometheus = "prometheus"
)

func initPrometheus(namespace, subsystem string, constLabels stdprometheus.Labels) *monitoring {
	generalErrors := prometheus.NewCounter(
		stdprometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        "general_errors_total",
			Help:        "Total number of errors that happen during execution (other than grpc and postgres).",
			ConstLabels: constLabels,
		},
		monitoringGeneralLabels,
	)
	rpcRequests := prometheus.NewCounter(
		stdprometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        "rpc_requests_total",
			Help:        "Total number of RPC requests made.",
			ConstLabels: constLabels,
		},
		monitoringRPCLabels,
	)
	rpcErrors := prometheus.NewCounter(
		stdprometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        "rpc_errors_total",
			Help:        "Total number of errors that happen during RPC calles.",
			ConstLabels: constLabels,
		},
		monitoringRPCLabels,
	)

	postgresQueries := prometheus.NewCounter(
		stdprometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        "postgres_queries_total",
			Help:        "Total number of SQL queries made.",
			ConstLabels: constLabels,
		},
		monitoringPostgresLabels,
	)
	postgresErrors := prometheus.NewCounter(
		stdprometheus.CounterOpts{
			Namespace:   namespace,
			Subsystem:   subsystem,
			Name:        "postgres_errors_total",
			Help:        "Total number of errors that happen during SQL queries.",
			ConstLabels: constLabels,
		},
		monitoringPostgresLabels,
	)

	return &monitoring{
		enabled: true,
		general: monitoringGeneral{
			enabled: true,
			errors:  generalErrors,
		},
		rpc: monitoringRPC{
			enabled:  true,
			requests: rpcRequests,
			errors:   rpcErrors,
		},
		postgres: monitoringPostgres{
			enabled: true,
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

func initStorage(env string, s Storage, l log.Logger) (Storage, error) {
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
