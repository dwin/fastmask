package fastmail

import (
	"context"
	"net/http"
	"testing"

	"github.com/icrowley/fake"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func Test_Masked_Email(t *testing.T) {
	appName := fake.CharactersN(10)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient(appName)
	require.NotNil(t, client)
	httpmock.ActivateNonDefault(client.httpC.GetClient()) // needed for to mock Resty.

	client.SetTokenAuthCredentials("fakeAccountID", "fakeAccessToken")

	ctx := context.TODO()

	createMaskedEmailInput := &MaskedEmail{
		ForDomain:   fake.DomainName(),
		Description: fake.Sentence(),
	}

	t.Run("Test Create Masked Email", func(t *testing.T) {
		defer httpmock.Reset()

		createResponder, err := httpmock.NewJsonResponder(http.StatusOK, httpmock.File("examples/create_masked_response.json"))
		require.NoError(t, err)

		httpmock.RegisterResponder(http.MethodPost, APIEndpoint, createResponder)

		result, err := client.CreateMaskedEmail(ctx, createMaskedEmailInput, true)
		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("Test Create Masked Email - Auth Failure", func(t *testing.T) {
		defer httpmock.Reset()

		createResponder, err := httpmock.NewJsonResponder(http.StatusUnauthorized, struct{}{})
		require.NoError(t, err)

		httpmock.RegisterResponder(http.MethodPost, APIEndpoint, createResponder)

		result, err := client.CreateMaskedEmail(ctx, createMaskedEmailInput, true)
		require.Error(t, err)
		require.Nil(t, result)
		require.ErrorIs(t, err, ErrUnauthorized)
	})
}
