package prompt

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const PROMPT core.Provide = "PROMPT"

type PromptConfig struct {
	Metrics []prometheus.Collector
	Opt     promhttp.HandlerOpts
}

func promptHandler(module core.Module) core.Controller {
	ctrl := module.NewController("/metrics")
	pro, ok := module.Ref(PROMPT).(*PromptConfig)
	if !ok || pro == nil {
		panic(errors.New("prometheus not config"))
	}

	if len(pro.Metrics) > 0 {
		reg := prometheus.NewRegistry()
		for _, metric := range pro.Metrics {
			reg.Register(metric)
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

func Register(opt promhttp.HandlerOpts) core.Modules {
	return func(module core.Module) core.Module {
		promptModule := module.New(core.NewModuleOptions{})

		promptModule.NewProvider(core.ProviderOptions{
			Name: PROMPT,
			Value: &PromptConfig{
				Opt: opt,
			},
		})

		promptModule.Controllers(promptHandler)

		return promptModule
	}
}
