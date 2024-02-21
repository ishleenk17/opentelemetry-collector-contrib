package hostmetrics

import (
	"fmt"
	"path"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var scraperToElasticDataset = map[string]string{
	"cpu":        "system.cpu",
	"disk":       "system.diskio",
	"filesystem": "system.filesystem",
	"load":       "system.load",
	"memory":     "system.memory",
	"network":    "system.network",
	"paging":     "system.memory",
	"processes":  "system.process.summary",
	"process":    "system.process",
}

// AddElasticSystemMetrics computes additional metrics for compatibility with the Elastic system integration.
// The `scopeMetrics` input should be metrics generated by a specific hostmetrics scraper.
// `scopeMetrics` are modified in place.
func AddElasticSystemMetrics(scopeMetrics pmetric.ScopeMetrics, rm pcommon.Resource, storage map[string]any) error {
	scope := scopeMetrics.Scope()
	scraper := path.Base(scope.Name())

	dataset, ok := scraperToElasticDataset[scraper]
	if !ok {
		return fmt.Errorf("no dataset defined for scaper '%s'", scraper)
	}

	currentTime := time.Now().UnixMilli()
	if lastScrape, ok := storage["lastScrape"]; ok {
		collectionPeriod := currentTime - lastScrape.(int64)
		scopeMetrics.Scope().Attributes().PutDouble("metricset.period", float64(collectionPeriod))
	}
	storage["lastScrape"] = currentTime

	switch scraper {
	case "cpu":
		return addCPUMetrics(scopeMetrics.Metrics(), rm, dataset)
	case "memory":
		return addMemoryMetrics(scopeMetrics.Metrics(), rm, dataset)
	case "load":
		return addLoadMetrics(scopeMetrics.Metrics(), rm, dataset)
	case "process":
		return addProcessMetrics(scopeMetrics.Metrics(), rm, dataset)
	case "processes":
		return addProcessSummaryMetrics(scopeMetrics.Metrics(), rm, dataset)
	default:
		return fmt.Errorf("no matching transform function found for scope '%s'", scope.Name())
	}
}
