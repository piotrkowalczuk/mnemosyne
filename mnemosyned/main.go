package main

import (
	"errors"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

var config configuration

func init() {
	config.init()
}

func main() {
	var (
		monitor *monitoring
		storage Storage
	)
	config.parse()

	logger := initLogger(config.logger.adapter, config.logger.format, config.logger.level, sklog.KeySubsystem, config.subsystem)
	postgres := initPostgres(
		config.storage.postgres.connectionString,
		logger,
	)

	hostname, err := os.Hostname()
	if err != nil {
		sklog.Fatal(logger, errors.New("mnemosyned: getting hostname failed"))
	}

	switch config.monitoring.engine {
	case "":
		sklog.Fatal(logger, errors.New("mnemosyned: monitoring is mandatory, at least for now"))
	case monitoringEnginePrometheus:
		monitor = initMonitoring(initPrometheus(config.namespace, config.subsystem, prometheus.Labels{"server": hostname}), logger)
	default:
		sklog.Fatal(logger, errors.New("mnemosyned: unknown monitoring engine"))
	}

	switch config.storage.engine {
	case storageEngineInMemory:
		sklog.Fatal(logger, errors.New("mnemosyned: in memory storage is not implemented yet"))
	case storageEnginePostgres:
		storage = initStorage(initPostgresStorage(config.storage.postgres.tableName, postgres, monitor), logger)
	case storageEngineRedis:
		sklog.Fatal(logger, errors.New("mnemosyned: redis storage is not implemented yet"))
	default:
		sklog.Fatal(logger, errors.New("mnemosyned: unknown storage engine"))
	}

	listenOn := config.host + ":" + strconv.FormatInt(int64(config.port), 10)
	listen, err := net.Listen("tcp", listenOn)
	if err != nil {
		sklog.Fatal(logger, err)
	}

	var opts []grpc.ServerOption
	if config.tls.enabled {
		creds, err := credentials.NewServerTLSFromFile(config.tls.certFile, config.tls.keyFile)
		if err != nil {
			sklog.Fatal(logger, err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))
	gRPCServer := grpc.NewServer(opts...)
	mnemosyneServer := newRPCServer(logger, storage, monitor)
	mnemosyne.RegisterRPCServer(gRPCServer, mnemosyneServer)

	sklog.Info(logger, "rpc api is running", "host", config.host, "port", config.port, "subsystem", config.subsystem, "namespace", config.namespace)

	go func() {
		sklog.Fatal(logger, http.ListenAndServe(address(config.host, config.port+1), nil))
	}()
	gRPCServer.Serve(listen)
}
