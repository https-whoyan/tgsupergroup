package internal

import "time"

const (
	basicURL      = "https://api.telegram.org/bot"
	getMeEndpoint = "/getMe"
	sendMessage   = "/sendMessage"
	createTopic   = "/createForumTopic"
	deleteMessage = "/deleteMessage"
)

const (
	chatIDJson      = "chat_id"
	msgThreadIDJson = "message_thread_id"
	messageIDJson   = "message_id"
	chatNameJson    = "name"
	messageTextJson = "text"
)

const (
	maxRequestsBySecond = 10
	requestTiming       = time.Second / maxRequestsBySecond
)
