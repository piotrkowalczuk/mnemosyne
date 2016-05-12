package mnemosyned

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/sklog"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const (
	MonitoringEnginePrometheus = "prometheus"
)

func initPrometheus(namespace, subsystem string, constLabels stdprometheus.Labels) *monitoring {
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

func initPostgres(connectionString string, logger log.Logger) (*sql.DB, error) {
	postgres, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failure: %s", err.Error())
	}

	sklog.Info(logger, "postgres connected", "address", connectionString)

	return postgres, nil
}

const (
	LoggerAdapterStdOut = "stdout"
	LoggerAdapterNone   = "none"
	LoggerFormatJSON    = "json"
	LoggerFormatHumane  = "humane"
	LoggerFormatLogFmt  = "logfmt"
)

func initLogger(adapter, format string, level int, context ...interface{}) log.Logger {
	var (
		l log.Logger
		a io.Writer
	)

	switch adapter {
	case LoggerAdapterStdOut:
		a = os.Stdout
	case LoggerAdapterNone:
		a = ioutil.Discard
	default:
		stdlog.Fatal("unsupported logger adapter")
	}

	switch format {
	case LoggerFormatHumane:
		l = sklog.NewHumaneLogger(a, sklog.DefaultHTTPFormatter)
	case LoggerFormatJSON:
		l = log.NewJSONLogger(a)
	case LoggerFormatLogFmt:
		l = log.NewLogfmtLogger(a)
	default:
		stdlog.Fatal("unsupported logger format")
	}

	l = log.NewContext(l).With(context...)

	sklog.Info(l, "logger has been initialized successfully", "adapter", adapter, "format", format, "level", level)

	return l
}

func initStorage(s Storage, l log.Logger) (Storage, error) {
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
