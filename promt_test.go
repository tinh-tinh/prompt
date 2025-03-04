package prompt_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/prompt"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Promp(t *testing.T) {
	appModule := func() core.Module {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{prompt.Register(&prompt.Config{})},
		})

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

func Test_Counter(t *testing.T) {
	middleware := func(ref core.RefProvider, ctx core.Ctx) error {
		counter := prompt.InjectCounter(ref, "http_requests_total")
		fmt.Println(counter)
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
	fmt.Println(string(data))
}
