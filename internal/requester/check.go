package requester

import (
	"context"
	"errors"
	"net/http"
)

func (r *requester) GetChat(ctx context.Context, chatID ChatID) (*Chat, error) {
	resp, err := r.sendMessageToChat(ctx, chatID, checkMessageAsset)
	if err != nil {
		return nil, err
	}
	err = r.deleteMessage(ctx, chatID, resp.Result.MessageID)
	if err != nil {
		return nil, err
	}
	return &Chat{
		ChatID:   chatID,
		ChatType: ChatType(resp.Result.ChatResponse.Type),
	}, nil
}

func (r *requester) CheckTopic(ctx context.Context, topic *Topic) (bool, error) {
	var (
		chatID    = topic.ChatID
		topicID   = topic.ThreadID
		topicName = topic.Name
	)
	resp, err := r.sendMessageToTopic(ctx, chatID, topicID, checkMessageAsset)
	if err != nil {
		if errors.Is(err, threadNotFoundErr) {
			return false, nil
		}
		return false, err
	}
	err = r.deleteMessage(ctx, chatID, resp.Result.MessageID)
	if err != nil {
		return false, err
	}
	if resp.Result.ReplyToMessage.ForumTopicResponse.Name != topicName {
		return false, nil
	}
	return true, nil
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
	if httpResp.StatusCode != http.StatusOK {
		return errors.New(httpResp.Status)
	}
	return nil
}
