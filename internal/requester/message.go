package requester

import (
	"context"
	"net/http"
)

func (r *requester) SendMessageToChat(ctx context.Context, chatID ChatID, messageText string, args ...interface{}) error {
	_, err := r.sendMessageToChat(ctx, chatID, messageText, args...)
	return err
}

func (r *requester) SendMessageToTopic(
	ctx context.Context, chatID ChatID, topicID TopicThreadID, messageText string, args ...interface{},
) error {
	_, err := r.sendMessageToTopic(ctx, chatID, topicID, messageText, args...)
	return err
}

func (r *requester) sendMessageToChat(
	ctx context.Context, chatID ChatID, messageText string, args ...interface{},
) (*sendMessageResponse, error) {
	req, err := r.newRequest(ctx, http.MethodGet, sendMessage, queryArgs{
		chatIDJson:      toStr(chatID),
		messageTextJson: r.escapeF(messageText, args),
	})
	if err != nil {
		return nil, err
	}
	var dst sendMessageResponse
	err = r.send(req, &dst)
	if err != nil {
		return nil, err
	}
	err = r.checkBasicResponse(dst.basicResponse)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

func (r *requester) sendMessageToTopic(
	ctx context.Context, chatID ChatID, topicID TopicThreadID, messageText string, args ...interface{},
) (*sendMessageResponse, error) {
	req, err := r.newRequest(ctx, http.MethodGet, sendMessage, queryArgs{
		chatIDJson:      toStr(chatID),
		messageTextJson: r.escapeF(messageText, args),
		msgThreadIDJson: toStr(topicID),
	})
	if err != nil {
		return nil, err
	}
	var dst sendMessageResponse
	err = r.send(req, &dst)
	if err != nil {
		return nil, err
	}
	err = r.checkBasicResponse(dst.basicResponse)
	if err != nil {
		return nil, err
	}
	return &dst, nil
}

type sendMessageResponse struct {
	basicResponse
	Result *struct {
		messageResponse
		ChatResponse   chatResponse           `json:"chat"`
		ReplyToMessage replyToMessageResponse `json:"reply_to_message"`
	} `json:"result"`
	MessageThreadID *TopicThreadID `json:"message_thread_id"`
}

type chatResponse struct {
	ChatID ChatID `json:"chat_id"`
	Type   string `json:"type"`
}

type messageResponse struct {
	MessageID MessageID `json:"message_id"`
}

type replyToMessageResponse struct {
	ForumTopicResponse struct {
		Name TopicName `json:"name"`
	} `json:"forum_topic_created"`
}
