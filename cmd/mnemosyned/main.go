package main

import (
	"net"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
	"github.com/piotrkowalczuk/sklog"
	"google.golang.org/grpc/grpclog"
)

var config configuration

func main() {
	config.init()
	config.parse()

	logger := initLogger(config.logger.adapter, config.logger.format, config.logger.level)
	rpcListener := initListener(logger, config.host, config.port)
	debugListener := initListener(logger, config.host, config.port+1)

	daemon, err := mnemosyned.NewDaemon(&mnemosyned.DaemonOpts{
		SessionTTL:        config.session.ttl,
		SessionTTC:        config.session.ttc,
		Storage:           config.storage,
		Monitoring:        config.monitoring.enabled,
		PostgresAddress:   config.postgres.address,
		PostgresTable:     config.postgres.table,
		PostgresSchema:    config.postgres.schema,
		TLS:               config.tls.enabled,
		TLSCertFile:       config.tls.certFile,
		TLSKeyFile:        config.tls.keyFile,
		ClusterListenAddr: config.cluster.listen,
		ClusterSeeds:      config.cluster.seeds,
		RPCListener:       rpcListener,
		Logger:            logger,
		DebugListener:     debugListener,
	})
	if err != nil {
		sklog.Fatal(logger, err)
	}

	grpclog.SetLogger(sklog.NewGRPCLogger(logger))
	if err := daemon.Run(); err != nil {
		sklog.Fatal(logger, err)
	}
	defer daemon.Close()

	done := make(chan struct{})
	<-done
}

func initListener(logger log.Logger, host string, port int) net.Listener {
	on := host + ":" + strconv.FormatInt(int64(port), 10)
	listener, err := net.Listen("tcp", on)
	if err != nil {
		sklog.Fatal(logger, err)
	}
	return listener
}
