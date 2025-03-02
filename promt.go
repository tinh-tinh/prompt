package prompt

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const PROMPT core.Provide = "PROMPT"

func promptHandler(module core.Module) core.Controller {
	ctrl := module.NewController("/metrics")

	reg := prometheus.NewRegistry()
	handler := promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{},
	)

	ctrl.Handler("", handler)

	return ctrl
}

func Register(cs ...prometheus.Collector) core.Modules {
	return func(module core.Module) core.Module {
		promptModule := module.New(core.NewModuleOptions{})

		promptModule.NewProvider(core.ProviderOptions{
			Name:  PROMPT,
			Value: cs,
		})

		promptModule.Controllers(promptHandler)

		return promptModule
	}
}
