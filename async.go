package tgsupergroup

import (
	"context"
	"github.com/https-whoyan/tgsupergroup/errors"
)

// Handler type for async versions of sending messages
type ErrorHandlerFunc = func(ctx context.Context, err error)

var mockErrorHandler = func(ctx context.Context, err error) {}

// Async version of SendMessageToTopic
func (b *Bot) AsyncSendMessageToTopic(
	ctx context.Context, onError ErrorHandlerFunc, topicName TopicName, messageText string, args ...interface{},
) {
	if onError == nil {
		onError = mockErrorHandler
	}
	go func() {
		if b.chat == nil {
			onError(ctx, errors.ErrNotProvidedChatID)
			return
		}
		if !b.chat.IsSuperGroup() {
			onError(ctx, errors.ErrChatIsNotSuperGroup)
			return
		}
		err := b.SendMessageToTopicByChatID(ctx, b.chat.ChatID, topicName, messageText, args...)
		if err != nil {
			onError(ctx, err)
		}
	}()
}

// Async version of SendMessageToTopicByChatID
func (b *Bot) AsyncSendMessageToTopicByChatID(
	ctx context.Context, onError ErrorHandlerFunc,
	chatID ChatID, topicName TopicName, messageText string, args ...interface{},
) {
	if onError == nil {
		onError = mockErrorHandler
	}
	go func() {
		currChat, err := b.findChat(ctx, chatID)
		if err != nil {
			onError(ctx, err)
			return
		}
		if !currChat.IsSuperGroup() {
			onError(ctx, errors.ErrChatIsNotSuperGroup)
			return
		}
		topicID, isContains, err := b.getTopic(ctx, chatID, topicName)
		if err != nil {
			onError(ctx, err)
			return
		}
		var topic = NewTopic(chatID, topicName, EmptyThreadID)
		defer func() { b.safeTopicToLocalCacheIfNeed(topic) }()
		if isContains {
			topic.ThreadID = topicID
			err = b.requester.SendMessageToTopic(ctx, chatID, topic.ThreadID, messageText, args...)
			if err != nil {
				onError(ctx, err)
			}
			return
		}
		topicIDPtr, createTopicErr := b.requester.CreateTopic(ctx, topic)
		if createTopicErr != nil {
			onError(ctx, createTopicErr)
			return
		}
		topicID = *topicIDPtr
		topic.ThreadID = topicID
		if b.storage != nil {
			err = b.storage.Save(ctx, topic)
			if err != nil {
				onError(ctx, err)
				return
			}
		}
		err = b.requester.SendMessageToTopic(ctx, chatID, topicID, messageText, args...)
		if err != nil {
			onError(ctx, err)
		}
	}()
}

// Async version of SendMessageToChat
func (b *Bot) AsyncSendMessageToChat(
	ctx context.Context, onError ErrorHandlerFunc, chatID ChatID, messageText string, args ...interface{},
) {
	if onError == nil {
		onError = mockErrorHandler
	}
	go func() {
		botChat, err := b.findChat(ctx, chatID)
		if err != nil {
			onError(ctx, err)
			return
		}
		if botChat == nil {
			onError(ctx, errors.ErrNotProvidedChatID)
		}
		err = b.requester.SendMessageToChat(ctx, chatID, messageText, args...)
		if err != nil {
			onError(ctx, err)
		}
	}()
}

// Async version of SendMessage
func (b *Bot) AsyncSendMessage(ctx context.Context, onError ErrorHandlerFunc, messageText string, args ...interface{}) {
	if onError == nil {
		onError = mockErrorHandler
	}
	go func() {
		if b.chat == nil {
			onError(ctx, errors.ErrNotProvidedChatID)
			return
		}
		err := b.requester.SendMessageToChat(ctx, b.chat.ChatID, messageText, args...)
		if err != nil {
			onError(ctx, err)
		}
	}()
}

// Async version of SendMessageToTopicByID
func (b *Bot) AsyncSendMessageToTopicByID(
	ctx context.Context, onError ErrorHandlerFunc,
	chatID ChatID, topicID TopicThreadID, messageText string, args ...interface{},
) {
	if onError == nil {
		onError = mockErrorHandler
	}
	go func() {
		currChat, err := b.findChat(ctx, chatID)
		if err != nil {
			onError(ctx, err)
			return
		}
		if !currChat.IsSuperGroup() {
			onError(ctx, errors.ErrChatIsNotSuperGroup)
			return
		}
		err = b.requester.SendMessageToTopic(ctx, chatID, topicID, messageText, args...)
		if err != nil {
			onError(ctx, err)
		}
	}()
}
