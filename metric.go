package prompt

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Metric struct {
	Name      string
	Collector prometheus.Collector
}

func InjectCounter(module core.RefProvider, name string) prometheus.Counter {
	metric, ok := module.Ref(GetMetricName(name)).(prometheus.Counter)
	if !ok {
		return nil
	}

	return metric
}

func InjectCounterVec(module core.RefProvider, name string) *prometheus.CounterVec {
	metric, ok := module.Ref(GetMetricName(name)).(*prometheus.CounterVec)
	if !ok {
		return nil
	}

	return metric
}

func InjectGauge(module core.RefProvider, name string) prometheus.Gauge {
	metric, ok := module.Ref(GetMetricName(name)).(prometheus.Gauge)
	if !ok {
		return nil
	}

	return metric
}

func InjectGaugeVec(module core.RefProvider, name string) *prometheus.GaugeVec {
	metric, ok := module.Ref(GetMetricName(name)).(*prometheus.GaugeVec)
	if !ok {
		return nil
	}

	return metric
}

func InjectHistogram(module core.RefProvider, name string) prometheus.Histogram {
	metric, ok := module.Ref(GetMetricName(name)).(prometheus.Histogram)
	if !ok {
		return nil
	}

	return metric
}

func InjectHistogramVec(module core.RefProvider, name string) *prometheus.HistogramVec {
	metric, ok := module.Ref(GetMetricName(name)).(*prometheus.HistogramVec)
	if !ok {
		return nil
	}

	return metric
}

func InjectSummary(module core.RefProvider, name string) prometheus.Summary {
	metric, ok := module.Ref(GetMetricName(name)).(prometheus.Summary)
	if !ok {
		return nil
	}

	return metric
}

func InjectSummaryVec(module core.RefProvider, name string) *prometheus.SummaryVec {
	metric, ok := module.Ref(GetMetricName(name)).(*prometheus.SummaryVec)
	if !ok {
		return nil
	}

	return metric
}
