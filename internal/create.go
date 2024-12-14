package internal

import (
	"context"
	"net/http"
)

type createTopicResponse struct {
	basicResponse
	Result *struct {
		ThreadID TopicThreadID `json:"message_thread_id"`
	}
}

func (r *requester) CreateTopic(ctx context.Context, topic *Topic) (*TopicThreadID, error) {
	req, err := r.newRequest(ctx, http.MethodPost, createTopic, queryArgs{
		chatIDJson:   toStr(topic.ChatID),
		chatNameJson: toStr(topic.Name),
	})
	if err != nil {
		return nil, err
	}
	var resp createTopicResponse
	err = r.send(req, &resp)
	if err != nil {
		return nil, err
	}
	err = r.checkBasicResponse(resp.basicResponse)
	if err != nil {
		return nil, err
	}
	return &resp.Result.ThreadID, nil
}
