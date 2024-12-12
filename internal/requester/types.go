package requester

import (
	"github.com/https-whoyan/tgsupergroup/errors"
	"github.com/https-whoyan/tgsupergroup/types"
)

type (
	Chat          = types.Chat
	ChatID        = types.ChatID
	ChatType      = types.ChatType
	Topic         = types.Topic
	Topics        = types.Topics
	TopicThreadID = types.TopicThreadID
	TopicName     = types.TopicName
	ParseMode     = types.ParseMode
	MessageID     = types.MessageID
)

var (
	toStr           = types.ToStr
	errChatNotFound = errors.ErrChatNotFound
)
