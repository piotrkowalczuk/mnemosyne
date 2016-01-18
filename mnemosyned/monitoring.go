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
	rpc      monitoringRPC
	postgres monitoringPostgres
}

type monitoringRPC struct {
	requests metrics.Counter
	errors   metrics.Counter
}

type monitoringPostgres struct {
	queries metrics.Counter
	errors  metrics.Counter
}
