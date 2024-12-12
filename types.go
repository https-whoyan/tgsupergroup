package tgsupergroup

import (
	types2 "github.com/https-whoyan/tgsupergroup/types"
)

type (
	Chat   = types2.Chat
	ChatID = types2.ChatID

	Topic         = types2.Topic
	Topics        = types2.Topics
	TopicThreadID = types2.TopicThreadID
	TopicName     = types2.TopicName

	ParseMode = types2.ParseMode
)

const (
	ParseModeHTML       = types2.ParseModeHTML
	ParseModeMarkdown   = types2.ParseModeMarkdown
	ParseModeMarkdownV2 = types2.ParseModeMarkdownV2

	SuperGroupType = types2.SuperGroupType
)

var (
	EmptyThreadID = types2.EmptyTopicID
)
