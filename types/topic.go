package types

type TopicThreadID = uint
type TopicName = string

const EmptyTopicID TopicThreadID = 0

type Topic struct {
	ChatID   int64         `json:"chatID"`
	ThreadID TopicThreadID `json:"threadID"`
	Name     TopicName     `json:"name"`
}

func NewTopic(chatID ChatID, topicName TopicName, id TopicThreadID) *Topic {
	return &Topic{
		ChatID:   chatID,
		ThreadID: id,
		Name:     topicName,
	}
}

type Topics map[TopicName]*Topic

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

func (t *Topics) Len() int {
	return len(*t)
}
