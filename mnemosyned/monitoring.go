package main

import "github.com/go-kit/kit/metrics"

var (
	monitoringRPCLabels = []string{
		"method",
	}
	monitoringPostgresLabels = []string{
		"query",
	}
)

type monitoring struct {
	rpc struct {
		requests metrics.Counter
		errors   metrics.Counter
	}
	postgres struct {
		queries metrics.Counter
		errors  metrics.Counter
	}
}
