package tgsupergroup

import (
	"github.com/https-whoyan/tgsupergroup/types"
)

type (
	Chat   = types.Chat
	ChatID = types.ChatID

	Topic         = types.Topic
	Topics        = types.Topics
	TopicThreadID = types.TopicThreadID
	TopicName     = types.TopicName

	ParseMode = types.ParseMode
)

const (
	ParseModeHTML       = types.ParseModeHTML
	ParseModeMarkdown   = types.ParseModeMarkdown
	ParseModeMarkdownV2 = types.ParseModeMarkdownV2

	SuperGroupType = types.SuperGroupType
)

var (
	NewTopic      = types.NewTopic
	EmptyThreadID = types.EmptyTopicID
)
