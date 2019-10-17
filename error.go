package chat

import "errors"

var (
	// ErrNoSuchUser results when a message is sent to an offline or
	// non-existent user.
	ErrNoSuchUser = errors.New("no such user found")
)
