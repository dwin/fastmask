package fastmail

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var (
	// APIEndpoint is the Fastmail API endpoint.
	APIEndpoint = "https://api.fastmail.com/jmap/api/"
	// APIAuthEndpoint is the Fastmail authentication endpoint.
	APIAuthEndpoint = "https://www.fastmail.com/jmap/authenticate/"
)

var usingValueForMaskedEmail = []string{
	"urn:ietf:params:jmap:core",
	"https://www.fastmail.com/dev/maskedemail",
}

type Client struct {
	httpC  *resty.Client
	config *ClientConfig
	creds  *Credentials
}

type ClientConfig struct {
	AppName    string
	APIBaseURL string
}

type Credentials struct {
	accountID   string
	accessToken string
}

// Returns NewClient with the given values. 'accountID' is the Fastmail account ID, this is
// not the same as the email address.
func NewClient(appName string) *Client {
	httpC := resty.New()
	httpC.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		if r.StatusCode() == http.StatusOK || r.StatusCode() == http.StatusCreated {
			return nil
		}

		if r.StatusCode() == http.StatusUnauthorized {
			return ErrUnauthorized
		}

		return APIError{Msg: "unexpected response", Code: r.StatusCode(), Status: r.Status(), Detail: r.String()}
	})

	return &Client{
		httpC: httpC,
		config: &ClientConfig{
			AppName:    appName,
			APIBaseURL: APIEndpoint,
		},
	}
}

func (c *Client) SetTokenAuthCredentials(accountID, accessToken string) *Client {
	c.creds = &Credentials{
		accountID:   accountID,
		accessToken: accessToken,
	}

	c.httpC.SetAuthToken(accessToken)

	return c
}

func (c *Client) sendRequest(ctx context.Context, r *JMAPRequest) (*JMAPResponse, error) {
	var JMAPResponse JMAPResponse

	request := c.httpC.R().SetBody(r)
	request.SetContext(ctx)
	request.SetResult(&JMAPResponse)

	if _, err := request.Post(c.config.APIBaseURL); err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return &JMAPResponse, nil
}
