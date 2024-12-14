package tgsupergroup

import (
	"context"
	"github.com/https-whoyan/tgsupergroup/errors"
)

/*
Use this method if you provided chatID in the options during initialization.

Otherwise, use SendMessageToTopicByChatID method.
*/
func (b *Bot) SendMessageToTopic(
	ctx context.Context, topicName TopicName, messageText string, args ...interface{},
) error {
	if b.chat == nil {
		return errors.ErrNotProvidedChatID
	}
	if !b.chat.IsSuperGroup() {
		return errors.ErrChatIsNotSuperGroup
	}
	return b.SendMessageToTopicByChatID(ctx, b.chat.ChatID, topicName, messageText, args...)
}

// Sends a message to a group's topic by its name.
func (b *Bot) SendMessageToTopicByChatID(
	ctx context.Context, chatID ChatID, topicName TopicName, messageText string, args ...interface{},
) error {
	currChat, err := b.findChat(ctx, chatID)
	if err != nil {
		return err
	}
	if !currChat.IsSuperGroup() {
		return errors.ErrChatIsNotSuperGroup
	}
	topicID, isContains, err := b.getTopic(ctx, chatID, topicName)
	if err != nil {
		return err
	}
	var topic = NewTopic(chatID, topicName, EmptyThreadID)
	defer func() { b.safeTopicToLocalCacheIfNeed(topic) }()
	if isContains {
		topic.ThreadID = topicID
		return b.requester.SendMessageToTopic(ctx, chatID, topic.ThreadID, messageText, args...)
	}
	topicIDPtr, createTopicErr := b.requester.CreateTopic(ctx, topic)
	if createTopicErr != nil {
		return createTopicErr
	}
	topicID = *topicIDPtr
	topic.ThreadID = topicID
	if b.storage != nil {
		err = b.storage.Save(ctx, topic)
		if err != nil {
			return err
		}
	}
	return b.requester.SendMessageToTopic(ctx, chatID, topicID, messageText, args...)
}

// Just send message to chatID
func (b *Bot) SendMessageToChat(ctx context.Context, chatID ChatID, messageText string, args ...interface{}) error {
	botChat, err := b.findChat(ctx, chatID)
	if err != nil {
		return err
	}
	if botChat == nil {
		return errors.ErrChatNotFound
	}
	err = b.requester.SendMessageToChat(ctx, chatID, messageText, args...)
	return err
}

// Just send message to prepared chatID
func (b *Bot) SendMessage(ctx context.Context, messageText string, args ...interface{}) error {
	if b.chat == nil {
		return errors.ErrNotProvidedChatID
	}
	return b.requester.SendMessageToChat(ctx, b.chat.ChatID, messageText, args...)
}

// May be useful if you have a threadID
func (b *Bot) SendMessageToTopicByID(
	ctx context.Context, chatID ChatID, topicID TopicThreadID, messageText string, args ...interface{},
) error {
	currChat, err := b.findChat(ctx, chatID)
	if err != nil {
		return err
	}
	if !currChat.IsSuperGroup() {
		return errors.ErrChatIsNotSuperGroup
	}
	return b.requester.SendMessageToTopic(ctx, chatID, topicID, messageText, args...)
}
