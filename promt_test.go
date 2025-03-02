package prompt_test

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/prompt"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Promp(t *testing.T) {
	appModule := func() core.Module {
		app := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{prompt.Register()},
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
}
