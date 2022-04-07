package fastmail

import (
	"encoding/json"
	"fmt"
	"time"
)

type MaskedEmail struct {
	ID            string     `json:"id,omitempty" mapstructure:"id"`
	State         string     `json:"state,omitempty" mapstructure:"state"`
	Email         string     `json:"email,omitempty" mapstructure:"email"`
	Description   string     `json:"description,omitempty" mapstructure:"description"`
	ForDomain     string     `json:"forDomain,omitempty" mapstructure:"forDomain"`
	URL           string     `json:"url,omitempty" mapstructure:"url"`
	CreatedBy     string     `json:"createdBy,omitempty" mapstructure:"createdBy"`
	CreatedAt     *time.Time `json:"createdAt,omitempty" mapstructure:"createdAt"`
	LastMessageAt *time.Time `json:"lastMessageAt,omitempty" mapstructure:"lastMessageAt"`
}

type JMAPRequest struct {
	Using       []string     `json:"using,omitempty"`
	MethodCalls []MethodCall `json:"methodCalls,omitempty"`
}

type JMAPResponse struct {
	LatestClientVersion string          `json:"latestClientVersion,omitempty"`
	MethodResponses     [][]interface{} `json:"methodResponses,omitempty"`
	SessionState        string          `json:"sessionState,omitempty"`
}

type MethodCall struct {
	Name    string
	Payload interface{}
	ID      string
}

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
		return fmt.Errorf("unmarshal method call response: %w", err)
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
	AccountID string                 `mapstructure:"accountId"`
	Created   map[string]MaskedEmail `mapstructure:"created"`
	Updated   map[string]interface{} `mapstructure:"updated"`
	Destroyed []interface{}          `mapstructure:"destroyed"`
	NewState  interface{}            `mapstructure:"newState"`
	OldState  interface{}            `mapstructure:"oldState"`
}

func (m *MethodResponseMaskedEmailSet) GetCreatedItem() (MaskedEmail, error) {
	for i := range m.Created {
		return m.Created[i], nil
	}

	return MaskedEmail{}, ErrNoItemsReturned
}
