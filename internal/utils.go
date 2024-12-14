package internal

import (
	"errors"
	"fmt"
	"strings"

	myErrors "github.com/https-whoyan/tgsupergroup/errors"
)

func (r *requester) escapeF(messageText string, args ...interface{}) string {
	return r.parseMode.EscapeText(fmt.Sprintf(messageText, args...))
}

const (
	messageThreadNotFoundDescription = "Bad Request: message thread not found"
	chatNotFoundDescription          = "Bad Request: chat not found"
	notFoundDescription              = "Not Found"
	messageIsNotForum                = "Bad Request: the chat is not a forum"
	noRootsDescription               = "Bad Request: not enough rights"
)

var (
	threadNotFoundErr = errors.New("thread not found")
)

func (*requester) checkBasicResponse(resp basicResponse) (parsedErr error) {
	if resp.Ok {
		return nil
	}
	respErrDesc := *resp.Desc
	switch respErrDesc {
	case messageThreadNotFoundDescription:
		return threadNotFoundErr
	case chatNotFoundDescription:
		return errChatNotFound
	case messageIsNotForum:
		return myErrors.ErrChatIsNotSuperGroup
	}
	if strings.HasPrefix(respErrDesc, noRootsDescription) {
		return myErrors.ErrNotEnoughPrivileges
	}
	return errors.New(respErrDesc)
}
