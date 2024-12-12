package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const checkMessageAsset = "Test by %v"

type errResponse struct {
	ErrCode *int    `json:"error_code"`
	Desc    *string `json:"description"`
}

type basicResponse struct {
	errResponse
	Ok bool `json:"ok"`
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

func (r *requester) deleteMessage(ctx context.Context, chatID ChatID, messageID MessageID) error {
	req, err := r.newRequest(ctx, http.MethodPost, deleteMessage, queryArgs{
		chatIDJson:    toStr(chatID),
		messageIDJson: toStr(messageID),
	})
	if err != nil {
		return err
	}
	httpResp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	if httpResp.StatusCode != 200 {
		return errors.New(httpResp.Status)
	}
	return nil
}

func (r *requester) checkTopic(ctx context.Context, topic *Topic) (bool, error) {
	req, err := r.newRequest(ctx, http.MethodGet, sendMessage, queryArgs{
		chatIDJson:      toStr(topic.ChatID),
		messageTextJson: fmt.Sprintf(checkMessageAsset, r.botName),
		msgThreadIDJson: toStr(topic.ThreadID),
	})
	if err != nil {
		return false, err
	}
	var resp sendMessageResponse
	err = r.send(req, &resp)
	if err != nil {
		return false, err
	}
	if !resp.Ok {
		if *resp.Desc == messageThreadNotFound {
			return false, nil
		}
		return false, errors.New(*resp.Desc)
	}
	err = r.deleteMessage(ctx, topic.ChatID, resp.Result.MessageID)
	if err != nil {
		return false, err
	}
	if resp.Result.ReplyToMessage.ForumTopicResponse.Name != topic.Name {
		return false, nil
	}
	return true, nil
}
