package fastmail

import (
	"context"
	"net/http"
	"testing"

	"github.com/icrowley/fake"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func Test_Client(t *testing.T) {
	appName := fake.CharactersN(10)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient(appName)
	httpmock.ActivateNonDefault(client.httpC.GetClient()) // needed for to mock Resty.

	ctx := context.TODO()

	require.NotNil(t, client)
	require.NotNil(t, client.config)
	require.Equal(t, appName, client.config.AppName)

	t.Run("Send Request - Error", func(t *testing.T) {
		errorResponder := httpmock.NewStringResponder(http.StatusInternalServerError, "test error message")

		httpmock.RegisterResponder(http.MethodPost, APIEndpoint, errorResponder)

		resp, err := client.sendRequest(ctx, &JMAPRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
		require.ErrorContains(t, err, "unexpected response")
		require.ErrorContains(t, err, "test error message")
	})
}
