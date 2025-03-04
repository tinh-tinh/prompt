package prompt

import (
	"errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const PROMPT core.Provide = "PROMPT"

type Config struct {
	Metrics []Metric
	Opt     promhttp.HandlerOpts
}

func promptHandler(module core.Module) core.Controller {
	ctrl := module.NewController("/metrics")
	pro, ok := module.Ref(PROMPT).(*Config)
	if !ok || pro == nil {
		panic(errors.New("prometheus not config"))
	}

	if len(pro.Metrics) > 0 {
		reg := prometheus.NewRegistry()
		for _, metric := range pro.Metrics {
			reg.Register(metric.Collector)
		}
		handler := promhttp.HandlerFor(
			reg,
			pro.Opt,
		)

		ctrl.Handler("", handler)
	} else {
		ctrl.Handler("", promhttp.Handler())
	}

	return ctrl
}

func Register(config *Config) core.Modules {
	return func(module core.Module) core.Module {
		promptModule := module.New(core.NewModuleOptions{})

		promptModule.NewProvider(core.ProviderOptions{
			Name:  PROMPT,
			Value: config,
		})

		if len(config.Metrics) > 0 {
			for _, metric := range config.Metrics {
				promptModule.NewProvider(core.ProviderOptions{
					Name:  GetMetricName(metric.Name),
					Value: metric.Collector,
				})
			}
		}

		for _, v := range promptModule.GetDataProviders() {
			fmt.Println(v.GetName())
		}

		promptModule.Controllers(promptHandler)

		return promptModule
	}
}

func GetMetricName(name string) core.Provide {
	modelName := "Metric_" + name

	return core.Provide(modelName)
}
