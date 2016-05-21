package mnemosyned

import "github.com/go-kit/kit/metrics"

var (
	monitoringGeneralLabels = []string{
		"action",
	}
	monitoringRPCLabels = []string{
		"method",
	}
	monitoringPostgresLabels = []string{
		"query",
	}
)

type monitoring struct {
	enabled  bool
	general  monitoringGeneral
	rpc      monitoringRPC
	postgres monitoringPostgres
}

type monitoringGeneral struct {
	enabled bool
	errors  metrics.Counter
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
