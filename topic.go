package tgsupergroup

import (
	"context"
	"strconv"
)

const superGroupType = "supergroup"

type TopicThreadID = uint
type TopicName = string

const EmptyTopicID TopicThreadID = 0

type Topic struct {
	ChatID   int64         `json:"chatID"`
	ThreadID TopicThreadID `json:"threadID"`
	Name     TopicName     `json:"name"`
}

type Topics map[TopicName]*Topic

func idFromBytes(bytes []byte) (TopicThreadID, error) {
	threadID, err := strconv.ParseUint(string(bytes), 10, 64)
	if err != nil {
		return EmptyTopicID, err
	}
	return TopicThreadID(threadID), nil
}

func NewTopic(chatID ChatID, topicName TopicName, id TopicThreadID) *Topic {
	return &Topic{
		ChatID:   chatID,
		ThreadID: id,
		Name:     topicName,
	}
}

func (t *Topics) Safe(topic *Topic) {
	(*t)[topic.Name] = topic
}

func (t *Topics) GetID(name TopicName) TopicThreadID {
	if t == nil {
		return EmptyTopicID
	}
	topic, ok := (*t)[name]
	if !ok {
		return EmptyTopicID
	}
	if topic == nil {
		return EmptyTopicID
	}
	return topic.ThreadID
}

func (b *Bot) safeTopicToLocalCacheIfNeed(topic *Topic) {
	if b.topicsCache[topic.ChatID] != nil {
		id := b.topicsCache[topic.ChatID].GetID(topic.Name)
		if id == EmptyTopicID {
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
		if topicID == EmptyTopicID {
			return b.getTopicByStorageOrSpam(ctx, chatID, topicName)
		}
		return topicID, true, nil
	}
	return b.getTopicByStorageOrSpam(ctx, chatID, topicName)
}

func (b *Bot) getTopicByStorageOrSpam(
	ctx context.Context, chatID ChatID, topicName TopicName,
) (topicID TopicThreadID, found bool, err error) {
	if b.cacher != nil {
		topic, err := b.cacher.GetByName(ctx, chatID, topicName)
		if err != nil {
			return EmptyTopicID, false, err
		}
		if topic != nil {
			topicID = topic.ThreadID
			return topicID, true, b.cacher.Save(ctx, topic)
		}
	}
	topicID, found, err = b.findBySpam(ctx, chatID, topicName)
	if err != nil {
		return
	}
	if found {
		err = b.cacher.Save(ctx, NewTopic(chatID, topicName, EmptyTopicID))
	}
	return
}

func (b *Bot) findBySpam(
	ctx context.Context, chatID ChatID, topicName TopicName,
) (topicID TopicThreadID, contains bool, err error) {
	for i := 0; i <= int(b.spamCount); i++ {
		ok, checkTopicErr := b.requester.checkTopic(
			ctx, NewTopic(chatID, topicName, TopicThreadID(i)),
		)
		if checkTopicErr != nil {
			return EmptyTopicID, false, checkTopicErr
		}
		if !ok {
			continue
		}
		return TopicThreadID(i), true, nil
	}
	return EmptyTopicID, false, nil
}
