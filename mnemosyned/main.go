package main

import (
	"errors"
	"net"
	"os"
	"strconv"

	"github.com/piotrkowalczuk/mnemosyne"
	"github.com/piotrkowalczuk/sklog"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	config configuration
)

func init() {
	config.Init()
}

func main() {
	config.Parse()

	initLogger(config.logger.adapter, config.logger.format, config.logger.level, sklog.KeySubsystem, config.subsystem)
	initPostgres(
		config.storage.postgres.connectionString,
		config.storage.postgres.retry,
		logger,
	)

	hostname, err := os.Hostname()
	if err != nil {
		sklog.Fatal(logger, errors.New("mnemosyned: getting hostname failed"))
	}

	switch config.monitoring.engine {
	case "":
		sklog.Fatal(logger, errors.New("mnemosyned: monitoring is mandatory, at least for now..."))
	case monitoringEnginePrometheus:
		initMonitoring(initPrometheus(config.namespace, config.subsystem, prometheus.Labels{"server": hostname}), logger)
	default:
		sklog.Fatal(logger, errors.New("mnemosyned: unknown monitoring engine"))
	}

	switch config.storage.engine {
	case storageEngineInMemory:
		sklog.Fatal(logger, errors.New("mnemosyned: in memory storage is not implemented yet"))
	case storageEnginePostgres:
		initStorage(initPostgresStorage(config.storage.postgres.tableName, postgres, monitor), logger)
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
	//	if *tls {
	//		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	//		if err != nil {
	//			grpclog.Fatalf("Failed to generate credentials %v", err)
	//		}
	//		opts = []grpc.ServerOption{grpc.Creds(creds)}
	//	}
	grpclog.SetLogger(sklog.NewGRPCLogger(logger))
	gRPCServer := grpc.NewServer(opts...)

	mnemosyneServer := &rpcServer{
		logger:  logger,
		storage: storage,
	}
	mnemosyne.RegisterRPCServer(gRPCServer, mnemosyneServer)

	gRPCServer.Serve(listen)
}
