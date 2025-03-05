package prompt_test

import (
	"io"
	"math/rand"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/prompt"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Counter(t *testing.T) {
	middleware := func(ref core.RefProvider, ctx core.Ctx) error {
		counter := prompt.InjectCounterVec(ref, "http_requests_total")
		if counter != nil {
			method := ctx.Req().Method
			path := ctx.Req().URL.Path

			counter.WithLabelValues(path, method).Inc()
		}
		return ctx.Next()
	}
	appModule := func() core.Module {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{prompt.Register(&prompt.Config{
				Metrics: []prompt.Metric{{
					Name: "http_requests_total",
					Collector: prometheus.NewCounterVec(prometheus.CounterOpts{
						Name: "http_requests_total",
						Help: "Total number of HTTP requests received",
					}, []string{"path", "method"}),
				}},
			})},
		}).UseRef(middleware)

		return app
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	res, err := testClient.Get(testServer.URL + "/api")
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)

	res, err = testClient.Get(testServer.URL + "/api/metrics")
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	require.NotEmpty(t, string(data))
}

func Test_Gauge(t *testing.T) {
	middleware := func(ref core.RefProvider, ctx core.Ctx) error {
		gauge := prompt.InjectGauge(ref, "http_active_requests")
		if gauge != nil {
			gauge.Inc()
		}
		time.Sleep(1 * time.Second)
		_ = ctx.Next()

		gauge.Dec()
		return nil
	}
	appModule := func() core.Module {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{prompt.Register(&prompt.Config{
				Metrics: []prompt.Metric{{
					Name: "http_active_requests",
					Collector: prometheus.NewGauge(
						prometheus.GaugeOpts{
							Name: "http_active_requests",
							Help: "Number of active connections to the service",
						},
					),
				}},
			})},
		}).UseRef(middleware)

		return app
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	res, err := testClient.Get(testServer.URL + "/api")
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)

	res, err = testClient.Get(testServer.URL + "/api/metrics")
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	require.NotEmpty(t, string(data))
}

func Test_Histogram(t *testing.T) {
	middleware := func(ref core.RefProvider, ctx core.Ctx) error {
		histogram := prompt.InjectHistogramVec(ref, "http_request_duration_seconds")
		if histogram != nil {
			now := time.Now()

			delay := time.Duration(rand.Intn(900)) * time.Millisecond
			time.Sleep(delay)

			method := ctx.Req().Method
			path := ctx.Req().URL.Path

			histogram.With(prometheus.Labels{
				"method": method, "path": path,
			}).
				Observe(time.Since(now).Seconds())
		}
		return ctx.Next()
	}
	appModule := func() core.Module {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{prompt.Register(&prompt.Config{
				Metrics: []prompt.Metric{{
					Name: "http_request_duration_seconds",
					Collector: prometheus.NewHistogramVec(prometheus.HistogramOpts{
						Name:    "http_request_duration_seconds",
						Help:    "Duration of HTTP requests",
						Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10},
					}, []string{"path", "method"}),
				}},
			})},
		}).UseRef(middleware)

		return app
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	res, err := testClient.Get(testServer.URL + "/api")
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)

	res, err = testClient.Get(testServer.URL + "/api/metrics")
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	require.NotEmpty(t, string(data))
}

func Test_Summary(t *testing.T) {
	middleware := func(ref core.RefProvider, ctx core.Ctx) error {
		summary := prompt.InjectSummary(ref, "post_request_duration_seconds")
		if summary != nil {
			now := time.Now()

			summary.Observe(time.Since(now).Seconds())
		}
		return ctx.Next()
	}
	appModule := func() core.Module {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{prompt.Register(&prompt.Config{
				Metrics: []prompt.Metric{{
					Name: "post_request_duration_seconds",
					Collector: prometheus.NewSummary(prometheus.SummaryOpts{
						Name: "post_request_duration_seconds",
						Help: "Duration of requests to https://jsonplaceholder.typicode.com/posts",
						Objectives: map[float64]float64{
							0.5:  0.05,  // Median (50th percentile) with a 5% tolerance
							0.9:  0.01,  // 90th percentile with a 1% tolerance
							0.99: 0.001, // 99th percentile with a 0.1% tolerance
						},
					}),
				}},
			})},
		}).UseRef(middleware)

		return app
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	res, err := testClient.Get(testServer.URL + "/api")
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)

	res, err = testClient.Get(testServer.URL + "/api/metrics")
	require.Nil(t, err)
	require.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	require.NotEmpty(t, string(data))
}
