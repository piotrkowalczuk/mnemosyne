package mnemosyned

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
	enabled  bool
	rpc      monitoringRPC
	postgres monitoringPostgres
}

type monitoringRPC struct {
	enabled  bool
	requests metrics.Counter
	errors   metrics.Counter
}

type monitoringPostgres struct {
	enabled bool
	queries metrics.Counter
	errors  metrics.Counter
}
