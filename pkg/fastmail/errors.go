package fastmail

import (
	"errors"
	"fmt"
)

var (
	ErrAccessTokenNotFound = errors.New("access token not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNoItemsReturned     = errors.New("no items returned")
	ErrMFARequired         = errors.New("mfa required for login")
)

type APIError struct {
	Code   int
	Status string
	Msg    string
	Detail string
}

func (a APIError) Error() string {
	return fmt.Sprintf("fastmail api unexpected status: '%d - %s', msg: '%s', detail: '%s'", a.Code, a.Status, a.Msg, a.Detail)
}

type MethodResponseError struct {
	Actual   int
	Expected int
}

func (m MethodResponseError) Error() string {
	return fmt.Sprintf("fastmail api returned unexpected method response , have method %d responses and expected: %d", m.Actual, m.Expected)
}
