package tgsupergroup

import (
	"context"
	"os"
	"strconv"
	"testing"
)

func newBot() *Bot {
	loadDotEnv()
	token := os.Getenv("BOT_TOKEN")
	testChatID, err := strconv.ParseInt(os.Getenv("TEST_CHAT_ID"), 10, 64)
	if err != nil {
		panic(err)
	}
	storage := NewRedisStorage(newClient())
	bot, err := NewBot(token, WithStorage(storage), WithChatID(testChatID))
	if err != nil {
		panic(err)
	}
	return bot
}

func TestBot_SendMessage(t *testing.T) {
	bot := newBot()
	t.Run("send_to_contains_topic", func(t *testing.T) {
		err := bot.SendMessageToTopic(ctx, "test name", "test")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("send_to_not_contains_topic", func(t *testing.T) {
		err := bot.SendMessageToTopic(context.Background(), "Not contains topic", "test")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("send to chat", func(t *testing.T) {
		err := bot.SendMessage(ctx, "test")
		if err != nil {
			t.Error(err)
		}
	})
	t.Run("send_to_topics", func(t *testing.T) {
		err := bot.SendMessageToTopic(ctx, "test name", "test")
		if err != nil {
			t.Error(err)
		}
		err = bot.SendMessageToTopic(ctx, "Not contains topic", "test")
		if err != nil {
			t.Error(err)
		}
	})
}
