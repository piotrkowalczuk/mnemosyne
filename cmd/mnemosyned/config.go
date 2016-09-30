package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"time"

	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
)

var (
	version string
)

type configuration struct {
	host    string
	port    int
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
		adapter string
		format  string
		level   int
	}
	session struct {
		ttl time.Duration
		ttc time.Duration
	}
	monitoring struct {
		enabled bool
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

	flag.StringVar(&c.host, "host", "127.0.0.1", "host")
	flag.IntVar(&c.port, "port", 8080, "port")
	flag.StringVar(&c.cluster.listen, "cluster.listen", "", "")
	flag.Var(&c.cluster.seeds, "cluster.seeds", "")
	flag.StringVar(&c.catalog.http, "catalog.http", "http://localhost:8500/v1/catalog/service/mnemosyned", "address of service catalog ")
	flag.StringVar(&c.catalog.dns, "catalog.dns", "", "dns server address that can resolve SRV lookup")
	flag.DurationVar(&c.session.ttl, "ttl", mnemosyned.DefaultTTL, "session time to live, after which session is deleted")
	flag.DurationVar(&c.session.ttc, "ttc", mnemosyned.DefaultTTC, "session time to cleanup, how offten cleanup will be performed")
	flag.StringVar(&c.logger.adapter, "log.adapter", loggerAdapterStdOut, "logger adapter")
	flag.StringVar(&c.logger.format, "log.format", loggerFormatJSON, "logger format")
	flag.IntVar(&c.logger.level, "log.level", 6, "logger level")
	flag.BoolVar(&c.monitoring.enabled, "monitoring", false, "toggle application monitoring")
	flag.StringVar(&c.storage, "storage", mnemosyned.StorageEnginePostgres, "storage engine") // TODO: change to in memory when implemented
	flag.StringVar(&c.postgres.address, "postgres.address", "postgres://localhost?sslmode=disable", "storage postgres connection string")
	flag.StringVar(&c.postgres.table, "postgres.table", "session", "postgres table name")
	flag.StringVar(&c.postgres.schema, "postgres.schema", "mnemosyne", "postgres schema name")
	flag.BoolVar(&c.tls.enabled, "tls", false, "tls enable flag")
	flag.StringVar(&c.tls.certFile, "tls.cert", "", "path to tls cert file")
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
