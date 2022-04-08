package fastmail

import (
	"context"
	"net/http"
	"testing"

	"github.com/icrowley/fake"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	// nolint:gosec // fakeAccessToken is a fake access token matching the one in the successful_auth_response.json file.
	fakeAccessToken = "fmb1-testtest-mv6PnlZnjmuB57v7ywfmxHM47NTiQyqN-bKhQxhpmPr-l7A2DI4DPm0Etn3ltjT7IiUnWJHtc0mv"
	fakeAccountID   = "z3U9VLb8H"
)

var (
	appName  = fake.CharactersN(10)
	ctx      = context.TODO()
	loginID  = fake.CharactersN(10)
	username = fake.EmailAddress()
	password = fake.CharactersN(18)
)

func Test_Auth(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient(appName)
	httpmock.ActivateNonDefault(client.httpC.GetClient()) // needed for to mock Resty.

	loginIDResponder, err := httpmock.NewJsonResponder(http.StatusOK, map[string]interface{}{
		"loginId": loginID,
	})
	require.NoError(t, err)

	// Test: Successfully authenticate with username and password, and get the account ID and access token.
	t.Run("Successful Login Flow - No MFA", func(t *testing.T) {
		defer httpmock.Reset()

		successResponder, err := httpmock.NewJsonResponder(http.StatusOK, httpmock.File("examples/successful_auth_response.json"))
		require.NoError(t, err)

		httpmock.RegisterResponder(http.MethodPost, APIAuthEndpoint, loginIDResponder.Then(successResponder))

		resp, err := client.LoginUsernamePasswordMFA(ctx, username, password, "")
		require.NoError(t, err, "LoginUsernamePasswordMFA should not return an error")
		require.NotNil(t, resp, "LoginUsernamePasswordMFA should return a non-nil response")

		accountIDRes, ok := resp.GetMailAccountID()
		require.True(t, ok, "account ID is expected and was not found in response")
		require.Equal(t, fakeAccountID, accountIDRes)
		require.Equal(t, fakeAccessToken, resp.GetAccessToken())
	})

	t.Run("Failed Login - Require MFA, not provided", func(t *testing.T) {
		defer httpmock.Reset()

		requireMFAResponder, err := httpmock.NewJsonResponder(http.StatusOK, httpmock.File("examples/require_mfa_response.json"))
		require.NoError(t, err)

		httpmock.RegisterResponder(http.MethodPost, APIAuthEndpoint, loginIDResponder.Then(requireMFAResponder))

		resp, err := client.LoginUsernamePasswordMFA(ctx, username, password, "")
		require.ErrorIs(t, err, ErrMFARequired)
		require.Nil(t, resp, "LoginUsernamePasswordMFA should return a nil response when error is returned")
	})

	t.Run("Successful Login Flow - With MFA", func(t *testing.T) {
		defer httpmock.Reset()

		requireMFAResponder, err := httpmock.NewJsonResponder(http.StatusOK, httpmock.File("examples/require_mfa_response.json"))
		require.NoError(t, err)

		successResponder, err := httpmock.NewJsonResponder(http.StatusOK, httpmock.File("examples/successful_auth_response.json"))
		require.NoError(t, err)

		httpmock.RegisterResponder(http.MethodPost, APIAuthEndpoint, loginIDResponder.Then(requireMFAResponder).Then(successResponder))

		resp, err := client.LoginUsernamePasswordMFA(ctx, username, password, "123456")
		require.NoError(t, err, "LoginUsernamePasswordMFA should not return an error")
		require.NotNil(t, resp, "LoginUsernamePasswordMFA should return a non-nil response")

		accountIDRes, ok := resp.GetMailAccountID()
		require.True(t, ok, "account ID is expected and was not found in response")
		require.Equal(t, fakeAccountID, accountIDRes)
		require.Equal(t, fakeAccessToken, resp.GetAccessToken())
	})
}
