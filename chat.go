package tgsupergroup

import (
	"context"
	"errors"
	"net/http"
)

const chatNotFoundDescription = "Bad Request: chat not found"

type chat struct {
	chatID   ChatID
	chatType string
}

func (b *Bot) findChat(ctx context.Context, chatID ChatID) (*chat, error) {
	cachedChat, contains := b.chatCacher[chatID]
	if contains {
		if cachedChat == nil {
			return nil, ErrChatNotFound
		}
		return cachedChat, nil
	}
	cachedChat, err := b.requester.getChat(ctx, chatID)
	if err != nil {
		return nil, err
	}
	b.chatCacher[chatID] = cachedChat
	if cachedChat == nil {
		return nil, ErrChatNotFound
	}
	return cachedChat, nil
}

func (r *requester) getChat(ctx context.Context, chatID ChatID) (*chat, error) {
	req, err := r.newRequest(ctx, http.MethodPost, sendMessage, queryArgs{
		chatIDJson:      toStr(chatID),
		messageTextJson: "test",
	})
	if err != nil {
		return nil, err
	}
	var resp sendMessageResponse
	err = r.send(req, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		if *resp.Desc == chatNotFoundDescription {
			return nil, nil
		}
		return nil, errors.New(*resp.Desc)
	}
	outChat := &chat{
		chatID:   chatID,
		chatType: resp.Result.ChatResponse.Type,
	}
	err = r.deleteMessage(ctx, chatID, resp.Result.MessageID)
	if err != nil {
		return outChat, err
	}
	return outChat, nil
}
