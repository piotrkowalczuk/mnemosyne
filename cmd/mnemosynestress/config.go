package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.0.0"

type configuration struct {
	verbose bool
	cluster struct {
		static struct {
			enabled bool
			members arrayFlags
		}
		discovery struct {
			enabled bool
			http    string
			dns     string
		}
	}
	tls struct {
		enabled bool
		cert    string
		key     string
	}
}

func (c *configuration) init() {
	if c == nil {
		*c = configuration{}
	}

	flag.BoolVar(&c.verbose, "verbose", false, "")
	flag.BoolVar(&c.cluster.static.enabled, "cluster.static", true, "")
	flag.Var(&c.cluster.static.members, "cluster.static.members", "")
	flag.BoolVar(&c.cluster.discovery.enabled, "cluster.discovery", false, "")
	flag.StringVar(&c.cluster.discovery.http, "cluster.discovery.http", "http://localhost:8500/v1/catalog/service/mnemosyned", "address of service catalog ")
	flag.StringVar(&c.cluster.discovery.dns, "cluster.discovery.dns", "", "dns server address that can resolve SRV lookup")
	flag.BoolVar(&c.tls.enabled, "tls", false, "tls enable flag")
	flag.StringVar(&c.tls.cert, "tls.cert", "", "path to tls cert file")
	flag.StringVar(&c.tls.key, "tls.key", "", "path to tls key file")
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
	*i = append(*i, value)
	return nil
}
