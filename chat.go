package tgsupergroup

import (
	"context"
	"github.com/https-whoyan/tgsupergroup/errors"
)

// Preload Chat Topics using Storage
func (b *Bot) PreloadChatTopics(ctx context.Context, chatIDs ...ChatID) error {
	if b.storage == nil {
		return nil
	}
	for _, id := range chatIDs {
		var err error
		chatCache, contains := b.chatCache[ChatID(id)]
		if !contains {
			chatByID, getChatErr := b.requester.GetChat(ctx, ChatID(id))
			if getChatErr != nil {
				return getChatErr
			}
			if chatByID.IsSuperGroup() {
				return errors.ErrChatIsNotSuperGroup
			}
			b.chatCache[ChatID(id)] = chatByID
		}
		if chatCache == nil {
			return errors.ErrChatNotFound
		}
		topicsCache, contains := b.topicsCache[ChatID(id)]
		if topicsCache == nil || !contains {
			topicsCache, err = b.storage.GetAll(ctx, id)
			if err != nil {
				return err
			}
			b.topicsCache[ChatID(id)] = topicsCache
		}
	}
	return nil
}

func (b *Bot) findChat(ctx context.Context, chatID ChatID) (*Chat, error) {
	cachedChat, contains := b.chatCache[chatID]
	if contains {
		if cachedChat == nil {
			return nil, errors.ErrChatNotFound
		}
		return cachedChat, nil
	}
	cachedChat, err := b.requester.GetChat(ctx, chatID)
	if err != nil {
		return nil, err
	}
	b.chatCache[chatID] = cachedChat
	return cachedChat, nil
}
