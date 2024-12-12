package tgsupergroup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type createTopicResponse struct {
	basicResponse
	Result *struct {
		ThreadID TopicThreadID `json:"message_thread_id"`
	}
}

func (r *requester) createTopic(ctx context.Context, topic *Topic) (TopicThreadID, error) {
	req, err := r.newRequest(ctx, http.MethodPost, createTopic, queryArgs{
		chatIDJson:      toStr(topic.ChatID),
		messageTextJson: fmt.Sprintf(MessageCreateAsset, r.botName),
		msgThreadIDJson: toStr(topic.ThreadID),
	})
	if err != nil {
		return EmptyTopicID, err
	}
	var createResp createTopicResponse
	err = r.send(req, &createResp)
	if err != nil {
		return EmptyTopicID, err
	}
	if !createResp.Ok {
		return EmptyTopicID, errors.New(*createResp.Desc)
	}
	return createResp.Result.ThreadID, nil
}
