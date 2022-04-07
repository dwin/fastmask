package fastmail

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

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
		return nil, fmt.Errorf("expected 1 method response, got %d", len(res.MethodResponses))
	}
	// nolint:gomnd // ignore here.
	if len(res.MethodResponses[0]) != 3 {
		return nil, fmt.Errorf("expected 3 method response items, got %d", len(res.MethodResponses[0]))
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
