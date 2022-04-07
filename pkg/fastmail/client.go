package fastmail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	httpC  *http.Client
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
		httpC: http.DefaultClient,
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

	return c
}

func (c *Client) setAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.creds.accessToken)
}

func setContentTypeHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func (c *Client) sendRequest(ctx context.Context, r *JMAPRequest) (*JMAPResponse, error) {
	requestJSON, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.APIBaseURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	c.setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request error: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, APIError{Status: resp.Status, Msg: "", Detail: string(respBody)}
	}

	var JMAPResponse JMAPResponse

	if err := json.Unmarshal(respBody, &JMAPResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &JMAPResponse, nil
}
