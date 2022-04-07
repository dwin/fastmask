package fastmail

import (
	"errors"
	"fmt"
)

var (
	ErrAccessTokenNotFound = errors.New("access token not found")
	ErrNoItemsReturned     = errors.New("no items returned")
	ErrMFARequired         = errors.New("mfa required for login")
)

type APIError struct {
	Status string
	Msg    string
	Detail string
}

func (a APIError) Error() string {
	return fmt.Sprintf("Fastmail API unexpected status code: '%s', msg: '%s', detail: '%s'", a.Status, a.Msg, a.Detail)
}
