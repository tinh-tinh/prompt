package prompt

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

// Counter
func MakeCounterProvider(opt prometheus.CounterOpts, labelNames []string) core.Providers {
	counter := prometheus.NewCounterVec(opt, labelNames)

	return func(module core.Module) core.Provider {
		prd := module.NewProvider(core.ProviderOptions{
			Name: core.Provide(opt.Name),
			Factory: func(param ...interface{}) interface{} {
				pro, ok := param[0].(*PromptConfig)
				if ok && pro != nil {
					pro.Metrics = append(pro.Metrics, counter)
				}
				return counter
			},
			Inject: []core.Provide{PROMPT},
		})

		return prd
	}
}

func InjectCounter(ref core.RefProvider, name string) *prometheus.CounterVec {
	counter, ok := ref.Ref(core.Provide(name)).(*prometheus.CounterVec)
	if !ok || counter == nil {
		return nil
	}
	return counter
}

// Gauge
func MakeGaugeProvider(opt prometheus.GaugeOpts, labelNames []string) core.Providers {
	gauge := prometheus.NewGaugeVec(opt, labelNames)

	return func(module core.Module) core.Provider {
		prd := module.NewProvider(core.ProviderOptions{
			Name: core.Provide(opt.Name),
			Factory: func(param ...interface{}) interface{} {
				pro, ok := param[0].(*PromptConfig)
				if ok && pro != nil {
					pro.Metrics = append(pro.Metrics, gauge)
				}
				return gauge
			},
			Inject: []core.Provide{PROMPT},
		})

		return prd
	}
}

func InjectGauge(ref core.RefProvider, name string) *prometheus.GaugeVec {
	gauge, ok := ref.Ref(core.Provide(name)).(*prometheus.GaugeVec)
	if !ok || gauge == nil {
		return nil
	}
	return gauge
}

// Histogram
func MakeHistogramProvider(opt prometheus.HistogramOpts, labelNames []string) core.Providers {
	histogram := prometheus.NewHistogramVec(opt, labelNames)

	return func(module core.Module) core.Provider {
		prd := module.NewProvider(core.ProviderOptions{
			Name: core.Provide(opt.Name),
			Factory: func(param ...interface{}) interface{} {
				pro, ok := param[0].(*PromptConfig)
				if ok && pro != nil {
					pro.Metrics = append(pro.Metrics, histogram)
				}
				return histogram
			},
			Inject: []core.Provide{PROMPT},
		})

		return prd
	}
}

func InjectHistogram(ref core.RefProvider, name string) *prometheus.HistogramVec {
	histogram, ok := ref.Ref(core.Provide(name)).(*prometheus.HistogramVec)
	if !ok || histogram == nil {
		return nil
	}
	return histogram
}

// Summary
func MakeSummaryProvider(opt prometheus.SummaryOpts, labelNames []string) core.Providers {
	summary := prometheus.NewSummaryVec(opt, labelNames)

	return func(module core.Module) core.Provider {
		prd := module.NewProvider(core.ProviderOptions{
			Name: core.Provide(opt.Name),
			Factory: func(param ...interface{}) interface{} {
				pro, ok := param[0].(*PromptConfig)
				if ok && pro != nil {
					pro.Metrics = append(pro.Metrics, summary)
				}
				return summary
			},
			Inject: []core.Provide{PROMPT},
		})

		return prd
	}
}

func InjectSummary(ref core.RefProvider, name string) *prometheus.SummaryVec {
	summary, ok := ref.Ref(core.Provide(name)).(*prometheus.SummaryVec)
	if !ok || summary == nil {
		return nil
	}
	return summary
}
