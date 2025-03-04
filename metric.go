package prompt

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Metric struct {
	Name      string
	Collector prometheus.Collector
}

func InjectCounter(module core.RefProvider, name string) *prometheus.CounterVec {
	fmt.Println(GetMetricName(name))
	metric, ok := module.Ref(GetMetricName(name)).(*prometheus.CounterVec)
	if !ok {
		return nil
	}

	return metric
}
