package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/piotrkowalczuk/mnemosyne/internal/service/logger"
	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapgrpc"
	"google.golang.org/grpc/grpclog"
)

var config configuration

func main() {
	config.init()
	config.parse()

	l, err := logger.Init(logger.Opts{
		Environment: config.logger.environment,
		Level:       zapcore.Level(config.logger.level),
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rpcListener := initListener(l, config.host, config.port)
	debugListener := initListener(l, config.host, config.port+1)

	daemon, err := mnemosyned.NewDaemon(&mnemosyned.DaemonOpts{
		SessionTTL:        config.session.ttl,
		SessionTTC:        config.session.ttc,
		Storage:           config.storage,
		Monitoring:        config.monitoring.enabled,
		PostgresAddress:   config.postgres.address + "&application_name=mnemosyned_" + version,
		PostgresTable:     config.postgres.table,
		PostgresSchema:    config.postgres.schema,
		TLS:               config.tls.enabled,
		TLSCertFile:       config.tls.certFile,
		TLSKeyFile:        config.tls.keyFile,
		ClusterListenAddr: config.cluster.listen,
		ClusterSeeds:      config.cluster.seeds,
		RPCListener:       rpcListener,
		Logger:            l.Named("daemon"),
		DebugListener:     debugListener,
	})
	if err != nil {
		l.Fatal("daemon allocation failure", zap.Error(err))
	}
	grpclog.SetLogger(zapgrpc.NewLogger(l))
	if err := daemon.Run(); err != nil {
		l.Fatal("daemon run failure", zap.Error(err))
	}
	defer daemon.Close()

	done := make(chan struct{})
	<-done
}

func initListener(logger *zap.Logger, host string, port int) net.Listener {
	on := host + ":" + strconv.FormatInt(int64(port), 10)
	listener, err := net.Listen("tcp", on)
	if err != nil {
		logger.Fatal("listener initialization failure", zap.Error(err), zap.String("host", host), zap.Int("port", port))
	}
	return listener
}
