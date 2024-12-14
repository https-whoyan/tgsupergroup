# Telegram Supergroup Bot Library

---

A lightweight library for managing **topics** in Telegram supergroups. It provides tools for sending messages to topics by name and handling topic creation automatically when needed. This library is designed to simplify interaction with Telegram's topic-based messaging for bots.

---

### Features

- **Send messages to specific topics by name**.
- **Automatic topic management**: If the target topic doesn't exist, it will be created automatically (requires bot admin privileges).
- **Flexible storage options**:
    - In-memory storage for lightweight use cases.
    - Redis-based storage for scalable and persistent topic management.
- **Support for both synchronous and asynchronous message sending**.
- **Customizable options** for fine-tuning bot behavior.

---

### Installation

Install the library using `go get`:

```bash
go get github.com/https-whoyan/tgsupergroup@v1.0.0
```

---

### Usage Example

Here’s a quick example of how to use the library to send a message to a specific topic:

```go
package main

import (
	"context"
	"log"

	"github.com/https-whoyan/tgsupergroup"
)

func main() {
	// Replace with your bot's token
	var botToken string // Your token
	ctx := context.Background()

	// Initialize the bot
	bot, err := tgsupergroup.NewBot(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Replace with your chat's ID
	var yourChatID tgsupergroup.ChatID // Chat ID as int64

	// Send a message to a specific topic by its name
	err = bot.SendMessageToTopicByChatID(
		ctx, yourChatID, "example topic name", "example message",
	)
	if err != nil {
		log.Fatal(err)
	}
}
```

---

### Project Architecture

The library is structured for clarity and modularity. Below is an overview of the main directories and key files:

<pre>
<code style="display: block">
├── errors
|     └── Handles custom error definitions used across the library.
├── internal
|     └── Internal package for interacting with the Telegram Bot API. | Includes HTTP request handling and utility functions for marshaling/unmarshaling data. | 
├── types
|     └── Contains all data types used by the library. You can use is from tgsupergroup package
|
├── bot.go
|     └── The main entry point for the library, containing bot initialization and configuration options.
├── send.go
|     └── Implements synchronous message-sending methods for topics and chats.
├── async.go
|     └── Asynchronous versions of message-sending methods with error handling via callback.
└── storage.go
      └── Defines the Storage interface for managing topic-related data, with a Redis-based implementation.
</code>
</pre>
---
### Limitations

While the library validates most requirements programmatically, there are a few conditions to be aware of:

- The bot must have permission to send messages in the target group.
- To create topics, the bot must be an administrator in the group with sufficient privileges.
- If the library cannot find the topic by enumerating `ThreadID`s, it will attempt to create the topic. This might fail if the bot does not have the necessary permissions or if Telegram API limits are reached.

Ensure your bot has the required permissions to avoid errors during runtime.

---

### License

This library is licensed under the MIT License. See the `LICENSE` file for more details.
