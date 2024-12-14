package tgsupergroup

import (
	"context"
	"github.com/https-whoyan/tgsupergroup/types"
	"strconv"
)

func idFromBytes(bytes []byte) (TopicThreadID, error) {
	threadID, err := strconv.ParseUint(string(bytes), 10, 64)
	if err != nil {
		return EmptyThreadID, err
	}
	return TopicThreadID(threadID), nil
}

func (b *Bot) safeTopicToLocalCacheIfNeed(topic *Topic) {
	if b.topicsCache[topic.ChatID] != nil {
		id := b.topicsCache[topic.ChatID].GetID(topic.Name)
		if id == EmptyThreadID {
			b.topicsCache[topic.ChatID].Safe(topic)
		}
		return
	}
	newTopics := make(Topics)
	newTopics.Safe(topic)
	b.topicsCache[topic.ChatID] = &newTopics
}

func (b *Bot) getTopic(
	ctx context.Context, chatID ChatID, topicName TopicName,
) (topicID TopicThreadID, found bool, err error) {
	topics, ok := b.topicsCache[chatID]
	if ok {
		if topics == nil {
			return b.getTopicByStorageOrSpam(ctx, chatID, topicName)
		}
		topicID = topics.GetID(topicName)
		if topicID == types.EmptyTopicID {
			return b.getTopicByStorageOrSpam(ctx, chatID, topicName)
		}
		return topicID, true, nil
	}
	return b.getTopicByStorageOrSpam(ctx, chatID, topicName)
}

func (b *Bot) getTopicByStorageOrSpam(
	ctx context.Context, chatID ChatID, topicName TopicName,
) (topicID TopicThreadID, found bool, err error) {
	if b.storage == nil {
		return b.findBySpam(ctx, chatID, topicName)
	}
	topic, err := b.storage.GetByName(ctx, chatID, topicName)
	if err != nil {
		return EmptyThreadID, false, err
	}
	if topic != nil {
		topicID = topic.ThreadID
		return topicID, true, b.storage.Save(ctx, topic)
	}
	topicID, found, err = b.findBySpam(ctx, chatID, topicName)
	if err != nil {
		return
	}
	if found {
		err = b.storage.Save(ctx, NewTopic(chatID, topicName, EmptyThreadID))
	}
	return
}

func (b *Bot) findBySpam(
	ctx context.Context, chatID ChatID, topicName TopicName,
) (topicID TopicThreadID, contains bool, err error) {
	for i := 1; i <= int(b.spamCount); i++ {
		ok, checkTopicErr := b.requester.CheckTopic(
			ctx, NewTopic(chatID, topicName, TopicThreadID(i)),
		)
		if checkTopicErr != nil {
			return EmptyThreadID, false, checkTopicErr
		}
		if !ok {
			continue
		}
		return TopicThreadID(i), true, nil
	}
	return EmptyThreadID, false, nil
}
