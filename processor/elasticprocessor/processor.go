package elasticprocessor

import (
	"context"
	"strings"

	"github.com/tommyers-elastic/opentelemetry-collector-contrib/processor/elasticprocessor/internal/hostmetrics"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"

	"go.uber.org/zap"
)

type ElasticProcessor struct {
	cfg     *Config
	logger  *zap.Logger
	storage map[string]any
}

func newProcessor(set processor.CreateSettings, cfg *Config) *ElasticProcessor {
	return &ElasticProcessor{cfg: cfg, logger: set.Logger, storage: make(map[string]any)}
}

func (p *ElasticProcessor) processMetrics(_ context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		resourceMetric := md.ResourceMetrics().At(i)
		rm := resourceMetric.Resource()

		for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
			scopeMetric := resourceMetric.ScopeMetrics().At(j)

			if p.cfg.AddSystemMetrics {
				if strings.HasPrefix(scopeMetric.Scope().Name(), "otelcol/hostmetricsreceiver") {
					if err := hostmetrics.AddElasticSystemMetrics(scopeMetric, rm, p.storage); err != nil {
						p.logger.Error("error adding the hostmetrics data", zap.Error(err))
					}
				}
			}
		}
	}

	return md, nil
}

func (p *ElasticProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	return ld, nil
}

func (p *ElasticProcessor) processTraces(_ context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	return td, nil
}
