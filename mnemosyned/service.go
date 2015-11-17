package main

import (
	"database/sql"
	stdlog "log"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
	"github.com/piotrkowalczuk/sklog"
)

func initPostgres(connectionString string, retry int, logger log.Logger) *sql.DB {
	var err error
	var postgres *sql.DB

	// Because of recursion it needs to be checked to not spawn more than one.
	if postgres == nil {
		postgres, err = sql.Open("postgres", connectionString)
		if err != nil {
			sklog.Fatal(logger, err)
		}
	}

	// At this moment connection is not yet established.
	// Ping is required.
RetryLoop:
	for i := 0; i <= retry; i++ {
		if err := postgres.Ping(); err != nil {
			if i == retry {
				sklog.Fatal(logger, err)
			}

			sklog.Error(logger, err)
			time.Sleep(2 * time.Second)
		}
		break RetryLoop
	}

	sklog.Info(logger, "connection do postgres established successfully", "address", connectionString, "retry", retry)

	return postgres
}

func initMonitoring(fn func() (*monitoring, error), logger log.Logger) *monitoring {
	m, err := fn()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	return m
}

const (
	loggerAdapterStdOut = "stdout"
	loggerFormatJSON    = "json"
	loggerFormatHumane  = "humane"
	loggerFormatLogFmt  = "logfmt"
)

func initLogger(adapter, format string, level int, context ...interface{}) log.Logger {
	var l log.Logger

	if adapter != loggerAdapterStdOut {
		stdlog.Fatal("service: unsupported logger adapter")
	}

	switch format {
	case loggerFormatHumane:
		l = sklog.NewHumaneLogger(os.Stdout)
	case loggerFormatJSON:
		l = log.NewJSONLogger(os.Stdout)
	case loggerFormatLogFmt:
		l = log.NewLogfmtLogger(os.Stdout)
	default:
		stdlog.Fatal("mnemosyned: unsupported logger format")
	}

	l = log.NewContext(l).With(context...)

	sklog.Info(l, "logger has been initialized successfully", "adapter", adapter, "format", format, "level", level)

	return l
}

func initStorage(fn func() (Storage, error), logger log.Logger) Storage {
	s, err := fn()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	err = s.Setup()
	if err != nil {
		sklog.Fatal(logger, err)
	}

	return s
}
