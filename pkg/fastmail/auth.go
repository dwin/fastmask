package fastmail

import (
	"context"
	"encoding/json"
	"fmt"
)

type AuthResponse struct {
	AccountType      string                 `json:"accountType,omitempty"`
	UserID           string                 `json:"userId,omitempty"`
	SigningID        string                 `json:"signingId,omitempty"`
	SigningKey       string                 `json:"signingKey,omitempty"`
	IsAdmin          bool                   `json:"isAdmin,omitempty"`
	SessionKey       string                 `json:"sessionKey,omitempty"`
	PrimaryAccounts  map[string]string      `json:"primaryAccounts,omitempty"`
	AccessToken      string                 `json:"accessToken,omitempty"`
	DisplayName      string                 `json:"displayName"`
	APIURL           string                 `json:"apiUrl,omitempty"`
	DownloadURL      string                 `json:"downloadUrl,omitempty"`
	UploadURL        string                 `json:"uploadUrl"`
	IsLoginBlocked   bool                   `json:"isLoginBlocked"`
	State            string                 `json:"state"`
	IsReseller       bool                   `json:"isReseller"`
	IsAppStore       bool                   `json:"isAppStore"`
	IsVerified       bool                   `json:"isVerified"`
	MayMailSync      bool                   `json:"mayMailSync"`
	MayContactSync   bool                   `json:"mayContactSync"`
	FeatureLevel     int                    `json:"featureLevel"`
	IsReceiveBlocked bool                   `json:"isReceiveBlocked"`
	ProxyURL         string                 `json:"proxyUrl,omitempty"`
	IsSendBlocked    bool                   `json:"isSendBlocked,omitempty"`
	NewCalAlerts     bool                   `json:"newCalAlerts,omitempty"`
	Language         string                 `json:"language"`
	Capabilities     map[string]interface{} `json:"capabilities,omitempty"`
}

func (a *AuthResponse) GetMailAccountID() (accountID string, ok bool) {
	if a.PrimaryAccounts == nil {
		return "", false
	}

	accountID, ok = a.PrimaryAccounts["https://www.fastmail.com/dev/mail"]

	return
}

func (a *AuthResponse) GetAccessToken() string {
	return a.AccessToken
}

type AuthenticateUsernameRequest struct {
	Username string `json:"username,omitempty"`
}

type AuthFlowMessage struct {
	LoginID        string       `json:"loginId,omitempty"`
	MayTrustDevice bool         `json:"mayTrustDevice,omitempty"`
	Remember       bool         `json:"remember"`
	Type           string       `json:"type,omitempty"`
	Methods        []AuthMethod `json:"methods,omitempty"`
	Value          string       `json:"value,omitempty"`
}

func (a *AuthFlowMessage) TOTPRequired() bool {
	for i := range a.Methods {
		if a.Methods[i].Type == "totp" {
			return true
		}
	}

	return false
}

type AuthMethod struct {
	Type         string `json:"type"`
	PhoneNumbers []struct {
		Number     string `json:"number"`
		ID         string `json:"id"`
		IsCodeSent bool   `json:"isCodeSent"`
	} `json:"phoneNumbers,omitempty"`
}

// LoginUsernamePasswordMFA authenticates with the given username and password, mfaCode is optional
// based on account settings.
func (c *Client) LoginUsernamePasswordMFA(ctx context.Context, username, password, mfaCode string) (*AuthResponse, error) {
	// Send the username to get a loginId
	var loginIDResult AuthFlowMessage

	loginIDRequest := c.httpC.R().SetBody(AuthenticateUsernameRequest{Username: username})
	loginIDRequest.SetContext(ctx)
	loginIDRequest.SetResult(&loginIDResult)

	if resp, err := loginIDRequest.Post(APIAuthEndpoint); err != nil {
		return nil, APIError{Msg: "get loginID failed", Status: resp.Status(), Detail: resp.String()}
	}

	// Send Password
	passwordAuthInput := AuthFlowMessage{
		LoginID: loginIDResult.LoginID,
		Type:    "password",
		Value:   password,
	}

	var authResponse AuthResponse

	passwordAuthRequest := c.httpC.R().SetBody(passwordAuthInput)
	passwordAuthRequest.SetContext(ctx)
	passwordAuthRequest.SetResult(&authResponse)

	resp, err := passwordAuthRequest.Post(APIAuthEndpoint)
	if err != nil {
		return nil, APIError{Msg: "password auth failed", Status: resp.Status(), Detail: resp.String()}
	}

	// If access token is received then we are authenticated, otherwise check for MFA requirement.
	if authResponse.AccessToken != "" {
		return &authResponse, nil
	}

	var passwordAuthResult AuthFlowMessage

	if err := json.Unmarshal(resp.Body(), &passwordAuthResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal auth response: %w", err)
	}

	mfaIsRequired := passwordAuthResult.TOTPRequired()

	if mfaIsRequired && mfaCode == "" {
		return nil, ErrMFARequired
	}

	if mfaIsRequired {
		mfaAuthInput := AuthFlowMessage{
			LoginID: passwordAuthResult.LoginID,
			Type:    "totp",
			Value:   mfaCode,
		}

		mfaAuthRequest := c.httpC.R().SetBody(mfaAuthInput)
		mfaAuthRequest.SetContext(ctx)
		mfaAuthRequest.SetResult(&authResponse)

		if resp, err := mfaAuthRequest.Post(APIAuthEndpoint); err != nil {
			return nil, APIError{Msg: "mfa auth failed", Status: resp.Status(), Detail: resp.String()}
		}
	}

	if authResponse.AccessToken == "" {
		return nil, ErrAccessTokenNotFound
	}

	return &authResponse, nil
}
