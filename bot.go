package main

import (
	"context"
	"io"
	"net/http"
)

type Bot struct {
	io.Closer
	ctx      context.Context
	basicURL string
	token    string
	botName  string
	chat     *chat

	spamCount uint
	parseMode ParseMode

	httpCli *http.Client

	requester   *requester
	cacher      Storage
	chatCacher  map[ChatID]*chat
	topicsCache map[ChatID]*Topics
}

const maxSpamThreadIDs = 100

type Option func(*Bot)

func WithHTTPClient(httpCli *http.Client) Option {
	return func(bot *Bot) {
		bot.httpCli = httpCli
	}
}
func WithCacher(cacher Storage) Option {
	return func(bot *Bot) {
		bot.cacher = cacher
	}
}
func WithChatID(chatID int64) Option {
	return func(bot *Bot) {
		bot.chat = &chat{
			chatID: chatID,
		}
	}
}
func WithMaxSpamThreadIDs(maxThreadIDs uint) Option {
	return func(bot *Bot) {
		bot.spamCount = maxThreadIDs
	}
}
func WithContext(ctx context.Context) Option {
	return func(bot *Bot) {
		bot.ctx = ctx
	}
}
func WithParseMode(m ParseMode) Option {
	return func(bot *Bot) {
		bot.parseMode = m
	}
}
func WithBotName(name string) Option {
	return func(bot *Bot) {
		bot.botName = name
	}
}

func NewBot(token string, opts ...Option) (*Bot, error) {
	b := &Bot{
		ctx:         context.Background(),
		token:       token,
		httpCli:     http.DefaultClient,
		spamCount:   maxSpamThreadIDs,
		parseMode:   ParseModeMarkdownV2,
		chatCacher:  make(map[ChatID]*chat),
		topicsCache: make(map[int64]*Topics),
	}
	for _, opt := range opts {
		opt(b)
	}
	var err error
	b.requester = newRequester(token, b.httpCli, b.parseMode, b.botName)
	botName, pingErr := b.requester.getMe(b.ctx)
	if pingErr != nil {
		return nil, pingErr
	}
	if len(b.botName) == 0 {
		b.botName = botName
	}
	if b.chat != nil {
		b.chat, err = b.requester.getChat(b.ctx, b.chat.chatID)
		if err != nil {
			return nil, err
		}
		if b.chat == nil {
			return nil, ErrChatNotFound
		}
		b.chatCacher[b.chat.chatID] = b.chat
	}
	if b.cacher != nil && b.chat != nil {
		var allTopics *Topics
		allTopics, err = b.cacher.GetAll(b.ctx, b.chat.chatID)
		if err != nil {
			return nil, err
		}
		b.topicsCache[b.chat.chatID] = allTopics
	}
	return b, nil
}

func (b *Bot) Close() error {
	b.httpCli.CloseIdleConnections()
	return nil
}
