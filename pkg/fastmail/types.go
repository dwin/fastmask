package fastmail

import (
	"encoding/json"
	"fmt"
)

type JMAPRequest struct {
	Using       []string     `json:"using,omitempty"`
	MethodCalls []MethodCall `json:"methodCalls,omitempty"`
}

type JMAPResponse struct {
	LatestClientVersion string           `json:"latestClientVersion,omitempty"`
	MethodResponses     []MethodResponse `json:"methodResponses,omitempty"`
	SessionState        string           `json:"sessionState,omitempty"`
}

type MethodCall struct {
	Name    string
	Payload interface{}
	ID      string
}

type MethodResponse [3]interface{}

// MarshalJSON marshals a MethodCall into the format needed by the Fastmail API
// eg. ["MaskedEmail/set", {...}, "0"].
func (m *MethodCall) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal([]interface{}{m.Name, m.Payload, m.ID})
	if err != nil {
		return nil, fmt.Errorf("marshal method call response: %w", err)
	}

	return b, nil
}

func (m *MethodCall) UnmarshalJSON(b []byte) error {
	var v [3]interface{}

	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("unmarshal method call to interface error: %w", err)
	}

	name, ok := v[0].(string)
	if ok {
		m.Name = name
	}

	m.Payload = v[1]

	id, ok := v[2].(string)
	if ok {
		m.ID = id
	}

	return nil
}

type MethodResponseMaskedEmailSet struct {
	AccountID string                 `mapstructure:"accountId" json:"accountId,omitempty"`
	Created   map[string]MaskedEmail `mapstructure:"created" json:"created,omitempty"`
	Updated   map[string]interface{} `mapstructure:"updated" json:"updated,omitempty"`
	Destroyed []interface{}          `mapstructure:"destroyed" json:"destroyed,omitempty"`
	NewState  interface{}            `mapstructure:"newState" json:"newState,omitempty"`
	OldState  interface{}            `mapstructure:"oldState" json:"oldState,omitempty"`
}

func (m *MethodResponseMaskedEmailSet) GetCreatedItem() (MaskedEmail, error) {
	for i := range m.Created {
		return m.Created[i], nil
	}

	return MaskedEmail{}, ErrNoItemsReturned
}
