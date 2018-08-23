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
	flag.BoolVar(&c.grpc.debug, "grpc.debug", false, "If true, enables very verbose grpc debug mode. Useful to track connectivity issues.")
	flag.StringVar(&c.cluster.listen, "cluster.listen", "", "Complete instance address (including port).")
	flag.Var(&c.cluster.seeds, "cluster.seeds", "List of instances addresses that are pare of the cluster separated by the comma. Entry that overlaps with cluster.listen value will be ignored.")
	flag.StringVar(&c.catalog.http, "catalog.http", "http://localhost:8500/v1/catalog/service/mnemosyned", "Address of service catalog (exprimental).")
	flag.StringVar(&c.catalog.dns, "catalog.dns", "", "DNS server address that can resolve SRV lookup (experimental).")
	flag.DurationVar(&c.session.ttl, "ttl", storage.DefaultTTL, "Session time to live, after which session is deleted.")
	flag.DurationVar(&c.session.ttc, "ttc", storage.DefaultTTC, "Session time to cleanup, how often cleanup will be performed.")
	flag.StringVar(&c.logger.environment, "log.environment", "production", "Logger environment config (production, stackdriver or development).")
	flag.StringVar(&c.logger.level, "log.level", "info", "Logger level (debug, info, warn, error, dpanic, panic, fatal)")
	flag.StringVar(&c.storage, "storage", storage.EnginePostgres, "storage engine") // TODO: change to in memory when implemented
	flag.StringVar(&c.postgres.address, "postgres.address", "postgres://localhost?sslmode=disable", "storage postgres connection string")
	flag.StringVar(&c.postgres.table, "postgres.table", "session", "postgres table name")
	flag.StringVar(&c.postgres.schema, "postgres.schema", "mnemosyne", "postgres schema name")
	flag.BoolVar(&c.tls.enabled, "tls", false, "tls enable flag")
	flag.StringVar(&c.tls.certFile, "tls.crt", "", "path to tls cert file")
	flag.StringVar(&c.tls.keyFile, "tls.key", "", "path to tls key file")
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
