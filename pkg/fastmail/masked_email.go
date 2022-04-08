package fastmail

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// MaskedEmail represents a Fastmail masked email.
type MaskedEmail struct {
	ID            string `json:"id,omitempty" mapstructure:"id"`
	State         string `json:"state,omitempty" mapstructure:"state"`
	Email         string `json:"email,omitempty" mapstructure:"email"`
	Description   string `json:"description,omitempty" mapstructure:"description"`
	ForDomain     string `json:"forDomain,omitempty" mapstructure:"forDomain"`
	URL           string `json:"url,omitempty" mapstructure:"url"`
	CreatedBy     string `json:"createdBy,omitempty" mapstructure:"createdBy"`
	CreatedAt     string `json:"createdAt,omitempty" mapstructure:"createdAt"`
	LastMessageAt string `json:"lastMessageAt,omitempty" mapstructure:"lastMessageAt"`
}

// MaskedEmailPayload is the payload for the MaskedEmail/{set,update} method.
type MaskedEmailPayload struct {
	AccountID string                  `json:"accountId,omitempty"`
	Create    map[string]*MaskedEmail `json:"create,omitempty"`
	Set       map[string]*MaskedEmail `json:"set,omitempty"`
	Update    map[string]*MaskedEmail `json:"update,omitempty"`
	Destroy   []string                `json:"destroy,omitempty"`
}

// CreateMaskedEmail creates a new masked email for the given forDomain domain.
// If `enabled` is set to false, will only create a pending email and needs to be confirmed before it's usable.
func (c *Client) CreateMaskedEmail(ctx context.Context, maskedEmail *MaskedEmail, enabled bool) (*MaskedEmail, error) {
	maskedEmail.State = isEnabledToString(enabled)

	request := JMAPRequest{
		Using: usingValueForMaskedEmail,
		MethodCalls: []MethodCall{{
			Name: "MaskedEmail/set",
			Payload: MaskedEmailPayload{
				AccountID: c.creds.accountID,
				Create: map[string]*MaskedEmail{
					c.config.AppName: maskedEmail,
				},
			},
			ID: "0",
		}},
	}

	res, err := c.sendRequest(ctx, &request)
	if err != nil {
		return nil, fmt.Errorf("send request error: %w", err)
	}
	// nolint:gomnd // ignore here.
	if len(res.MethodResponses) != 1 {
		return nil, MethodResponseError{len(res.MethodResponses), 1}
	}
	// nolint:gomnd // ignore here.
	if len(res.MethodResponses[0]) != 3 {
		return nil, MethodResponseError{len(res.MethodResponses[0]), 3}
	}

	var payload MethodResponseMaskedEmailSet

	if err := mapstructure.Decode(res.MethodResponses[0][1], &payload); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	created, err := payload.GetCreatedItem()
	if err != nil {
		return nil, fmt.Errorf("error getting created item: %w", err)
	}

	return &created, nil
}

// DeleteMaskedEmails deletes the given masked emails by ID.
func (c *Client) DeleteMaskedEmails(ctx context.Context, ids ...string) error {
	request := JMAPRequest{
		Using: usingValueForMaskedEmail,
		MethodCalls: []MethodCall{{
			Name: "MaskedEmail/set",
			Payload: &MaskedEmailPayload{
				AccountID: c.creds.accountID,
				Destroy:   ids,
			},
			ID: "0",
		}},
	}

	_, err := c.sendRequest(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}
