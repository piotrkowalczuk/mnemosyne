package main

import (
	"flag"
	"fmt"
	"os"

	"time"

	"github.com/piotrkowalczuk/mnemosyne/mnemosyned"
)

const version = "0.2.1"

type configuration struct {
	host    string
	port    int
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
	flag.DurationVar(&c.session.ttl, "ttl", mnemosyned.DefaultTTL, "session time to live, after which session is deleted")
	flag.DurationVar(&c.session.ttc, "ttc", mnemosyned.DefaultTTC, "session time to cleanup, how offten cleanup will be performed")
	flag.StringVar(&c.logger.adapter, "log.adapter", loggerAdapterStdOut, "logger adapter")
	flag.StringVar(&c.logger.format, "log.format", loggerFormatJSON, "logger format")
	flag.IntVar(&c.logger.level, "log.level", 6, "logger level")
	flag.BoolVar(&c.monitoring.enabled, "monitoring", false, "toggle application monitoring")
	flag.StringVar(&c.storage, "storage", mnemosyned.StorageEnginePostgres, "storage engine") // TODO: change to in memory when implemented
	flag.StringVar(&c.postgres.address, "postgres.address", "postgres://localhost?sslmode=disable", "storage postgres connection string")
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
