package tgsupergroup

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Storage interface {
	GetAll(ctx context.Context, chatID ChatID) (*Topics, error)
	Save(ctx context.Context, topic *Topic) error
	GetByName(ctx context.Context, chatID ChatID, name TopicName) (*Topic, error)
}

type redisCacher struct {
	client     *redis.Client
	fieldAsset string
}

// Initiates Storage using the redis.Client
func NewRedisStorage(client *redis.Client) Storage {
	return &redisCacher{
		client:     client,
		fieldAsset: RedisKeyAsset,
	}
}

const RedisKeyAsset = "TGSGroups:%d" // ChatID

func (c redisCacher) getKey(chatID int64) string {
	return fmt.Sprintf(c.fieldAsset, chatID)
}

func (c redisCacher) GetAll(ctx context.Context, chatID ChatID) (*Topics, error) {
	mapped, err := c.client.HGetAll(ctx, c.getKey(chatID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var out = make(Topics, len(mapped))
	for topicName, value := range mapped {
		id, idsMarshalErr := idFromBytes([]byte(value))
		if idsMarshalErr != nil {
			return nil, idsMarshalErr
		}
		out[topicName] = NewTopic(chatID, topicName, id)
	}
	return &out, nil
}

func (c redisCacher) Save(ctx context.Context, topic *Topic) error {
	chatID := topic.ChatID
	bytes, err := c.client.HGet(ctx, c.getKey(chatID), topic.Name).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return c.client.HSet(ctx, c.getKey(chatID), topic.Name, topic.ThreadID).Err()
		}
		return err
	}
	id, err := idFromBytes([]byte(bytes))
	if err != nil {
		return err
	}
	// If already have, safe smallest
	id = min(id, topic.ThreadID)
	return c.client.HSet(ctx, c.getKey(chatID), topic.Name, id).Err()
}

func (c redisCacher) GetByName(ctx context.Context, chatID ChatID, name TopicName) (*Topic, error) {
	bytesIDs, err := c.client.HGet(ctx, c.getKey(chatID), name).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
	}
	id, err := idFromBytes(bytesIDs)
	if err != nil {
		return nil, err
	}
	return NewTopic(chatID, name, id), nil
}
