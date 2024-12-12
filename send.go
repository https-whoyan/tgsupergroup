package tgsupergroup

import "context"

func (b *Bot) SendMessageToTopic(ctx context.Context, topicName TopicName, messageText string) error {
	if b.chat == nil {
		return ErrNotProvidedChatID
	}
	if b.chat.chatType != superGroupType {
		return ErrChatIsNotSuperGroup
	}
	return b.SendMessageToTopicByChatID(ctx, b.chat.chatID, topicName, messageText)
}

func (b *Bot) SendMessageToTopicByChatID(
	ctx context.Context, chatID ChatID, topicName TopicName, messageText string,
) error {
	botChat, err := b.findChat(ctx, chatID)
	if err != nil {
		return err
	}
	if botChat.chatType != superGroupType {
		return ErrChatIsNotSuperGroup
	}
	topicID, isContains, err := b.getTopic(ctx, chatID, topicName)
	if err != nil {
		return err
	}
	var topic = NewTopic(chatID, topicName, EmptyTopicID)
	defer func() { b.safeTopicToLocalCacheIfNeed(topic) }()
	if isContains {
		topic.ThreadID = topicID
		err = b.requester.sendMessageToTopic(ctx, topic, messageText)
		return err
	}
	topicID, err = b.requester.createTopic(ctx, topic)
	if err != nil {
		return err
	}
	topic.ThreadID = topicID
	if b.cacher != nil {
		err = b.cacher.Save(ctx, topic)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) SendMessageToChat(ctx context.Context, chatID ChatID, messageText string) error {
	botChat, err := b.findChat(ctx, chatID)
	if err != nil {
		return err
	}
	if botChat == nil {
		return ErrChatNotFound
	}
	err = b.requester.sendMessageToChat(ctx, chatID, messageText)
	return err
}

func (b *Bot) SendMessage(ctx context.Context, messageText string) error {
	if b.chat == nil {
		return ErrNotProvidedChatID
	}
	return b.requester.sendMessageToChat(ctx, b.chat.chatID, messageText)
}
