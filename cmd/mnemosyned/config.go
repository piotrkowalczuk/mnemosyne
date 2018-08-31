package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"time"

	"github.com/piotrkowalczuk/mnemosyne/internal/storage"
)

var (
	version string
)

type configuration struct {
	host string
	port int
	grpc struct {
		debug bool
	}
	cluster struct {
		listen string
		seeds  arrayFlags
	}
	catalog struct {
		http string
		dns  string
	}
	tracing struct {
		agent struct {
			address string
		}
	}
	storage string
	logger  struct {
		environment string
		level       string
	}
	session struct {
		ttl time.Duration
		ttc time.Duration
	}
	postgres struct {
		address string
		table   string
		schema  string
	}
	tls struct {
		enabled  bool
		certFile string
		keyFile  string
	}
}

func (c *configuration) init() {
	if c == nil {
		*c = configuration{}
	}

	flag.StringVar(&c.host, "host", "127.0.0.1", "Host")
	flag.IntVar(&c.port, "port", 8080, "Port")
	// GRPC
	flag.BoolVar(&c.grpc.debug, "grpc.debug", false, "If true, enables very verbose gRPC to debug mode. Useful to track connectivity issues.")
	// CLUSTER
	flag.StringVar(&c.cluster.listen, "cluster.listen", "", "Complete instance address (including port).")
	flag.Var(&c.cluster.seeds, "cluster.seeds", "List of comma-separated instances addresses that are part of the cluster. An entry that overlaps with cluster.listen value will be ignored.")
	// CATALOG
	flag.StringVar(&c.catalog.http, "catalog.http", "http://localhost:8500/v1/catalog/service/mnemosyned", "Address of a service catalog (experimental).")
	flag.StringVar(&c.catalog.dns, "catalog.dns", "", "A DNS server address that can resolve SRV lookup (experimental).")
	// TRACING
	flag.StringVar(&c.tracing.agent.address, "tracing.agent.address", "", "Address of a tracing agent.")
	// SESSION
	flag.DurationVar(&c.session.ttl, "ttl", storage.DefaultTTL, "Session time to live, after which session is deleted.")
	flag.DurationVar(&c.session.ttc, "ttc", storage.DefaultTTC, "Session time to cleanup, how often cleanup will be performed.")
	// LOGGER
	flag.StringVar(&c.logger.environment, "log.environment", "production", "Logger environment config (production, stackdriver or development).")
	flag.StringVar(&c.logger.level, "log.level", "info", "Logger level (debug, info, warn, error, dpanic, panic, fatal)")
	// STORAGE
	flag.StringVar(&c.storage, "storage", storage.EnginePostgres, "Storage engine (postgres).") // TODO: change to in memory when implemented
	// POSTGRES
	flag.StringVar(&c.postgres.address, "postgres.address", "postgres://localhost?sslmode=disable", "Storage postgres connection string.")
	flag.StringVar(&c.postgres.table, "postgres.table", "session", "Postgres table name.")
	flag.StringVar(&c.postgres.schema, "postgres.schema", "mnemosyne", "Postgres schema name.")
	// TLS
	flag.BoolVar(&c.tls.enabled, "tls", false, "If true, TLS is enabled.")
	flag.StringVar(&c.tls.certFile, "tls.crt", "", "Path to TLS cert file.")
	flag.StringVar(&c.tls.keyFile, "tls.key", "", "Path to TLS key file.")
}

func (c *configuration) parse() {
	if !flag.Parsed() {
		ver := flag.Bool("version", false, "print version and exit")
		flag.Parse()
		if *ver {
			fmt.Printf("%s", version)
			os.Exit(0)
		}
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, strings.Split(value, ",")...)
	return nil
}
