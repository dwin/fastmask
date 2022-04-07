package fastmail

import (
	"context"

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
	return &Client{
		httpC: resty.New().SetBaseURL(APIEndpoint),
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

	if resp, err := request.Post(c.config.APIBaseURL); err != nil {
		return nil, APIError{Msg: "request failed", Status: resp.Status(), Detail: resp.String()}
	}

	return &JMAPResponse, nil
}
