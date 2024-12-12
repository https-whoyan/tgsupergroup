package main

import "errors"

var (
	ErrChatNotFound        = errors.New("chat not found")
	ErrChatIsNotSuperGroup = errors.New("chat is not super group")
)

const (
	messageThreadNotFound = "Bad Request: message thread not found"
)
