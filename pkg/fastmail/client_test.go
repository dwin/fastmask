package fastmail

import (
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/require"
)

func Test_Client(t *testing.T) {
	appName := fake.CharactersN(10)

	client := NewClient(appName)

	require.NotNil(t, client)
	require.NotNil(t, client.config)
	require.Equal(t, appName, client.config.AppName)
}
