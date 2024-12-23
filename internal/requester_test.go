package internal

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"

	myErrors "github.com/https-whoyan/tgsupergroup/errors"
	"github.com/https-whoyan/tgsupergroup/types"

	"github.com/joho/godotenv"
)

var ctx = context.Background()

func loadDotEnv() error {
	return godotenv.Load("../.env")
}

func newRequester(token string) Requester {
	return NewRequester(
		token,
		http.DefaultClient,
		types.ParseModeMarkdownV2,
	)
}

func newValidRequester() Requester {
	err := loadDotEnv()
	if err != nil {
		panic(err)
	}
	return newRequester(os.Getenv("BOT_TOKEN"))
}

func Test_GetMe(t *testing.T) {
	t.Run("invalid token", func(t *testing.T) {
		req := newRequester("")
		_, err := req.GetMe(ctx)
		if !errors.Is(err, myErrors.ErrInvalidToken) {
			t.Errorf("GetMe() error = %v, wantErr %v", err, myErrors.ErrInvalidToken)
		}
	})
	t.Run("valid token", func(t *testing.T) {
		req := newValidRequester()
		reqBotName := os.Getenv("BOT_NAME")
		botNameAns, err := req.GetMe(ctx)
		if err != nil {
			t.Errorf("GetMe() error = %v, wantErr %v", err, nil)
		}
		if reqBotName != botNameAns {
			t.Errorf("GetMe() = %v, want %v", botNameAns, reqBotName)
		}
	})
}

func Test_Send_To_Chat(t *testing.T) {
	t.Run("invalid chat ID", func(t *testing.T) {
		req := newValidRequester()
		err := req.SendMessageToChat(ctx, 0, "test msg")
		if !errors.Is(err, errChatNotFound) {
			t.Errorf("SendMessageToChat() error = %v, wantErr %v", err, errChatNotFound)
		}
	})
	t.Run("valid chat ID", func(t *testing.T) {
		req := newValidRequester()
		int64ChatID, err := strconv.ParseInt(os.Getenv("TEST_CHAT_ID"), 10, 64)
		if err != nil {
			t.Errorf("ParseInt() error = %v, wantErr %v", err, nil)
		}
		err = req.SendMessageToChat(ctx, int64ChatID, "test msg")
		if err != nil {
			t.Errorf("SendMessageToChat() error = %v, wantErr %v", err, nil)
		}
	})
}

func Test_Send_To_Topic(t *testing.T) {
	t.Run("invalid chatID", func(t *testing.T) {
		req := newValidRequester()
		err := req.SendMessageToTopic(ctx, 0, 0, "test msg")
		if !errors.Is(err, errChatNotFound) {
			t.Errorf("SendMessageToChat() error = %v, wantErr %v", err, errChatNotFound)
		}
	})
	t.Run("valid chat ID, invalid topicID", func(t *testing.T) {
		req := newValidRequester()
		int64ChatID, err := strconv.ParseInt(os.Getenv("TEST_CHAT_ID"), 10, 64)
		if err != nil {
			t.Errorf("ParseInt() error = %v, wantErr %v", err, nil)
		}
		err = req.SendMessageToTopic(ctx, int64ChatID, 100000, "test msg")
		if !errors.Is(err, threadNotFoundErr) {
			t.Errorf("SendMessageToChat() error = %v, wantErr %v", err, threadNotFoundErr)
		}
	})
	t.Run("valid chat ID, valid topicID", func(t *testing.T) {
		req := newValidRequester()
		int64ChatID, err := strconv.ParseInt(os.Getenv("TEST_CHAT_ID"), 10, 64)
		if err != nil {
			t.Errorf("ParseInt() error = %v, wantErr %v", err, nil)
		}
		uint64TopicID, err := strconv.ParseUint(os.Getenv("TEST_TOPIC_ID"), 10, 64)
		err = req.SendMessageToTopic(ctx, int64ChatID, TopicThreadID(uint64TopicID), "test msg")
		if err != nil {
			t.Errorf("SendMessageToTopic() error = %v, wantErr %v", err, nil)
		}
	})
}

func Test_Create(t *testing.T) {
	t.Run("no roots", func(t *testing.T) {
		req := newValidRequester()
		int64ChatID, err := strconv.ParseInt(os.Getenv("TEST_CHAT_ID_2"), 10, 64)
		if err != nil {
			t.Errorf("ParseInt() error = %v, wantErr %v", err, nil)
		}
		_, err = req.CreateTopic(ctx, &Topic{
			ChatID: int64ChatID,
			Name:   "TEST",
		})
		if !errors.Is(err, myErrors.ErrNotEnoughPrivileges) {
			t.Errorf("Create() error = %v, wantErr %v", err, myErrors.ErrNotEnoughPrivileges)
		}
	})
	t.Run("not superforum", func(t *testing.T) {
		req := newValidRequester()
		int64ChatID, err := strconv.ParseInt(os.Getenv("TEST_CHAT_ID_3"), 10, 64)
		if err != nil {
			t.Errorf("ParseInt() error = %v, wantErr %v", err, nil)
		}
		_, err = req.CreateTopic(ctx, &Topic{
			ChatID: int64ChatID,
			Name:   "TEST",
		})
		if !errors.Is(err, myErrors.ErrChatIsNotSuperGroup) {
			t.Errorf("Create() error = %v, wantErr %v", err, myErrors.ErrChatIsNotSuperGroup)
		}
	})
	t.Run("all valid", func(t *testing.T) {
		req := newValidRequester()
		int64ChatID, err := strconv.ParseInt(os.Getenv("TEST_CHAT_ID"), 10, 64)
		if err != nil {
			t.Errorf("ParseInt() error = %v, wantErr %v", err, nil)
		}
		_, err = req.CreateTopic(ctx, &Topic{
			ChatID: int64ChatID,
			Name:   "TEST",
		})
		if err != nil {
			t.Errorf("Create() error = %v, wantErr %v", err, nil)
		}
	})
}
