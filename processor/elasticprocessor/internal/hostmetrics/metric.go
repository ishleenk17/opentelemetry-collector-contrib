package hostmetrics

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type dataType int

const (
	Gauge dataType = iota
	Sum
)

type metric struct {
	dataType       dataType
	name           string
	timestamp      pcommon.Timestamp
	startTimestamp pcommon.Timestamp
	intValue       *int64
	doubleValue    *float64
}

func addMetrics(ms pmetric.MetricSlice, resource pcommon.Resource, dataset string, oteldp pmetric.NumberDataPoint, metrics ...metric) {
	ms.EnsureCapacity(ms.Len() + len(metrics))

	for _, metric := range metrics {
		m := ms.AppendEmpty()
		m.SetName(metric.name)

		var dp pmetric.NumberDataPoint
		switch metric.dataType {
		case Gauge:
			dp = m.SetEmptyGauge().DataPoints().AppendEmpty()
		case Sum:
			dp = m.SetEmptySum().DataPoints().AppendEmpty()
		}

		if metric.intValue != nil {
			dp.SetIntValue(*metric.intValue)
		} else if metric.doubleValue != nil {
			dp.SetDoubleValue(*metric.doubleValue)
		}

		dp.SetTimestamp(metric.timestamp)
		if metric.startTimestamp != 0 {
			dp.SetStartTimestamp(metric.startTimestamp)
		}

		if dataset == "system.process" {
			// Add resource attribute as an attribute to each datapoint
			addProcessAttributes(resource, dp)
		}
		if dataset == "system.network" {
			// Add the network name as an attribute to each datapoint being generated
			addNetworkDatapointAttributes(oteldp, dp)

		}

		dp.Attributes().PutStr("data_stream.dataset", dataset)
	}
}
