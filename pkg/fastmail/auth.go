package fastmail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
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
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar for login: %w", err)
	}

	c.httpC.Jar = cookieJar

	// Send the username to get a loginId
	usernameReqBody, err := json.Marshal(AuthenticateUsernameRequest{Username: username})
	if err != nil {
		return nil, fmt.Errorf("login loginId json marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, APIAuthEndpoint, bytes.NewBuffer(usernameReqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating loginId request: %w", err)
	}

	setContentTypeHeaders(req)

	resp, err := c.httpC.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login loginId request error: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("login read loginId response error: %w", err)
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, APIError{Msg: "Login flow get loginId failed", Status: resp.Status, Detail: string(respBody)}
	}

	var loginIDResp AuthFlowMessage

	if err := json.Unmarshal(respBody, &loginIDResp); err != nil {
		return nil, fmt.Errorf("login unmarshal loginId response error: %w", err)
	}

	// Send Password
	passwordAuthInput := AuthFlowMessage{
		LoginID: loginIDResp.LoginID,
		Type:    "password",
		Value:   password,
	}

	passwordReqBody, err := json.Marshal(&passwordAuthInput)
	if err != nil {
		return nil, fmt.Errorf("login marshal password request error: %w", err)
	}

	req, err = http.NewRequestWithContext(ctx, http.MethodPost, APIAuthEndpoint, bytes.NewBuffer(passwordReqBody))
	if err != nil {
		return nil, fmt.Errorf("login init password request error: %w", err)
	}

	setContentTypeHeaders(req)

	resp, err = c.httpC.Do(req)
	if err != nil {
		return nil, fmt.Errorf("login password request error: %w", err)
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("login password response ReadAll error: %w", err)
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, APIError{Msg: "Login flow send password failed", Status: resp.Status, Detail: string(respBody)}
	}

	mfaRequired := func() bool {
		for _, method := range loginIDResp.Methods {
			if method.Type == "totp" {
				return true
			}
		}

		return false
	}()

	if mfaRequired && mfaCode == "" {
		return nil, ErrMFARequired
	}

	// Send MFA if required
	if mfaRequired {
		mfaAuthInput := AuthFlowMessage{
			LoginID: loginIDResp.LoginID,
			Type:    "totp",
			Value:   mfaCode,
		}

		mfaReqBody, err := json.Marshal(&mfaAuthInput)
		if err != nil {
			return nil, fmt.Errorf("login MFA json marshal error: %w", err)
		}

		req, err = http.NewRequestWithContext(ctx, http.MethodPost, APIAuthEndpoint, bytes.NewBuffer(mfaReqBody))
		if err != nil {
			return nil, fmt.Errorf("login MFA init request error: %w", err)
		}

		setContentTypeHeaders(req)

		resp, err = c.httpC.Do(req)
		if err != nil {
			return nil, fmt.Errorf("login MFA request error: %w", err)
		}

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("login ReadAll MFA response error: %w", err)
		}

		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return nil, APIError{Msg: "Login flow send MFA failed, check MFA code", Status: resp.Status, Detail: string(respBody)}
		}
	}

	var authResponse AuthResponse

	if err := json.Unmarshal(respBody, &authResponse); err != nil {
		return nil, fmt.Errorf("login unmarshal auth response error: %w", err)
	}

	if authResponse.AccessToken == "" {
		return nil, ErrAccessTokenNotFound
	}

	return &authResponse, nil
}
