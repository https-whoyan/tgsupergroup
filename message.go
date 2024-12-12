package tgsupergroup

import (
	"context"
	"errors"
	"net/http"
)

var ErrNotProvidedChatID = errors.New("not provided chat ID")

func (r *requester) sendMessageToTopic(ctx context.Context, topic *Topic, messageText string) error {
	req, err := r.newRequest(ctx, http.MethodPost, sendMessage, queryArgs{
		chatIDJson:      toStr(topic.ChatID),
		messageTextJson: r.parseMode.EscapeText(messageText),
		msgThreadIDJson: toStr(topic.ThreadID),
	})
	if err != nil {
		return err
	}
	return r.send(req, nil)
}

func (r *requester) sendMessageToChat(ctx context.Context, chatID ChatID, messageText string) error {
	req, err := r.newRequest(ctx, http.MethodPost, sendMessage, queryArgs{
		chatIDJson:      toStr(chatID),
		messageTextJson: r.parseMode.EscapeText(messageText),
	})
	if err != nil {
		return err
	}
	return r.send(req, nil)
}
