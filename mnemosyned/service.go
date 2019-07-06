package mnemosyned

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/piotrkowalczuk/mnemosyne/internal/cluster"
	"github.com/piotrkowalczuk/mnemosyne/internal/constant"
	"github.com/uber/jaeger-client-go/config"
	zapjaeger "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
)

func initCluster(l *zap.Logger, addr string, seeds ...string) (*cluster.Cluster, error) {
	csr, err := cluster.New(cluster.Opts{
		Listen: addr,
		Seeds:  seeds,
		Logger: l,
	})
	if err != nil {
		return nil, err
	}

	l.Debug("cluster initialized",
		zap.Strings("seeds", seeds),
		zap.String("listen", addr),
		zap.String("cluster", csr.GoString()),
	)
	return csr, nil
}

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service, node, agentAddress string, log *zap.Logger) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Tags: []opentracing.Tag{{
			Key:   constant.Subsystem + ".listen",
			Value: node,
		}},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentAddress,
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(zapjaeger.NewLogger(log)))
	if err != nil {
		return nil, nil, err
	}
	return tracer, closer, nil
}
