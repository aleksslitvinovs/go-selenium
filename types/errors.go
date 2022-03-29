package types

import "errors"

var (
	ErrInvalidParameters = errors.New("invalid parameters")

	ErrFailedRequest = errors.New("failed to execute request")
)
