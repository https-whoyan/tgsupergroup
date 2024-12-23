package tgsupergroup

import (
	"context"
	"github.com/https-whoyan/tgsupergroup/internal"
	"io"
	"net/http"
)

/*
Main structure for operand topics.

Main method : SendMessageToTopicByChatID

The method calls the telegram API, which is responsible for sending a message to a specific topic, and for that it needs the ThreadID of the chat, not the chat name.
To find a id, bot has 3 levels of get this id.

The first one is local, the bot hashes chats and topics in the group.

The second is Storage, which stores the group's tops in any database. The library already implements storage using Redis. Storage is configured with the WithStorage option. To get storage via Redis, use the NewRedisStorage function
Otherwise, if nothing is found, the bot finds the right topic by enumerating the IDs. The maximum standard value is 100, but this is Option WithMaxSpamThreadIDs

If the desired chat is not found, a new one is created. Make sure that the bot has enough permissions to create topics.
*/
type Bot struct {
	io.Closer
	ctx      context.Context
	basicURL string
	chat     *Chat

	spamCount uint
	parseMode ParseMode

	httpCli *http.Client

	requester   internal.Requester
	storage     Storage
	chatCache   map[ChatID]*Chat
	topicsCache map[ChatID]*Topics
}

const maxSpamThreadIDs = 100

// List:
// WithHTTPClient, WithStorage, WithChatID, WithMaxSpamThreadIDs, WithContext, WithParseMode, WithBotName
type Option func(*Bot)

// Configure http.Client. Default: http.DefaultClient
func WithHTTPClient(httpCli *http.Client) Option {
	return func(bot *Bot) {
		bot.httpCli = httpCli
	}
}

// Configure Storage to save TopicThreadID. Use NewRedisStorage if you're comfortable using redis.
//
// May be nil. (No Storage)
func WithStorage(cacher Storage) Option {
	return func(bot *Bot) {
		bot.storage = cacher
	}
}

// Configure default chatID, if you need to operand only with this chat.
func WithChatID(chatID ChatID) Option {
	return func(bot *Bot) {
		bot.chat = &Chat{
			ChatID: chatID,
		}
	}
}

// Configure max attemps to find a need topic. Default: 100
func WithMaxSpamThreadIDs(maxThreadIDs uint) Option {
	return func(bot *Bot) {
		bot.spamCount = maxThreadIDs
	}
}

// Configure context.Context. Default: context.Background
func WithContext(ctx context.Context) Option {
	return func(bot *Bot) {
		bot.ctx = ctx
	}
}

// Configure ParseMode. Default: ParseModeMarkdownV2
func WithParseMode(m ParseMode) Option {
	return func(bot *Bot) {
		bot.parseMode = m
	}
}

// NewBot Returns an instance of the bot, look for Option's above
func NewBot(token string, opts ...Option) (*Bot, error) {
	b := &Bot{
		ctx:         context.Background(),
		httpCli:     http.DefaultClient,
		spamCount:   maxSpamThreadIDs,
		parseMode:   ParseModeMarkdownV2,
		chatCache:   make(map[ChatID]*Chat),
		topicsCache: make(map[int64]*Topics),
	}
	for _, opt := range opts {
		opt(b)
	}
	b.requester = internal.NewRequester(token, b.httpCli, b.parseMode)
	optionsErr := b.initOptions()
	if optionsErr != nil {
		return nil, optionsErr
	}
	return b, nil
}

func NewBotsGroup(tokens []string, opts ...Option) (*Bot, error) {
	b := &Bot{
		ctx:         context.Background(),
		spamCount:   maxSpamThreadIDs,
		parseMode:   ParseModeMarkdownV2,
		chatCache:   make(map[ChatID]*Chat),
		topicsCache: make(map[int64]*Topics),
	}
	for _, opt := range opts {
		opt(b)
	}
	req, err := internal.NewBotsGroup(b.ctx, b.parseMode, tokens...)
	if err != nil {
		return nil, err
	}
	b.requester = req
	optionsErr := b.initOptions()
	if optionsErr != nil {
		return nil, optionsErr
	}
	return b, nil
}

func (b *Bot) initOptions() error {
	var err error
	_, pingErr := b.requester.GetMe(b.ctx)
	if pingErr != nil {
		return pingErr
	}
	if b.chat != nil {
		b.chat, err = b.requester.GetChat(b.ctx, b.chat.ChatID)
		if err != nil {
			return err
		}
		b.chatCache[b.chat.ChatID] = b.chat
	}
	if b.storage != nil && b.chat != nil {
		var allTopics *Topics
		allTopics, err = b.storage.GetAll(b.ctx, b.chat.ChatID)
		if err != nil {
			return err
		}
		b.topicsCache[b.chat.ChatID] = allTopics
	}
	return nil
}

func (b *Bot) Close() error {
	return b.requester.Close()
}
