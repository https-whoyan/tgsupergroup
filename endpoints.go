package main

import "time"

const (
	basicURL      = "https://api.telegram.org/bot"
	pingEndpoint  = "/getMe"
	sendMessage   = "/sendMessage"
	createTopic   = "/createForumTopic"
	deleteMessage = "/deleteMessage"
)

const (
	chatIDJson      = "chat_id"
	msgThreadIDJson = "msg_thread_id"
	messageIDJson   = "message_id"
	chatNameJson    = "name"
	messageTextJson = "text"
)

const (
	maxRequestsBySecond = 30
	requestTiming       = time.Second / (maxRequestsBySecond + 1)
)
