package internal

import (
	"context"
	"io"
	"net/http"
	"sync"
)

type Requester interface {
	io.Closer
	GetMe(ctx context.Context) (botName string, err error)
	GetChat(ctx context.Context, chatID ChatID) (*Chat, error)
	SendMessageToChat(ctx context.Context, chatID ChatID, messageText string, args ...interface{}) error
	SendMessageToTopic(
		ctx context.Context, chatID ChatID, topicID TopicThreadID, messageText string, args ...interface{},
	) error
	CheckTopic(ctx context.Context, topic *Topic) (bool, error)
	CreateTopic(ctx context.Context, topic *Topic) (topicID *TopicThreadID, err error)
}

type requester struct {
	mu         sync.Mutex
	basicURL   string
	botToken   string
	parseMode  ParseMode
	httpClient *http.Client
}

func NewRequester(
	botToken string,
	httpCli *http.Client,
	parseMode ParseMode,
) Requester {
	if httpCli == nil {
		httpCli = http.DefaultClient
	}
	return &requester{
		basicURL:   basicURL,
		botToken:   botToken,
		httpClient: httpCli,
		parseMode:  parseMode,
	}
}
