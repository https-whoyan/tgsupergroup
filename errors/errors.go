package errors

import "errors"

var (
	ErrChatNotFound        = errors.New("chat not found")
	ErrChatIsNotSuperGroup = errors.New("chat is not supergroup")
	ErrNotProvidedChatID   = errors.New("not provided chat ID")
	ErrInvalidToken        = errors.New("invalid token")
)
